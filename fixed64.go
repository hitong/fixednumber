package fixed

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"strconv"
	"strings"
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

//Precision is the bit Precision of a fixed-point number
var Precision int64

var onceSet = &sync.Once{}
var precisionBitsNum = Precision
var decimalBitsMask uint64 = 1<<precisionBitsNum - 1

//MaxFixed64 is the maximum number of fixed points,and its value is 2**(63 - Precision) + 2 ** Precision - 2
const MaxFixed64 = Fixed64((1 << 63) - 1)

//SmallestFixed64 is the minimum value of a fixed-point number, and its value is 1/2**Precision - 2**(63-precisionBitsNum)
const SmallestFixed64 = ^Fixed64(0)

//Fixed64Zero is the zero value of fixed number
const Fixed64Zero = Fixed64(0)

//PrecisionNumber is the minimum precision value of a fixed-point number, and its value is 1 / 2**Precision
const PrecisionNumber = Fixed64(1)

// Fixed64 uses uint64 type to facilitate bit conversion.Fixed64 can use +-0.
// Range: [-(2**(63 - Precision) + 2 ** Precision - 2), 2**(63 - Precision) + 2 ** Precision - 2].
type Fixed64 uint64

//SetPrecisionOnce can only be successfully set once
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

// Float64ToFixed64 can convert the normalized number and normalized number conforming to IEEE 754
// double precision floating-point number standard to fixed-point number. When the function processes
// Nan and Inf, there may be a panic
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
				fixedP = (m & pMask) << (-1 * pBitsFlowNum)
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

// Str2Fixed64 produces the result by converting the string value to a two-part float,
// and then performing a combination of division and addition.The efficiency of the function may not be high,
// there should be a better implementation.
//
// Note: the function only supports processing data in (%d+.%d+)
func Str2Fixed64(value string)(Fixed64,error){
	s := uint64(0)
	if strings.Contains(value,"-"){
		value = value[1:]
		s = mask64S
	}
	values := strings.Split(value,".")

	if len(values) == 0 || len(values) > 2{
		return 0,errors.New(fmt.Sprintf("Fixed64: %s format is err ",value))
	}
	if len(values) == 1{
		values = append(values, "0")
	}

	if values[0] == ""{
		values[0] = "0"
	} else if values[1] == ""{
		values[1] = "0"
	}

	di,e1 := parseStringToFixed64(values[0])
	de,e2 := parseStringToFixed64(values[1])
	if e1 != nil || e2 != nil{
		return 0,errors.New(fmt.Sprintf("Fixed64：%s con not parse to Fixed64 ",value))
	}

	dee := Float64ToFixed64(math.Pow10(len(values[1])))

	return Fixed64(uint64(di.Add(de.Div(dee))) | s),nil
}

func parseStringToFixed64(v string)(Fixed64, error){
	if f,err := strconv.ParseFloat(v,64);err != nil{
		return 0,err
	} else {
		return Float64ToFixed64(f),nil
	}
}


//SafeFloat64ToFixed64 checks NaN and Inf before floating point conversion, which guarantees some accuracy
func SafeFloat64ToFixed64(value float64) (Fixed64, error) {
	if math.IsNaN(value) {
		return 0, errors.New("Float64 value is NaN ")
	}
	if math.IsInf(value, 0) {
		return 0, errors.New("Float64 value is Inf ")
	}

	return Float64ToFixed64(value), nil
}

func (fixed Fixed64) Add(oth Fixed64) Fixed64 {
	var fS = uint64(fixed) & mask64S
	var oS = uint64(oth) & mask64S
	fixed &^= mask64S
	oth &^= mask64S
	if fS == oS {
		hi,lo := add64(uint64(fixed),uint64(oth))
		if hi >> 31 > 0{
			panic("Fixed64: Add Overflow " + fixed.ToBase10N(18) + " " + oth.ToBase10N(18))
		}
		return Fixed64(hi << 32 | lo | fS)
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

func add64(x,y uint64)(hi,lo uint64){
	const mask32 = 1 << 32 - 1
	var x0 = x >> 32
	var x1 = x & mask32
	var y0 = y >> 32
	var y1 = y & mask32
	lo = x1 + y1
	var u = lo >> 32
	lo = lo & mask32
	hi = x0 + y0 + u
	return
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
	if hi >> precisionBitsNum > 0{
		panic("Fixed64: Number OverFlow")
	}

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
		return (int64((fixed&^mask64S)>>precisionBitsNum) + roundUp) * -1
	}

	return int64((fixed&^mask64S)>>precisionBitsNum) + roundUp
}

//5 decimal places are retained by default
func (fixed Fixed64) ToBase10() string {
	return fixed.ToBase10N(5)
}

//the maximum support is 18 decimals
func (fixed Fixed64) ToBase10N(n uint) string {
	if n > 18 {
		panic("Fixed64.ToBase10N: Not Support n > 19")
	}
	floatSlice := make([]byte, 0, 10)
	if fixed&mask64S > 0 {
		floatSlice = append(floatSlice, '-')
	}

	if n == 0 {
		return string(append(floatSlice, insertToFloatSliceBase10(uint64(fixed.Abs().Round()), 0)...))
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
		partD[len(partD)-1]++
		partD, upToTop = carryUpBase10(partD)
	}

	if upToTop {
		floatSlice = append(floatSlice, '1')
	}
	floatSlice = append(floatSlice, partD...)
	floatSlice = append(floatSlice, '.')
	floatSlice = append(floatSlice, partP...)
	return string(floatSlice)
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
			r[63-i]++
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
	//if p == nil || len(p) == 0 {
	//	return []byte{'0'}
	//}

	if p[len(p)-1] > '4' {
		p[len(p)-1] = 0
		//if len(p) == 1 {
		//	p[0] = '9' + 1
		//	return p
		//}
		p[len(p)-2]++
	}

	return p
}

func carryUpBase10(p []byte) ([]byte, bool) {
	i := len(p) - 1
	for ; i >= 1; i-- {
		if p[i] > '9' {
			p[i] = '0'
			p[i-1]++
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
