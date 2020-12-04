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
	var fS uint64 = uint64(fixed) >> 63
	var oS uint64 = uint64(oth) >> 63
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
		if fixed > oth {
			return Fixed54_10(uint64(int64(fixed)-int64(oth)) | fS<<63)
		} else {
			return Fixed54_10(uint64(int64(oth)-int64(fixed)) | oS<<63)
		}
	}
}

func (fixed Fixed54_10) Del(oth Fixed54_10) Fixed54_10 {
	return fixed.Add(-oth)
}

func (fixed Fixed54_10) Mul(oth Fixed54_10) Fixed54_10 {
	return fixed.Add(-oth)
}

func (fixed Fixed54_10) Div(oth Fixed54_10) Fixed54_10 {
	return fixed.Add(-oth)
}


func (fixed Fixed54_10)String()string{
	var signBit int64 = 1
	if fixed &mask64S == mask64S {
		signBit = -1
	}

	return fmt.Sprintf("%d.%d",int64((fixed &^mask64S) >> 10) * signBit, (1 << 10 - 1) & fixed)
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