package fixed

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"strings"
	"unsafe"
)

const (
	mask64S  = 1 << 63
	size64M  = 52
	mask64E  = (1<<11 - 1) << size64M
	mask64M  = (1 << size64M) - 1
	shift64E = size64M
	bias64   = (1 << (11 - 1)) - 1
)

const Precision = 20 //配置定点数精度 精度值为：1/2**Precision
const precisionBitsNum = Precision
const decimalBitsMask = 1<<precisionBitsNum - 1
const MaxFixed64 = Fixed64((1 << 63) - 1)
const MinFixed64 = ^Fixed64(0)
const Fixed64Zero = Fixed64(0)
const Fixed64Nan = Fixed64(1 << 63) //非数

type Fixed64 uint64

//规范数处理
func Float64ToFixed64(value float64) Fixed64 {
	var valueBits = math.Float64bits(value)
	s := valueBits & mask64S
	e := valueBits & mask64E
	m := valueBits&mask64M | (1 << size64M)
	realE := int64(e>>shift64E) - bias64
	fixedD := uint64(0)
	var result uint64
	var fixedP uint64 = 0

	if realE >= 0 {
		pBitsNum := size64M - realE
		if pBitsNum < 0 {
			if precisionBitsNum-pBitsNum > 63-53 {
				panic("Fixed number: part digital overflow") //超出定点整数表示范围
			}

			fixedD = m << -pBitsNum << precisionBitsNum
		} else {
			fixedD = m >> pBitsNum << precisionBitsNum
			pBitsFlowNum := pBitsNum - precisionBitsNum
			pMask := uint64((1 << pBitsNum) - 1)

			if pBitsFlowNum <= 0 {
				fixedP = m & pMask
			} else {
				fixedP = roundOdd(m&pMask, uint64(pBitsFlowNum)) >> pBitsFlowNum
			}
		}
	} else {
		allMBitsNum := 52 - realE
		pFlowBitsNum := allMBitsNum - precisionBitsNum
		if pFlowBitsNum > 53 {
			return 0
		}
		fixedP = roundOdd(m, uint64(pFlowBitsNum)) >> pFlowBitsNum
		//fixedP = roundOdd(m,uint64((size64M - precisionBitsNum) + (realE * -1))) >> ((size64M - precisionBitsNum) + (realE * -1))
	}

	result |= s
	result |= fixedD
	result |= fixedP
	return Fixed64(result)
}

func SafeFloat64ToFixed64(value float64) (Fixed64, error) {
	if math.IsNaN(value) {
		return 0, errors.New(string(Uint64Bits(math.Float64bits(value))) + " is NaN ")
	}
	if math.IsInf(value, 0) {
		return 0, errors.New(string(Uint64Bits(math.Float64bits(value))) + " is Inf ")
	}

	return Float64ToFixed64(value), nil
}

func (fixed Fixed64) Add(oth Fixed64) Fixed64 {
	var fS = uint64(fixed) & mask64S
	var oS = uint64(oth) & mask64S
	fixed &^= mask64S
	oth &^= mask64S
	if fS == oS {
		return Fixed64(uint64(fixed+oth) | fS)
	} else {
		if fixed > oth {
			return Fixed64(uint64(fixed-oth) | fS)
		} else if fixed < oth {
			return Fixed64(uint64(oth-fixed) | oS)
		} else {
			return 0
		}
	}
}

func (fixed Fixed64) Sub(oth Fixed64) Fixed64 {
	return fixed.Add(oth ^ mask64S)
}

func (fixed Fixed64) Mul(oth Fixed64) Fixed64 {
	var fS = fixed & mask64S
	var oS = oth & mask64S
	fixed &^= mask64S
	oth &^= mask64S

	hi, lo := bits.Mul64(uint64(fixed), uint64(oth))
	lo = roundOdd(lo, precisionBitsNum) >> precisionBitsNum

	hi = (hi & decimalBitsMask) << (64 - precisionBitsNum)
	if hi == 0 && lo == 0{
		return 0
	}

	return Fixed64(uint64(fS^oS) | hi | lo)
}

func (fixed Fixed64) Div(oth Fixed64) Fixed64 {
	var fS = fixed & mask64S
	var oS = oth & mask64S
	fixed &^= mask64S
	oth &^= mask64S

	quo, _ := bits.Div64(uint64(fixed>>(64-precisionBitsNum)), uint64(fixed<<precisionBitsNum), uint64(oth))
	if quo == 0{
		return 0
	}
	return Fixed64(uint64(fS^oS) | quo)
}

func (fixed Fixed64) Abs() Fixed64 {
	return fixed &^ mask64S
}

func (fixed Fixed64) Equal(oth Fixed64) bool {
	return fixed == oth
}

func (fixed Fixed64) Less(oth Fixed64) bool {
	return (fixed ^ mask64S) < (oth ^ mask64S)
}

func (fixed Fixed64) Great(oth Fixed64) bool {
	return (fixed ^ mask64S) > (oth ^ mask64S)
}

func (fixed Fixed64) Float64() float64 {
	number := uint64(fixed &^ mask64S)
	idx := bits.Len64(number)
	if idx != 0 {
		e := idx - precisionBitsNum - 1
		number = ((1 << idx) - 1) & number
		if idx > size64M {
			number >>= idx - size64M - 1
		} else {
			number <<= size64M - idx + 1
		}
		number = number &^ mask64E
		number |= uint64(fixed & mask64S)
		number |= uint64((e + bias64) << size64M)
	}

	return *(*float64)(unsafe.Pointer(&number))
}

func (fixed Fixed64) Int64() int64 {
	return int64((fixed&^mask64S)>>precisionBitsNum) * (1 - int64((fixed&mask64S)>>63))
}

func (fixed Fixed64) String() string {
	//return strconv.FormatFloat(fixed.Float64(), 'f', -1, 64)
	s := strings.TrimRightFunc(fmt.Sprintf("%f", fixed.Float64()), func(r rune) bool {
		if r == '0' || r == '.' {
			return true
		}
		return false
	})

	if s == "" {
		s = "0"
	}
	return s
}

func Uint64Bits(v uint64) (r []byte) {
	r = make([]byte, 64)
	for i := 0; i < 64; i++ {
		r[63-i] = '0'
		if (v & (1 << i)) > 0 {
			r[63-i] += 1
			v &^= 1 << i
		}
	}
	return
}

func roundOdd(v, precisionBitsNum uint64) uint64 {
	precisionBitsMask := uint64(1<<precisionBitsNum) - 1
	flow := v & precisionBitsMask
	cond := uint64(1 << (precisionBitsNum - 1))
	var endBit uint64 = 0
	if flow > cond || flow == cond && v&(1<<precisionBitsNum) > 0 { //向偶数舍入,统计误差最小
		endBit = 1
	}
	return (v &^ precisionBitsMask) + (endBit << precisionBitsNum)
}
