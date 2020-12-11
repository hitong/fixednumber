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

func TestConvAllBranch(t *testing.T){
	defer func() {
		if err := recover(); err != nil{
			t.Fail()
		}
		onceSet = &sync.Once{}
	}()
	SetPrecisionOnce(20)
	MaxFixed64.ToBase10s(18)
	PrecisionNumber.Add(MaxFixed64).ToBase10s(18)
	MaxFixed64.Mul(MaxFixed64).Float64()
	MaxFixed64.ToBase10N(0)
	insertToFloatSliceBase10(0,0)
	Fixed64Zero.Div(MaxFixed64)
	Float64ToFixed64(0.9).ToBase10s(1)
	Float64ToFixed64(1 << 40 + 0.15464)
	onceSet = &sync.Once{}
	SetPrecisionOnce(60)
	fmt.Println(Float64ToFixed64(0.05))
	onceSet = &sync.Once{}
	SetPrecisionOnce(5)
	Float64ToFixed64(2 << 53)
}

func TestFloat64ToFixed64Overflow(t *testing.T){
	defer func() {
		if err := recover();err != nil{
			if err == "Fixed number: part digital overflow"{
				return
			}
		}

		t.Fail()
	}()
	SetPrecisionOnce(20)
	Float64ToFixed64(math.Pow10(100))
}

func TestPanicFixedConvBase10(t *testing.T){
	defer func() {
		if err := recover(); err != nil{
			if err == "Fixed64.ToBase10N: Not Support n > 19"{
				return
			}
		}

		t.Fail()
	}()

	SetPrecisionOnce(20)
	defer func() {onceSet = &sync.Once{}}()
	Float64ToFixed64(1).ToBase10N(88)
}

func TestPanicPrecisionSet(t *testing.T){
	defer func() {
		err := recover()
		if err != "Precision overflow"{
			t.Fail()
		}
	}()
	SetPrecisionOnce(63)
	defer func() {onceSet = &sync.Once{}}()
}


func TestBasic(t *testing.T) {
	SetPrecisionOnce(20)
	defer func() {onceSet = &sync.Once{}}()
	f0 := Float64ToFixed64(123.456)
	f1 := Float64ToFixed64(123.456)

	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}

	if f0.Int64() != 123 {
		t.Error("should be equal", f0.Int64(), 123)
	}

	if f0.ToBase10s(3) != "123.456" {
		t.Error("should be equal", f0.String(), "123.456")
	}

	f0 = Float64ToFixed64(0.499)
	f1 = Float64ToFixed64(-0.5011)
	if f1.Int64() != 0{
		t.Error("should be round to equal ", 0 , f1.Int64())
	}

	if f0.Round() != 0{
		t.Error("should be round to equal ", 0 , string(f0.ToBase10()), f0.Round())
	}

	f0 = Float64ToFixed64(999.999)
	if f0.ToBase10s(2) != "1000.00"{
		t.Error("should be round to equal ", "1000.00" , f0.ToBase10s(2))
	}

	if f1.Round() != -1{
		t.Error("should be round to equal ",-1,f1.ToBase10s(3), f1.Round())
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
	if f0.ToBase10s(3) != "0.999" {
		t.Error("should be equal", f0, "0.999")
	}

	f0 = Float64ToFixed64(.331)
	f1 = Float64ToFixed64(.332)
	if !f0.Less(f1){
		t.Error("should be less:",f0,f1)
	}

	if !f1.Great(f0){
		t.Error("should be great:",f1,f0)
	}

	if !f0.Equal(f0.Mul(Float64ToFixed64(-1)).Abs()){
		t.Error("should be equal ",f1,f0.Mul(Float64ToFixed64(-1)).Abs())
	}

	f0 = Float64ToFixed64(1)
	f1 = Float64ToFixed64(0.1)
	f2 = Float64ToFixed64(10)
	if f0.Div(f1).ToBase10s(3) != f2.ToBase10s(3) {
		t.Error("should be equal ",f0.Div(f1),f2)
	}

	if f0.Div(f2).ToBase10s(3) != f1.ToBase10s(3){
		t.Error("should be equal ",f0.Div(f2),f1)
	}

	if !Fixed64Zero.Mul(MaxFixed64).Equal(MaxFixed64.Mul(Fixed64Zero)){
		t.Error("should be equal  ",Fixed64Zero.Mul(MaxFixed64),MaxFixed64.Mul(Fixed64Zero))
	}

	if !MaxFixed64.Sub(MaxFixed64).Equal(Fixed64Zero){
		t.Error("should be equal  ",MaxFixed64.Sub(MaxFixed64),Fixed64Zero)
	}
}

func TestUint64Bits(t *testing.T) {
	SetPrecisionOnce(20)
	defer func() {onceSet = &sync.Once{}}()
	if string(Uint64Bits(uint64(0b11000111001010101010))) != "0000000000000000000000000000000000000000000011000111001010101010"{
		t.Error("should be equal ",string(Uint64Bits(uint64(0b11000111001010101010))), "0000000000000000000000000000000000000000000011000111001010101010")
	}
}

func TestSafeFloat64ToFixed64(t *testing.T) {
	SetPrecisionOnce(20)
	defer func() {onceSet = &sync.Once{}}()
	if _, err := SafeFloat64ToFixed64(math.Float64frombits(^uint64(0))); err != nil {
		if err.Error() != "Float64 value is NaN "{
			t.Error("should be NaN:" ,^uint64(0))
		}
	}

	if _, err := SafeFloat64ToFixed64(math.Float64frombits((1 << 11 - 1) << 52)); err != nil {
		if err.Error() != "Float64 value is Inf "{
			t.Error("should be NaN:" ,^uint64(0))
		}
	}

	if v, _ := SafeFloat64ToFixed64(math.Float64frombits(0b0011110101010101001)); v !=0 {
		t.Error("The non-normalized number should be 0")
	}
}

func TestVerifyFixedCompute(t *testing.T){
	//save to file
	buildData("AddData.txt", "Add", 1000,3)
	buildData("SubData.txt", "Sub", 1000,3)
	buildData("MulData.txt", "Mul", 1000,3)
	buildData("DivData.txt", "Div", 1000,3)

	paths := []string{"AddData.txt", "SubData.txt", "DivData.txt", "MulData.txt"}
	SetPrecisionOnce(20)

	for _, path := range paths {
		if op, v1, v2, v3, err := readData(path); err != nil {
			panic(err)
		} else {
			for i := 0; i < len(op); i++ {
				fix1 := Float64ToFixed64(v1[i])
				fix2 := Float64ToFixed64(v2[i])
				switch op[i] {
				case "Add":
					if strings.Compare(string(fix1.Add(fix2).ToBase10N(3)),fmt.Sprintf("%.3f",v3[i]))!= 0{
						println("Add:" + fmt.Sprintf("%.f %s %s %s %.3f", fix1.Add(fix2).Float64()-v3[i], fix1, fix2, fix1.Add(fix2).ToBase10N(3), v3[i]))
					}
				case "Sub":
					if strings.Compare(string(fix1.Sub(fix2).ToBase10N(3)),fmt.Sprintf("%.3f",v3[i]))!= 0{
						println("Sub:" + fmt.Sprintf("%f %s %s %s %.3f", fix1.Sub(fix2).Float64()-v3[i], fix1, fix2, fix1.Sub(fix2).ToBase10N(3), v3[i]))
					}
				case "Div":
					if strings.Compare(string(fix1.Div(fix2).ToBase10N(3)),fmt.Sprintf("%.3f",v3[i]))!= 0{
						println("Div:" + fmt.Sprintf("%f %s %s %s %.3f", fix1.Div(fix2).Float64()-v3[i], fix1, fix2, fix1.Div(fix2).ToBase10N(3), v3[i]))
					}
				case "Mul":
					if strings.Compare(string(fix1.Mul(fix2).ToBase10N(3)),fmt.Sprintf("%.3f",v3[i]))!= 0 {
						println("Mul:" + fmt.Sprintf("%f %s %s %s %.3f", fix1.Mul(fix2).Float64()-v3[i], fix1, fix2, fix1.Mul(fix2).ToBase10N(3), v3[i]))
					}
				}
			}
		}
	}
}

func readData(fileName string) (op []string, v1 []float64, v2 []float64, v3 []float64, err error) {
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
		num1, err1 := strconv.ParseFloat(nums[0], 64)
		num2, err2 := strconv.ParseFloat(nums[1], 64)
		num3, err3 := strconv.ParseFloat(nums[2], 64)
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

func buildData(fileName, op string, num int,decimalSize int) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	var fmtFloat = "%." + fmt.Sprintf("%d",decimalSize) + "f "
	for i := 0; i < num; i++ {
		v1 := math.Floor((rand.Float64()+float64(rand.Int()%1024))*1000) / 1000
		v2 := math.Floor((rand.Float64()+float64(rand.Int()%1024))*1000) / 1000
		var v3 float64
		if op == "Add" {
			v3 = v1 + v2
		}
		if op == "Sub" {
			v3 = v1 - v2
		}
		if op == "Mul" {
			v3 = v1 * v2
		}
		if op == "Div" {
			v3 = v1 / v2
		}

		file.WriteString(fmt.Sprintf("%s:" + fmtFloat + fmtFloat + fmtFloat + "\n", op, v1, v2, v3))
	}
	file.Close()
}
