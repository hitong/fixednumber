# FixedNumber

[![codecov](https://codecov.io/gh/hitong/FixedNumber/branch/master/graph/badge.svg)](https://codecov.io/gh/hitong/FixedNumber/branch/master/graphs)
[![Go Report Card](https://goreportcard.com/badge/github.com/hitong/FixedNumber)](https://goreportcard.com/report/github.com/hitong/FixedNumber)
[![Build Status](https://travis-ci.org/hitong/FixedNumber.svg?branch=master)](https://travis-ci.org/hitong/FixedNumber)


## Usage
https://play.golang.org/p/iW1rFy0_Y_j
<pre><code>
package main

import (
	"fmt"

	fixed "github.com/hitong/fixednumber"
)

func main() {
	fixed.SetPrecisionOnce(20)              //Initialization of fixed point precision
	strFix, _ := fixed.Str2Fixed64("0.1")   // from string
	floatFix := fixed.Float64ToFixed64(0.1) // from float64
	fmt.Println("strFix", strFix)
	fmt.Println("floatFix", floatFix)

	//Basic calculation
	resFix := strFix.Add(floatFix)
	fmt.Println("add res", resFix)

	//forcibly convert data into uint64 for data compression and transmission
	convRes := uint64(resFix)
	convFix := fixed.Fixed64(convRes)
	fmt.Println("src data", resFix)
	fmt.Println("conv data", convFix)

	var n uint = 10
	fmt.Println(fmt.Sprintf("keep %d as a decimal %s", n, resFix.ToBase10N(n)))
}

</code></pre>
