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

var pow10tab = [...]uint64{
	1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
}

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

func Uint32Bits(v uint32) (r [32]int) {
	for i := 0; i < 32; i++ {
		if v == 0 {
			return
		}
		if (v & (1 << i)) > 0 {
			r[31-i] = 1
			v = v &^ (1 << i)
		}
	}

	return
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

//
//func (fixed Fixed64) Div(oth Fixed64) Fixed64 {
//	var fS = uint64(fixed) >> 63
//	var oS = uint64(oth) >> 63
//	fixed &^= mask64S
//	oth &^= mask64S
//
//	fixedInt := uint64(fixed >> PrecisionBitsNum) * pow10tab[6] + uint64(fixed & DecimalBitsMask) * pow10tab[3]
//	othInt :=  uint64(oth >> PrecisionBitsNum) * pow10tab[3]+ uint64(oth & DecimalBitsMask)
//
//	endIntWithDec := fixedInt / othInt
//	endIntNoDec := endIntWithDec /  pow10tab[Precision]
//	endDec := (endIntWithDec - endIntNoDec  * pow10tab[Precision]) % (pow10tab[Precision] + 1)
//
//	return Fixed64((fS &^ oS << 63) | (endIntNoDec << PrecisionBitsNum) | endDec)
//}

//
//func Float64ToFixed64(value float64) Fixed64 {
//	var v uint64 = 0
//	var valueBits = math.Float64bits(value)
//	v |= valueBits & mask64S
//	if value < 0{
//		value *= -1
//	}
//	num := int64(value * 1000)
//	d := num / 1000
//	t := num % 1000
//	//p := t * float64(pow10tab[Precision])
//	//p := int64(math.Floor(t * float64(pow10tab[Precision]) + 0.5))
//	v |= uint64(d) << PrecisionBitsNum
//	v |= uint64(t)
//	return Fixed64(v)
//}
//
//func (fixed Fixed64) Add(oth Fixed64) Fixed64 {
//	var fS = uint64(fixed) >> 63
//	var oS = uint64(oth) >> 63
//	fixed &^= mask64S
//	oth &^= mask64S
//	fixedH :=fixed &^ DecimalBitsMask
//	othH :=oth &^ DecimalBitsMask
//	fixedL := fixed & DecimalBitsMask
//	othL := oth & DecimalBitsMask
//	if fS == oS {
//		h := fixedH + othH
//		l := othL + fixedL
//		lUp := l / Fixed64(pow10tab[Precision])
//		h += lUp << PrecisionBitsNum
//		l -= lUp * Fixed64(pow10tab[Precision])
//		return Fixed64(uint64(h + l) | fS<<63)
//	} else {
//		h := int64(fixedH) - int64(othH)
//		l := int64(fixedL) - int64(othL)
//		endS := 0
//		if (fS == 1 && fixed > oth ) || (oS == 1 && fixed < oth) {
//			endS = 1
//		}
//
//		var end = uint64(endS) << 63
//		if h > 0 {
//			if l < 0{
//				h -= 1 << PrecisionBitsNum
//				l += int64(pow10tab[Precision])
//			}
//			end |= uint64(h + l)
//		} else if h < 0 {
//			if l > 0{
//				h += 1 << PrecisionBitsNum
//				l -= int64(pow10tab[Precision])
//			}
//
//			end |= uint64((h + l) * -1)
//		} else {
//			if l < 0{
//				l *= -1
//			}
//			end |= uint64(l)
//		}
//
//		return Fixed64(end)
//	}
//}
//
//func (fixed Fixed64) Mul(oth Fixed64) Fixed64 {
//	var fS = fixed & mask64S
//	var oS = oth & mask64S
//	fixed &^= mask64S
//	oth &^= mask64S
//
//	fixedInt := uint64(fixed>>PrecisionBitsNum)*pow10tab[Precision] + uint64(fixed&DecimalBitsMask)
//	othInt := uint64(oth>>PrecisionBitsNum)*pow10tab[Precision] + uint64(oth&DecimalBitsMask)
//
//	endIntWithDec := fixedInt * othInt
//	endIntNoDec := endIntWithDec / pow10tab[Precision<<1]
//	endDec := (endIntWithDec - endIntNoDec*pow10tab[Precision<<1]) % (pow10tab[Precision<<1] + 1) / pow10tab[Precision]
//
//	return Fixed64(uint64(fS^oS) | (endIntNoDec << PrecisionBitsNum) | endDec)
//}
