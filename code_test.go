package fixed

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
)

func TestDiv64(t *testing.T)  {
	fmt.Println(string(Float64ToFixed64(1/3.0).ToBase10()))
}

func TestBasic(t *testing.T) {
	fmt.Println(Float64ToFixed64(489.644).Mul(Float64ToFixed64(math.Float64frombits(1 << 63))))
	fmt.Println(Float64ToFixed64(456.6).Mul(1 << 63))
	fmt.Println(Float64ToFixed64(1).Div(Float64ToFixed64(3)))
	f0 := Float64ToFixed64(123.456)
	f1 := Float64ToFixed64(123.456)



	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}

	if f0.Int64() != 123 {
		t.Error("should be equal", f0.Int64(), 123)
	}

	if strings.TrimRightFunc(f0.String(), func(r rune) bool {
		if r == '0'{
			return true
		}
		return false
	}) != "123.456" {
		t.Error("should be equal", f0.String(), "123.456")
	}

	f0 = Float64ToFixed64(1)
	f1 = Float64ToFixed64(.5).Add(Float64ToFixed64(.5))
	f2 := Float64ToFixed64(.3).Add(Float64ToFixed64(.3)).Add(Float64ToFixed64(.4))

	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}
	if !f0.Equal(f2) {
		t.Error("should be equal", f0, f2)
	}

	f0 = Float64ToFixed64(.999)
	if f0.String() != "0.999" {
		t.Error("should be equal", f0, "0.999")
	}

	//f0 = Float64ToFixed64(1).Sub(Float64ToFixed64(1/3.0)).Sub(Float64ToFixed64(1/3.0)).Sub(Float64ToFixed64(1/3.0))
	//if f0.String() != "0" {
	//	t.Error("should be equal", f0, "0")
	//}
}

func TestUint64Bits(t *testing.T) {
	fmt.Println(Uint64Bits(math.Float64bits(rand.Float64())))
}

func TestSafeFloat64ToFixed64(t *testing.T) {
	if _, err := SafeFloat64ToFixed64(math.Float64frombits(^uint64(0))); err != nil {
	}
}

func TestFixed64_Add(t *testing.T) {
	Float64ToFixed64(-1.0).Add(Float64ToFixed64(1))
	var v1 = [...]float64{0, -1, 5, 9, 456, 789, 12.145, 45.11, -95.154, 456.01}
	var v2 = [...]float64{0, 1, 165.6, 48.1, 5.1, 5654.1, -4.1, -69.6, .26, 0}
	for i := 0; i < len(v1); i++ {
		fmt.Println("add ", Float64ToFixed64(v1[i]), " + ", Float64ToFixed64(v2[i]), " = ", Float64ToFixed64(v1[i]).Add(Float64ToFixed64(v2[i])))
	}
}

func TestFixed64_Del(t *testing.T) {
	var v1 = [...]float64{0, -1, 5, 9, 456, 789, 12.145, 45.11, -95.154, 456.01}
	var v2 = [...]float64{0, 1, 165.6, 48.1, 5.1, 5654.1, -4.1, -69.6, .26, 0}
	for i := 0; i < len(v1); i++ {
		fmt.Println("del ", Float64ToFixed64(v1[i]), " - ", Float64ToFixed64(v2[i]), " = ", Float64ToFixed64(v1[i]).Sub(Float64ToFixed64(v2[i])))
	}
}

func TestFixed64_Mul(t *testing.T) {
	var v1 = [...]float64{779.281257,0, -1, 5, 9, 456, 789, 12.145, 45.11, -95.154, 456.01}
	var v2 = [...]float64{828.6875,0, 1, 165.6, 48.1, 5.1, 5654.1, -4.1, -69.6, .26, 0}
	for i := 0; i < len(v1); i++ {
		fmt.Println("mul ", Float64ToFixed64(v1[i]), " * ", Float64ToFixed64(v2[i]), " = ", Float64ToFixed64(v1[i]).Mul(Float64ToFixed64(v2[i])))
	}
}

func TestFixed64_Div(t *testing.T) {
	var v1 = [...]float64{0, -1, 5, 9, 456, 789, 12.145, 45.11, -95.154, 456.01}
	var v2 = [...]float64{2, 1, 165.6, 48.1, 5.1, 5654.1, -4.1, -69.6, .26, 550}

	for i := 0; i < len(v1); i++ {
		fmt.Println("div ", Float64ToFixed64(v1[i]), " / ", Float64ToFixed64(v2[i]), " = ", Float64ToFixed64(v1[i]).Div(Float64ToFixed64(v2[i])))
	}
}

func TestFixed64(t *testing.T) {
	v := 630.506836 *237.22168
	fmt.Println(v)// 149569.890625 149570.131554
	//fmt.Println(fmt.Sprintf("%.3f",Fixed64ToFloat64(Float64ToFixed64(45687.2456))))
	fmt.Println(fmt.Sprintf("%.3f", Float64ToFixed64(489.15641).Mul(Float64ToFixed64(3.64)).Float64()))

	fmt.Println(Float64ToFixed64(0.14564).Sub(Float64ToFixed64(0.14564)))
	fmt.Println(Float64ToFixed64(0.14564).Mul(Float64ToFixed64(0.14564)))
	fmt.Println(Float64ToFixed64(45687.2456).Mul(Float64ToFixed64(0)))
	//fmt.Println(Float64ToFixed64(4568.015).Div(Float64ToFixed64(12)))
	fmt.Println(Float64ToFixed64(0.2455))
	//fmt.Println(Float64ToFixed64(9).Mul(Float64ToFixed64(4568.64485)))
	//fmt.Println(Float64ToFixed64(458.654).Add(Float64ToFixed64(45648.656)))
	//fmt.Println(Float64ToFixed64(3111984.465489))
	//fmt.Println(Float64ToFixed64(3111984.465489).Add(Float64ToFixed64(3111984.465489)))
}

func TestBuildData(t *testing.T) {
	BuildData("AddData.txt", "Add", 1000)
	BuildData("SubData.txt", "Sub", 1000)
	BuildData("MulData.txt", "Mul", 1000)
	BuildData("DivData.txt", "Div", 1000)
}

func TestReadData(t *testing.T) {
	paths := []string{"AddData.txt", "SubData.txt", "DivData.txt", "MulData.txt"}
	for _, path := range paths {
		if op, v1, v2, v3, err := ReadData(path); err != nil {
			panic(err)
		} else {
			for i := 0; i < len(op); i++ {
				//fmt.Println(fmt.Sprintf("%s:%f %f %f",op[i],v1[i],v2[i],v3[i]))
				fix1 := Float64ToFixed64(v1[i])
				fix2 := Float64ToFixed64(v2[i])
				fix3 := Float64ToFixed64(v3[i])
				switch op[i] {
				case "Add":
					if int64(fix1.Add(fix2).Float64())*1000 != int64((fix3.Float64())*1000) {
						//fmt.Println(int64(Fixed64ToFloat64(fix1.Add(fix2)) * 1000), int64(Fixed64ToFloat64(fix3) * 1000))
						println("Add:" + fmt.Sprintf("%.10f %s %s %s %f", fix1.Add(fix2).Float64()-v3[i], fix1, fix2, fix1.Add(fix2), v3[i]))
					}
				case "Sub":
					if fix1.Sub(fix2) != fix3 {
						println("Sub:" + fmt.Sprintf("%f %s %s %s %f", fix1.Sub(fix2).Float64()-v3[i], fix1, fix2, fix1.Sub(fix2), v3[i]))
					}
				case "Div":
					if fix1.Div(fix2) != fix3 {
						println("Div:" + fmt.Sprintf("%f %s %s %s %f", fix1.Div(fix2).Float64()-v3[i], fix1, fix2, fix1.Div(fix2), v3[i]))
					}
				case "Mul":
					if fix1.Mul(fix2) != fix3 {
						println("Mul:" + fmt.Sprintf("%f %s %s %s %f", fix1.Mul(fix2).Float64()-v3[i], fix1, fix2, fix1.Mul(fix2), v3[i]))
					}
				}
			}
		}
	}
}
