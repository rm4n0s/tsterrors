package main

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/rm4n0s/tsterrors/v1"
)

var errBankAccountEmpty tsterrors.ErrorTag = "AccountIsEmpty"
var errInvestmentLost tsterrors.ErrorTag = "InvestmentLost"

var fNf1 tsterrors.FunctionName = "f1"

func f1() error {
	terr := tsterrors.New(fNf1)
	n := rand.IntN(9) + 1
	if n%2 == 0 {
		return terr.Set(errBankAccountEmpty, errors.New("account is empty"))
	}
	return terr.Set(errInvestmentLost, errors.New("investment is lost"))
}

var fNf2 tsterrors.FunctionName = "f2"

func f2() error {
	terr := tsterrors.New(fNf2)
	err := f1()
	return terr.Pkg(err)

}

var fNf3 tsterrors.FunctionName = "f3"

func f3() error {
	terr := tsterrors.New(fNf3)
	err := f1()
	return terr.Pkg(err)
}

var fNf4 tsterrors.FunctionName = "f4"

func f4() error {
	terr := tsterrors.New(fNf4)
	n := rand.IntN(9) + 1
	if n%2 == 0 {
		return terr.Pkg(f2())
	}
	return terr.Pkg(f3())
}

func main() {
	err := f4()
	fmt.Println(err)

	terr := err.(*tsterrors.Error)
	fmt.Printf("%#v \n", terr)
	fmt.Println(terr.StackTrace())

	if terr.IsRoute(errInvestmentLost, fNf4, fNf3, fNf1) {
		fmt.Println("The money in your account didn't do well")
	} else if terr.IsRoute(errBankAccountEmpty, fNf4, fNf2, fNf1) {
		fmt.Println("Aand it's gone")
	} else {
		fmt.Println("This line is for bank members only")
	}
}
