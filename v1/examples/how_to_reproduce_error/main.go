package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rm4n0s/tsterrors/v1"
)

func main() {
	var errInvestmentLost tsterrors.ErrorTag = "InvestmentLost"

	// from microservice 1
	err1 := tsterrors.New(tsterrors.FunctionName("function_1"))
	err1.Set(errInvestmentLost, errors.New("investment lost"))
	err2 := tsterrors.New(tsterrors.FunctionName("function_2"))
	err2.Pkg(err1)
	err3 := tsterrors.New(tsterrors.FunctionName("function_3"))
	err3.Pkg(err2)
	err4 := tsterrors.New(tsterrors.FunctionName("function_4"))
	err4.Pkg(err3)

	out, _ := json.Marshal(err4)

	// send error to microservice 2 as a response to HTTP Request
	in := &tsterrors.Error{}
	json.Unmarshal(out, &in)

	err5 := tsterrors.New(tsterrors.FunctionName("function_5"))
	err5.Pkg(in)
	err6 := tsterrors.New(tsterrors.FunctionName("function_6"))
	err6.Pkg(err5)
	err7 := tsterrors.New(tsterrors.FunctionName("function_7"))
	err7.Pkg(err6)

	// microservice 2 sends the error as json in an email to sysadmin, with its stacktrace
	fmt.Println("Stacktrace: ", err7.StackTrace())
	out, _ = json.Marshal(err7)
	fmt.Println("Json", string(out))

	// then the sysadmin gives the json to QA engineer to reproduce the bug
	// the engineer puts the json in Error object

	in = &tsterrors.Error{}
	json.Unmarshal(out, &in)

	// then the engineer produces, from Go's REPL, a code to validate similar errors
	code, _ := tsterrors.ProduceValidationCode(in)
	fmt.Println("Code: ", code)

	// the engineer copies the code into a test and tries to reproduce with different inputs
	TestWeirdBug := func() {
		reproducableError := err7
		var function7 tsterrors.FunctionName = "function_7"
		var function6 tsterrors.FunctionName = "function_6"
		var function5 tsterrors.FunctionName = "function_5"
		var function4 tsterrors.FunctionName = "function_4"
		var function3 tsterrors.FunctionName = "function_3"
		var function2 tsterrors.FunctionName = "function_2"
		var function1 tsterrors.FunctionName = "function_1"
		var investmentlost tsterrors.ErrorTag = "InvestmentLost"
		if reproducableError.IsRoute(investmentlost, function7, function6, function5, function4, function3, function2, function1) {
			fmt.Println("Yes! the engineer reproduced it! ")
		}
	}
	TestWeirdBug()
}
