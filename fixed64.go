package fixed

import (
	"errors"
	"math"
	"math/bits"
	"sync"
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

var Precision int64 = 0 //Fixed point precision:1/2**Precision

var onceSet = &sync.Once{}
var precisionBitsNum = Precision
var decimalBitsMask uint64 = 1<<precisionBitsNum - 1

const MaxFixed64 = Fixed64((1 << 63) - 1) //Maximum number of fixed points: 2**(63-precisionBitsNum) - 1/2**Precision
const SmallestFixed64 = ^Fixed64(0)       //Smallest fixed point:  1/2**Precision - 2**(63-precisionBitsNum)
const Fixed64Zero = Fixed64(0)            //Zero
const PrecisionNumber = Fixed64(1)

type Fixed64 uint64

//Precision can only be set once
func SetPrecisionOnce(precision uint64) {
	onceSet.Do(func() {
		if precision > 62 {
			panic("Precision overflow")
		}
		Precision = int64(precision)
		precisionBitsNum = Precision
		decimalBitsMask = 1<<precisionBitsNum - 1
	})
}

//Convert normalizing float point number to fixed point number
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
				panic("Fixed number: part digital overflow")
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
		if pFlowBitsNum <= 0 {
			//pFlowBitsNum *= -1
			fixedP = m << (pFlowBitsNum * -1)
		} else {
			fixedP = roundOdd(m, uint64(pFlowBitsNum)) >> pFlowBitsNum
		}
	}

	result |= s
	result |= fixedD
	result |= fixedP
	return Fixed64(result)
}

//Convert normalizing float point number to fixed point number with error back
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
	lo = roundOdd(lo, uint64(precisionBitsNum)) >> precisionBitsNum

	hi = (hi & decimalBitsMask) << (64 - precisionBitsNum)
	if hi == 0 && lo == 0 {
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
	if quo == 0 {
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

//Converts Fixed64 To float64
func (fixed Fixed64) Float64() float64 {
	number := uint64(fixed &^ mask64S)
	idx := int64(bits.Len64(number))
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
	if fixed&mask64S > 0 {
		return int64((fixed&^mask64S)>>precisionBitsNum) * -1
	}

	return int64((fixed &^ mask64S) >> precisionBitsNum)
}

func (fixed Fixed64) Round() int64 {
	decimal := uint64(fixed) & decimalBitsMask
	var roundUp int64
	if precisionBitsNum > 0 && decimal >= (1<<(precisionBitsNum-1)) {
		roundUp = 1
	}
	if fixed&mask64S > 0 {
		return int64((fixed&^mask64S)>>precisionBitsNum)*-1 + roundUp
	}

	return int64((fixed&^mask64S)>>precisionBitsNum) + roundUp
}

//5 decimal places are retained by default
func (fixed Fixed64) ToBase10() []byte {
	return fixed.ToBase10N(5)
}

func (fixed Fixed64)ToBase10s(n uint) string{
	return string(fixed.ToBase10N(n))
}

//the maximum support is 18 decimals
func (fixed Fixed64) ToBase10N(n uint) []byte {
	if n > 18 {
		panic("Fixed64.ToBase10N: Not Support n > 19")
	}
	floatSlice := make([]byte, 0, 10)
	if fixed&mask64S > 0 {
		floatSlice = append(floatSlice, '-')
	}

	if n == 0 {
		return append(floatSlice, insertToFloatSliceBase10(uint64(fixed.Abs().Round()), 0)...)
	}

	n = n + 1 // To round to the end
	number := uint64(fixed &^ mask64S)
	d := int64(number >> precisionBitsNum)
	p := number & decimalBitsMask

	hi, lo := bits.Mul64(p, uint64(math.Pow10(int(n)))) //todo optimize data processing methods
	quo, _ := bits.Div64(hi, lo, 1<<precisionBitsNum)

	partD := insertToFloatSliceBase10(uint64(d), 0)
	partP, upToTop := carryUpBase10(roundEndBase10(insertToFloatSliceBase10(quo, int(n)))[:n-1])
	if upToTop {
		partD[len(partD)-1] += 1
		partD, upToTop = carryUpBase10(partD)
	}

	if upToTop {
		floatSlice = append(floatSlice, '1')
	}
	floatSlice = append(floatSlice, partD...)
	floatSlice = append(floatSlice, '.')
	floatSlice = append(floatSlice, partP...)
	return floatSlice
}

func insertToFloatSliceBase10(v uint64, n int) []byte {
	const zeroStr = "00000000000000000000" //len 20

	var ret []byte
	if n > 0 {
		ret = make([]byte, 0, n)
	}
	for v > 0 {
		ret = append(ret, byte(v%10)+'0')
		v /= 10
	}
	if n > len(ret) {
		ret = append(ret, zeroStr[:n-len(ret)]...)
	}
	if len(ret) > 0 {
		for i := 0; i < len(ret)>>1; i++ {
			ret[i], ret[len(ret)-1-i] = ret[len(ret)-1-i], ret[i]
		}
	} else {
		ret = append(ret, '0')
	}

	return ret
}

func (fixed Fixed64) String() string {
	//return strconv.FormatFloat(fixed.Float64(), 'f', -1, 64)
	//s := strings.TrimRightFunc(fmt.Sprintf("%f", fixed.Float64()), func(r rune) bool {
	//	if r == '0' || r == '.' {
	//		return true
	//	}
	//	return false
	//})
	//
	//if s == "" {
	//	s = "0"
	//}
	return string(fixed.ToBase10())
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

// round to even
func roundOdd(v, precisionBitsNum uint64) uint64 {
	precisionBitsMask := uint64(1<<precisionBitsNum) - 1
	flow := v & precisionBitsMask
	cond := uint64(1 << (precisionBitsNum - 1))
	var endBit uint64
	if flow > cond || flow == cond && v&(1<<precisionBitsNum) > 0 {
		endBit = 1
	}

	return (v &^ precisionBitsMask) + (endBit << precisionBitsNum)
}

func roundEndBase10(p []byte) []byte {
	if p == nil || len(p) == 0 {
		return []byte{'0'}
	}

	if p[len(p)-1] > '4' {
		p[len(p)-1] = 0
		if len(p) == 1 {
			p[0] = '9' + 1
			return p
		}
		p[len(p)-2] += 1
	}

	return p
}

func carryUpBase10(p []byte) ([]byte, bool) {
	i := len(p) - 1
	for ; i >= 1; i-- {
		if p[i] > '9' {
			p[i] = '0'
			p[i-1] += 1
		} else {
			return p, false
		}
	}
	upToTop := p[0] == '9'+1
	if upToTop {
		p[0] = '0'
	}

	return p, upToTop
}
