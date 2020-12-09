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
)

func ReadData(fileName string) (op []string, v1 []float64, v2 []float64, v3 []float64, err error) {
	var file *os.File
	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
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

func BuildData(fileName, op string, num int) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
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

		file.WriteString(fmt.Sprintf("%s:%f %f %f \n", op, v1, v2, v3))
	}
	file.Close()
}
