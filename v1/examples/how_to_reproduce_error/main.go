package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rm4n0s/tsterrors/v1"
)

func main() {
	var errInvestmentLost tsterrors.ErrorTag = "InvestmentLost"

	// microservice 2 does a HTTP request to microservice 1
	// microservice 1 runs function_4() but it fails at function_1()
	// The error propagates from function_1() to function_4()
	err1 := tsterrors.New(tsterrors.FunctionName("function_1"))
	err1.Set(errInvestmentLost, errors.New("investment lost"))
	err2 := tsterrors.New(tsterrors.FunctionName("function_2"))
	err2.Pkg(err1)
	err3 := tsterrors.New(tsterrors.FunctionName("function_3"))
	err3.Pkg(err2)
	err4 := tsterrors.New(tsterrors.FunctionName("function_4"))
	err4.Pkg(err3)

	// microservice 1 transforms the error from function_4() to json
	out, _ := json.Marshal(err4)

	// microservice 1 sends the json to microservice 2 as a response
	// microservice 2 transforms the json to an error
	in := &tsterrors.Error{}
	json.Unmarshal(out, &in)

	// microservice 2 propagates the error to function_7()
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

	// the sysadmin receives the email and copy pastes the json to QA engineer to reproduce the error in his/her test enviroment

	// the engineer puts the json in Error object
	in = &tsterrors.Error{}
	json.Unmarshal(out, &in)

	// so he/she can use the ProduceValidationCode utility to create the validation code for his/her test code
	code, _ := tsterrors.ProduceValidationCode(in)
	fmt.Println("Code: ", code)

	// the engineer copies the code into a test and tries to reproduce the same error with different inputs
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
