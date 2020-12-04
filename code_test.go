package FixedNumber

import (
	"testing"
)

func TestFixed54_10_Add(t *testing.T) {
	var in  = 45678.664
	fixed := Float64ToFixed54_10(in)
	if fixed.Add(fixed).Add(fixed).Add(fixed).Add(fixed) != Float64ToFixed54_10(in * 5){
		t.Fail()
	}

	if fixed.Del(fixed) != Float64ToFixed54_10(in - in){
		t.Fail()
	}
}
