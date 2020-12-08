package fixed

import (
	"math"
	"math/bits"
	"strconv"
	"unsafe"
)

const (
	mask64S  = 1 << 63
	SizeM    = 52
	mask64E  = (1<<11 - 1) << SizeM
	mask64M  = (1 << SizeM) - 1
	shift64E = SizeM
	bias64   = (1 << (11 - 1)) - 1
)

const Precision = 30
const PrecisionBitsNum = Precision
const DecimalBitsMask = 1<<PrecisionBitsNum - 1

type Fixed64 uint64

//规范数处理
func Float64ToFixed64(value float64) Fixed64 {
	var valueBits = math.Float64bits(value)
	s := valueBits & mask64S
	e := valueBits & mask64E
	m := valueBits&mask64M | (1 << SizeM)
	realE := int64(e>>shift64E) - bias64 //todo 限制范围位 [-PrecisionBitsNum - 1,+52]
	fixedD := uint64(0)
	var result uint64
	var fixedP uint64 = 0
	//todo 定点数溢出, Panic?,Error?
	if realE >= 0 {
		fixedD = m >> (SizeM - realE) << PrecisionBitsNum
		fixedP = ((1<<PrecisionBitsNum - 1) << (SizeM - realE - PrecisionBitsNum) & m) >> (SizeM - realE - PrecisionBitsNum)
	} else {
		fixedP = m >> (realE * -1) >> (SizeM - PrecisionBitsNum)
	}

	result |= s
	result |= fixedD
	result |= fixedP
	return Fixed64(result)
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
		} else {
			return Fixed64(uint64(oth-fixed) | oS)
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
	lo >>= PrecisionBitsNum //todo 舍入规则定义
	hi = (hi & DecimalBitsMask) << (64 - PrecisionBitsNum)

	return Fixed64(uint64(fS^oS) | hi | lo)
}

func (fixed Fixed64) Div(oth Fixed64) Fixed64 {
	var fS = fixed & mask64S
	var oS = oth & mask64S
	fixed &^= mask64S
	oth &^= mask64S

	quo, _ := bits.Div64(uint64(fixed>>(64-PrecisionBitsNum)), uint64(fixed<<PrecisionBitsNum), uint64(oth))
	//lo >>= PrecisionBitsNum //todo 舍入规则定义
	//hi = (hi & DecimalBitsMask) << (64 - PrecisionBitsNum)

	return Fixed64(uint64(fS^oS) | quo)
}

func (fixed Fixed64) String() string {
	return strconv.FormatFloat(Fixed64ToFloat64(fixed), 'f', -1, 64)
}

func Fixed64ToFloat64(fixed Fixed64) float64 {
	number := uint64(fixed &^ mask64S)
	idx := getLastBitIdx(number)
	if idx != 0 {
		e := idx - PrecisionBitsNum - 1
		number = ((1 << idx) - 1) & number
		if idx > SizeM {
			number >>= idx - SizeM - 1
		} else {
			number <<= SizeM - idx + 1
		}
		number = number &^ mask64E
		number |= uint64(fixed & mask64S)
		number |= uint64((e + bias64) << SizeM)
	}

	return *(*float64)(unsafe.Pointer(&number))
}

func getLastBitIdx(v uint64) int {
	idx := 0
	for v != 0 {
		v >>= 1
		idx++
	}
	return idx
}

func Uint64Bits(v uint64) (r [64]int) {
	for i := 0; i < 64; i++ {
		if (v & (1 << i)) > 0 {
			r[63-i] = 1
			v &^= 1 << i
		}
	}
	return
}
