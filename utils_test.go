package tsterrors

import (
	"errors"
	"testing"
)

func TestProduceValidationCode(t *testing.T) {
	var errInvestmentLost ErrorTag = "investment-lost"

	err1 := New(FunctionName("function_1"))
	err1.Set(errInvestmentLost, errors.New("investment lost"))
	err2 := New(FunctionName("function_2"))
	err2.Pkg(err1)
	err3 := New(FunctionName("function_3"))
	err3.Pkg(err2)

	code, err := ProduceValidationCode(err3)
	if err != nil {
		t.Error(err)
	}

	expectedRes := `var function3 tsterrors.FunctionName = "function_3"
var function2 tsterrors.FunctionName = "function_2"
var function1 tsterrors.FunctionName = "function_1"
var investmentlost tsterrors.ErrorTag = "investment-lost"
if terr.IsRoute(investmentlost, function3, function2, function1) { // todo your validation }`

	if code != expectedRes {
		t.Error(code, expectedRes)
	}
}

func TestProduceValidationCodeFailures(t *testing.T) {
	err1 := New(FunctionName("function_1"))
	_, err := ProduceValidationCode(err1)
	if err == nil {
		t.Error("should fail on empty error")
	}

	_, err = ProduceValidationCode(nil)
	if err == nil {
		t.Error("should fail on nil error")
	}
}
