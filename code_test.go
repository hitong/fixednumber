package fixed

import (
	"fmt"
	"testing"
)

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
	var v1 = [...]float64{0, -1, 5, 9, 456, 789, 12.145, 45.11, -95.154, 456.01}
	var v2 = [...]float64{0, 1, 165.6, 48.1, 5.1, 5654.1, -4.1, -69.6, .26, 0}
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
	//fmt.Println(fmt.Sprintf("%.3f",Fixed64ToFloat64(Float64ToFixed64(45687.2456))))
	fmt.Println(fmt.Sprintf("%.3f", Fixed64ToFloat64(Float64ToFixed64(489.15641).Mul(Float64ToFixed64(3.64)))))

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
					if int64(Fixed64ToFloat64(fix1.Add(fix2))*1000) != int64(Fixed64ToFloat64(fix3)*1000) {
						//fmt.Println(int64(Fixed64ToFloat64(fix1.Add(fix2)) * 1000), int64(Fixed64ToFloat64(fix3) * 1000))
						println("Add:" + fmt.Sprintf("%.10f %s %s %s %s", Fixed64ToFloat64(fix1.Add(fix2))-Fixed64ToFloat64(fix3), fix1, fix2, fix1.Add(fix2), fix3))
					}
				case "Sub":
					if fix1.Sub(fix2) != fix3 {
						println("Sub:" + fmt.Sprintf("%f %s %s %s %s", Fixed64ToFloat64(fix1.Sub(fix2))-Fixed64ToFloat64(fix3), fix1, fix2, fix1.Sub(fix2), fix3))
					}
				case "Div":
					if fix1.Div(fix2) != fix3 {
						println("Div:" + fmt.Sprintf("%f %s %s %s %s", Fixed64ToFloat64(fix1.Div(fix2))-Fixed64ToFloat64(fix3), fix1, fix2, fix1.Div(fix2), fix3))
					}
				case "Mul":
					if fix1.Mul(fix2) != fix3 {
						println("Mul:" + fmt.Sprintf("%f %s %s %s %s", Fixed64ToFloat64(fix1.Mul(fix2))-Fixed64ToFloat64(fix3), fix1, fix2, fix1.Mul(fix2), fix3))
					}
				}
			}
		}
	}
}
