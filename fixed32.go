package FixedNumber

import (
	"fmt"
	"math"
)

//type IFixedNumber interface {
//	Add(IFixedNumber) IFixedNumber
//	Del(IFixedNumber) IFixedNumber
//	Mul(IFixedNumber) IFixedNumber
//	Div(IFixedNumber) IFixedNumber
//}

const (
	mask64E = (1 << 8) - 1
	mask64B = (1 << 23) - 1
	mask64S = 1 << 63
	shift64E = 23
)

type int128 struct {
	high int64
	low int64
}

var pow10tab = [...]uint64{
	1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
}

var zeroStr = [...]string{
	"00","00","0","",
}

type Fixed54_10 uint64

func Float64ToFixed54_10(value float64) Fixed54_10 {
	var v uint64 = 0
	var valueBits = math.Float64bits(value)
	v |= valueBits & mask64S
	if value < 0{
		value *= -1
	}
	d := int64(value)
	t := value - float64(d)
	p := int64(math.Floor(t * 1000 + 0.5))
	v |= uint64(d) << 10
	v |= uint64(p)
	return Fixed54_10(v)
}

func (fixed Fixed54_10) Add(oth Fixed54_10) Fixed54_10 {
	var fS = uint64(fixed) >> 63
	var oS = uint64(oth) >> 63
	fixed &^= mask64S
	oth &^= mask64S
	fixedH :=fixed &^ (1 << 10 - 1)
	othH :=oth &^ (1 << 10 - 1)
	fixedL := fixed & (1 << 10 - 1)
	othL := oth & (1 << 10 - 1)
	if fS == oS {
		h := fixedH + othH
		l := othL + fixedL
		lUp := l / 1000
		h += lUp << 10
		l -= lUp * 1000
		return Fixed54_10(uint64(h + l) | fS<<63)
	} else {
		h := int64(fixedH) - int64(othH)
		l := int64(fixedL) - int64(othL)
		endS := 0
		if (fS == 1 && fixed > oth ) || (oS == 1 && fixed < oth) {
			endS = 1
		}

		var end = uint64(endS) << 63
		if h > 0 {
			if l < 0{
				h -= 1024
				l += 1000
			}
			end |= uint64(h + l)
		} else if h < 0 {
			if l > 0{
				h += 1024
				l -= 1000
			}

			end |= uint64((h + l) * -1)
		} else {
			if l < 0{
				l *= -1
			}
			end |= uint64(l)
		}

		return Fixed54_10(end)
	}
}

func (fixed Fixed54_10) Del(oth Fixed54_10) Fixed54_10 {
	return fixed.Add(oth ^ mask64S)
}

func (fixed Fixed54_10) Mul(oth Fixed54_10) Fixed54_10 {
	var fS = uint64(fixed) >> 63
	var oS = uint64(oth) >> 63
	fixed &^= mask64S
	oth &^= mask64S

	fixedInt := uint64(fixed >> 10) * 1000 + uint64(fixed & (1 << 10 - 1))
	othInt :=  uint64(oth >> 10) * 1000+ uint64(oth & (1 << 10 - 1))

	endIntWithDec := fixedInt * othInt
	endIntNoDec := endIntWithDec /  pow10tab[6]
	endDec := (endIntWithDec - endIntNoDec  * pow10tab[6]) % (pow10tab[6] + 1) / pow10tab[3]

	return Fixed54_10((fS &^ oS << 63) | (endIntNoDec << 10) | endDec)
}

func (fixed Fixed54_10) Div(oth Fixed54_10) Fixed54_10 {
	var fS = uint64(fixed) >> 63
	var oS = uint64(oth) >> 63
	fixed &^= mask64S
	oth &^= mask64S

	fixedInt := uint64(fixed >> 10) * pow10tab[6] + uint64(fixed & (1 << 10 - 1)) * pow10tab[3]
	othInt :=  uint64(oth >> 10) * pow10tab[3]+ uint64(oth & (1 << 10 - 1))

	endIntWithDec := fixedInt / othInt
	endIntNoDec := endIntWithDec /  pow10tab[3]
	endDec := (endIntWithDec - endIntNoDec  * pow10tab[3]) % (pow10tab[3] + 1)

	return Fixed54_10((fS &^ oS << 63) | (endIntNoDec << 10) | endDec)
}


func (fixed Fixed54_10)String()string{
	var signBit int64 = 1
	if fixed &mask64S == mask64S {
		signBit = -1
	}
	if signBit == -1{
		return fmt.Sprintf("-%d.%s%d",int64((fixed &^mask64S) >> 10) ,getZero(fixed), (1 << 10 - 1) & fixed)
	} else {
		return fmt.Sprintf("%d.%s%d",int64((fixed &^mask64S) >> 10) ,getZero(fixed), (1 << 10 - 1) & fixed)
	}
}

func getZero(deci Fixed54_10)string{
	var i = 0
	v := uint64(deci & (1 << 10 - 1))
	for ; v > pow10tab[i];i++{}

	return zeroStr[i]
}

func Uint32Bits(v uint32) (r [32]int){
	for i := 0 ; i < 32 ;i++{
		if v == 0{
			return
		}
		if (v & (1 << i)) > 0 {
			r[31 - i] = 1
			v = v &^ (1 << i)
		}
	}

	return
}

func Uint64Bits(v uint64)(r [64]int){
	for i := 0; i < 64; i++{
		if (v & (1 << i )) > 0{
			r[63 - i] = 1
			v &^= 1 << i
		}
	}
	return
}