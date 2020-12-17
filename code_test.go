package fixed

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestConvAllBranch(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
		onceSet = &sync.Once{}
	}()
	SetPrecisionOnce(20)
	MaxFixed64.ToBase10N(18)
	MaxFixed64.ToBase10N(0)
	insertToFloatSliceBase10(0, 0)
	Fixed64Zero.Div(MaxFixed64)
	Float64ToFixed64(0.9).ToBase10N(1)
	Float64ToFixed64(1<<40 + 0.15464).Float64()
	Fixed64(1).Float64()
	onceSet = &sync.Once{}
	SetPrecisionOnce(60)
	Float64ToFixed64(0.05).String()
	onceSet = &sync.Once{}
	SetPrecisionOnce(5)
	Float64ToFixed64(2 << 53)
}

func TestFloat64AddOverFlow(t *testing.T){
	defer func() {
		if err := recover();err == nil || !strings.Contains(err.(string), "Fixed64: Add Overflow"){
			t.Fail()
		}

		onceSet = &sync.Once{}
	}()
	SetPrecisionOnce(20)
	f0 := Fixed64((1 << 63)- 1)
	f1 := Fixed64(1)
	f0.Add(f1)
}

func TestFloat64MulOverFlow(t *testing.T){
	defer func() {
		if err := recover();err == nil || err != "Fixed64: Number OverFlow"{
			t.Fail()
		}

		onceSet = &sync.Once{}
	}()
	SetPrecisionOnce(20)
	f0 := Float64ToFixed64(5744874.6666664)
	f1 := Float64ToFixed64(12124678.6666664)
	f0.Mul(f1)
}

func TestFloat64ToFixed64Overflow(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			if err == "Fixed number: part digital overflow" {
				return
			}
		}
		t.Fail()
	}()
	SetPrecisionOnce(20)
	Float64ToFixed64(math.Pow10(100))
}

func TestPanicFixedConvBase10(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			if err == "Fixed64.ToBase10N: Not Support n > 19" {
				return
			}
		}

		t.Fail()
	}()

	SetPrecisionOnce(20)
	defer func() { onceSet = &sync.Once{} }()
	Float64ToFixed64(1).ToBase10N(88)
}

func TestPanicPrecisionSet(t *testing.T) {
	defer func() {
		err := recover()
		if err != "Precision overflow" {
			t.Fail()
		}
	}()
	SetPrecisionOnce(63)
	defer func() { onceSet = &sync.Once{} }()
}

func TestStr2Fixed64(t *testing.T) {
	SetPrecisionOnce(20)
	defer func() { onceSet = &sync.Once{} }()
	var values = []string{"1248.154","0.111","3.133","999.489","48487",".144848","1458.","-1564.154","-0.1"}
	for _,value := range values{
		if fix,err := Str2Fixed64(value); err != nil{
			t.Error(err)
		} else {
			fmt.Println(value,fix)
		}
	}
	values = []string{"12.48.154",".",".00.1215"}
	for _,value := range values{
		if _,err := Str2Fixed64(value); err == nil{
			t.Error("should has err ", value)
		}
	}
}

func TestBasic(t *testing.T) {
	SetPrecisionOnce(20)
	var v = Fixed64((1 << 63) - 1).Div(Fixed64(4561356462))
	fmt.Println(v)
	defer func() { onceSet = &sync.Once{} }()
	f0 := Float64ToFixed64(123.456)
	f1 := Float64ToFixed64(123.456)

	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}

	if f0.Int64() != 123 {
		t.Error("should be equal", f0.Int64(), 123)
	}

	if f0.ToBase10N(3) != "123.456" {
		t.Error("should be equal", f0.String(), "123.456")
	}

	f0 = Float64ToFixed64(0.499)
	f1 = Float64ToFixed64(-0.5011)
	if f1.Int64() != 0 {
		t.Error("should be round to equal ", 0, f1.Int64())
	}

	if f0.Round() != 0 {
		t.Error("should be round to equal ", 0, string(f0.ToBase10()), f0.Round())
	}

	f0 = Float64ToFixed64(999.999)
	if f0.ToBase10N(2) != "1000.00" {
		t.Error("should be round to equal ", "1000.00", f0.ToBase10N(2))
	}

	if f1.Round() != -1 {
		t.Error("should be round to equal ", -1, f1.ToBase10N(3), f1.Round())
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
	if f0.ToBase10N(3) != "0.999" {
		t.Error("should be equal", f0, "0.999")
	}

	f0 = Float64ToFixed64(.331)
	f1 = Float64ToFixed64(.332)
	if !f0.Less(f1) {
		t.Error("should be less:", f0, f1)
	}

	if !f1.Great(f0) {
		t.Error("should be great:", f1, f0)
	}

	if !f0.Equal(f0.Mul(Float64ToFixed64(-1)).Abs()) {
		t.Error("should be equal ", f1, f0.Mul(Float64ToFixed64(-1)).Abs())
	}

	f0 = Float64ToFixed64(1)
	f1 = Float64ToFixed64(0.1)
	f2 = Float64ToFixed64(10)
	if f0.Div(f1).ToBase10N(3) != f2.ToBase10N(3) {
		t.Error("should be equal ", f0.Div(f1), f2)
	}

	if f0.Div(f2).ToBase10N(3) != f1.ToBase10N(3) {
		t.Error("should be equal ", f0.Div(f2), f1)
	}

	if !Fixed64Zero.Mul(MaxFixed64).Equal(MaxFixed64.Mul(Fixed64Zero)) {
		t.Error("should be equal  ", Fixed64Zero.Mul(MaxFixed64), MaxFixed64.Mul(Fixed64Zero))
	}

	if !MaxFixed64.Sub(MaxFixed64).Equal(Fixed64Zero) {
		t.Error("should be equal  ", MaxFixed64.Sub(MaxFixed64), Fixed64Zero)
	}

	f0 = Fixed64((1 << 63) - 1)
	f1 = Fixed64(1 << 20)
	if !f0.Mul(f1).Equal(f0){
		t.Error("should be equal ",f0.Mul(f1), f0)
	}
}

func TestUint64Bits(t *testing.T) {
	SetPrecisionOnce(20)
	defer func() { onceSet = &sync.Once{} }()
	if string(Uint64Bits(uint64(0b11000111001010101010))) != "0000000000000000000000000000000000000000000011000111001010101010" {
		t.Error("should be equal ", string(Uint64Bits(uint64(0b11000111001010101010))), "0000000000000000000000000000000000000000000011000111001010101010")
	}
}

func TestSafeFloat64ToFixed64(t *testing.T) {
	SetPrecisionOnce(20)
	var a =  49877.641
	var b = 4588.648
	fmt.Println(Float64ToFixed64(a).Add(Float64ToFixed64(b)))
	fmt.Println(Float64ToFixed64(a).Sub(Float64ToFixed64(b)))
	fmt.Println(Float64ToFixed64(a).Mul(Float64ToFixed64(b)))
	fmt.Println(Float64ToFixed64(a).Div(Float64ToFixed64(b)))
	defer func() { onceSet = &sync.Once{} }()
	if _, err := SafeFloat64ToFixed64(math.Float64frombits(^uint64(0))); err != nil {
		if err.Error() != "Float64 value is NaN " {
			t.Error("should be NaN:", ^uint64(0))
		}
	}

	if _, err := SafeFloat64ToFixed64(math.Float64frombits((1<<11 - 1) << 52)); err != nil {
		if err.Error() != "Float64 value is Inf " {
			t.Error("should be NaN:", ^uint64(0))
		}
	}

	if v, _ := SafeFloat64ToFixed64(math.Float64frombits(0b0011110101010101001)); v != 0 {
		t.Error("The non-normalized number should be 0")
	}
}

func TestVerifyFixedCompute(t *testing.T) {
	//save to file
	SetPrecisionOnce(20)
	buildData("AddData.txt", "Add", 1000, 3)
	buildData("SubData.txt", "Sub", 1000, 3)
	buildData("MulData.txt", "Mul", 1000, 3)
	buildData("DivData.txt", "Div", 1000, 3)

	paths := []string{"AddData.txt", "SubData.txt", "DivData.txt", "MulData.txt"}


	for _, path := range paths {
		if op, v1, v2, v3, err := readData(path); err != nil {
			panic(err)
		} else {
			for i := 0; i < len(op); i++ {
				fix1 := Fixed64(v1[i])
				fix2 := Fixed64(v2[i])
				switch op[i] {
				case "Add":
					if !fix1.Add(fix2).Equal(Fixed64(v3[i])) {
						t.Error("Add:" + fmt.Sprintf("%f != %f", fix1.Add(fix2).Float64(), Fixed64(v3[i]).Float64()))
					}
				case "Sub":
					if !fix1.Sub(fix2).Equal(Fixed64(v3[i])) {
						t.Error("Sub:" + fmt.Sprintf("%f != %f", fix1.Add(fix2).Float64(), Fixed64(v3[i]).Float64()))
					}
				case "Div":
					if !fix1.Div(fix2).Equal(Fixed64(v3[i])) {
						t.Error("Div:" + fmt.Sprintf("%f != %f", fix1.Add(fix2).Float64(), Fixed64(v3[i]).Float64()))
					}
				case "Mul":
					if !fix1.Mul(fix2).Equal(Fixed64(v3[i])) {
						t.Error("Mul:" + fmt.Sprintf("%f != %f", fix1.Add(fix2).Float64(), Fixed64(v3[i]).Float64()))
					}
				}
			}
		}
	}
}

func readData(fileName string) (op []string, v1 []uint64, v2 []uint64, v3 []uint64, err error) {
	var file *os.File
	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0666)
	defer file.Close()
	reader := bufio.NewReader(file)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	for {
		var data string
		data, err = reader.ReadString('\n')
		if err != nil {
			err = nil
			break
		}
		opAndData := strings.Split(data, ":")
		if len(opAndData) != 2 {
			return nil, nil, nil, nil, errors.New("data in wrong format:" + data)
		}
		nums := strings.Split(opAndData[1], " ")
		if len(nums) != 4 {
			return nil, nil, nil, nil, errors.New("data in wrong format")
		}
		num1, err1 := strconv.ParseUint(nums[0],10, 64)
		num2, err2 :=  strconv.ParseUint(nums[1],10, 64)
		num3, err3 :=  strconv.ParseUint(nums[2],10, 64)
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, nil, nil, nil, errors.New("parse number error")
		}
		op = append(op, opAndData[0])
		v1 = append(v1, num1)
		v2 = append(v2, num2)
		v3 = append(v3, num3)
	}
	return
}

func buildData(fileName, op string, num int, decimalSize int) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	//var fmtFloat = "%." + fmt.Sprintf("%d", decimalSize) + "f "
	var fmtFloat = "%d "
	for i := 0; i < num; i++ {
		f1 := Float64ToFixed64(math.Floor((rand.Float64()+float64(rand.Int()%1024))*1000) / 1000.0)
		f2 := Float64ToFixed64(math.Floor((rand.Float64()+float64(rand.Int()%1024))*1000) / 1000.0)
		var f3 Fixed64
		if op == "Add" {
			f3 =f1.Add(f2)
		}
		if op == "Sub" {
			f3 =f1.Sub(f2)
		}
		if op == "Mul" {
			f3 =f1.Mul(f2)
		}
		if op == "Div" {
			f3 =f1.Div(f2)
		}
		file.WriteString(fmt.Sprintf("%s:"+fmtFloat+fmtFloat+fmtFloat+"\n", op, uint64(f1), uint64(f2), uint64(f3)))
	}
	file.Close()
}
