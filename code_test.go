package FixedNumber

import (
	"fmt"
	"testing"
)

func TestFixed54_10_Add(t *testing.T) {
	Float64ToFixed54_10(-1.0).Add(Float64ToFixed54_10(1))
	var v1  = [...]float64{0,-1,5,9,456,789,12.145,45.11,-95.154,456.01}
	var v2  = [...]float64{0,1,165.6,48.1,5.1,5654.1,-4.1,-69.6,.26,0}
	for i := 0; i < len(v1); i++{
		fmt.Println("add " ,Float64ToFixed54_10(v1[i]) , " + ", Float64ToFixed54_10(v2[i]) ," = ", Float64ToFixed54_10(v1[i]).Add(Float64ToFixed54_10(v2[i])))
	}
}

func TestFixed54_10_Del(t *testing.T) {
	var v1  = [...]float64{0,-1,5,9,456,789,12.145,45.11,-95.154,456.01}
	var v2  = [...]float64{0,1,165.6,48.1,5.1,5654.1,-4.1,-69.6,.26,0}
	for i := 0; i < len(v1); i++{
		fmt.Println("del " ,Float64ToFixed54_10(v1[i]) , " + ", Float64ToFixed54_10(v2[i]) ," = ", Float64ToFixed54_10(v1[i]).Del(Float64ToFixed54_10(v2[i])))
	}
}

func TestFixed54_10_Mul(t *testing.T) {
	var v1  = [...]float64{0,-1,5,9,456,789,12.145,45.11,-95.154,456.01}
	var v2  = [...]float64{0,1,165.6,48.1,5.1,5654.1,-4.1,-69.6,.26,0}
	for i := 0; i < len(v1); i++{
		fmt.Println("mul " ,Float64ToFixed54_10(v1[i]) , " + ", Float64ToFixed54_10(v2[i]) ," = ", Float64ToFixed54_10(v1[i]).Mul(Float64ToFixed54_10(v2[i])))
	}
}

func TestFixed54_10_Div(t *testing.T) {
	var v1  = [...]float64{0,-1,5,9,456,789,12.145,45.11,-95.154,456.01}
	var v2  = [...]float64{2,1,165.6,48.1,5.1,5654.1,-4.1,-69.6,.26,550}
	for i := 0; i < len(v1); i++{
		fmt.Println("div " ,Float64ToFixed54_10(v1[i]) , " + ", Float64ToFixed54_10(v2[i]) ," = ", Float64ToFixed54_10(v1[i]).Div(Float64ToFixed54_10(v2[i])))
	}
}
