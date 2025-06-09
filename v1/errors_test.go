package tsterrors

import (
	"errors"
	"testing"
)

func TestNewFailures(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the New() should have paniced")
		}
	}()

	New("")
}

func TestNewSuccess(t *testing.T) {
	terr := New("TestNewSuccess")
	if terr.CurrentFunction != "TestNewSuccess" {
		t.Error("failed to save function name")
	}
}

func TestNewAuto(t *testing.T) {
	terr := NewAuto()

	if "github.com/rm4n0s/tsterrors/v1.TestNewAuto" != terr.CurrentFunction {
		t.Error("failed to save function name")
	}
}

func TestSetFailures(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the Set() should have paniced")
		}
	}()
	err1 := New(FunctionName("function_1"))
	err1.State = StatePackage
	err1.Set("InvestmentLost", errors.New("investment lost"))
}

func TestSetSuccess(t *testing.T) {
	err1 := New(FunctionName("function_1"))
	err := err1.Set("InvestmentLost", errors.New("investment lost"))
	terr := err.(*Error)

	if terr.ErrorTag != "InvestmentLost" || terr.State != StateFirst {
		t.Error("didn't save correctly errortag and state")
	}
}

func TestPkgFailOnNonTstError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the Pkg() should have paniced")
		}
	}()
	err1 := New(FunctionName("function_1"))
	err1.Pkg(errors.New("investment lost"))
}

func TestPkgFailOnNonIsEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the Pkg() should have paniced")
		}
	}()
	err1 := New(FunctionName("function_1"))
	err1.Set("InvestmentLost", errors.New("investment lost"))
	err2 := New(FunctionName("function_2"))
	err2.Pkg(err1)
	// this will panic
	err2.Pkg(err1)
}

func TestRoute(t *testing.T) {
	var errInvestmentLost ErrorTag = "investment-lost"

	err1 := New(FunctionName("function_1"))
	err1.Set(errInvestmentLost, errors.New("investment lost"))
	err2 := New(FunctionName("function_2"))
	err2.Pkg(err1)
	err3 := New(FunctionName("function_3"))
	err3.Pkg(err2)

	err := err3.Route("function_3", "function_2")
	if err.CurrentFunction != "function_2" || err.PrvErr.CurrentFunction != "function_1" {
		t.Error("didn't route correctly ")
	}

	err = err3.Route("function_3", "function_2", "function_1")
	if err.CurrentFunction != "function_1" {
		t.Error("didn't route correctly ")
	}

	err = err3.Route()
	if err != nil {
		t.Error("should be nil")
	}
}

func TestIsRoute(t *testing.T) {
	var errInvestmentLost ErrorTag = "investment-lost"

	err1 := New(FunctionName("function_1"))
	err1.Set(errInvestmentLost, errors.New("investment lost"))
	err2 := New(FunctionName("function_2"))
	err2.Pkg(err1)
	err3 := New(FunctionName("function_3"))
	err3.Pkg(err2)

	if !err3.IsRoute(errInvestmentLost, "function_3", "function_2", "function_1") {
		t.Error("didn't check route correctly ")
	}

	if err3.IsRoute("", "function_3", "function_2", "function_1") {
		t.Error("didn't check route correctly ")
	}

	if err3.IsRoute(errInvestmentLost, "function_3", "function_1") {
		t.Error("didn't check route correctly ")
	}

	if err3.IsRoute(errInvestmentLost) {
		t.Error("didn't check route correctly ")
	}

	if !err1.IsRoute(errInvestmentLost) {
		t.Error("didn't check route correctly ")
	}
}

func TestError(t *testing.T) {
	var errInvestmentLost ErrorTag = "investment-lost"
	err := errors.New("investment lost")
	err1 := New(FunctionName("function_1"))
	err1.Set(errInvestmentLost, err)
	err2 := New(FunctionName("function_2"))
	err2.Pkg(err1)
	err3 := New(FunctionName("function_3"))
	err3.Pkg(err2)

	if err3.Error() != err.Error() {
		t.Error("the package didn't inherit the first error's string")
	}
}

func TestStacktrace(t *testing.T) {
	var errInvestmentLost ErrorTag = "investment-lost"
	err := errors.New("investment lost")
	err1 := New(FunctionName("function_1"))
	err1.Set(errInvestmentLost, err)
	err2 := New(FunctionName("function_2"))
	err2.Pkg(err1)
	err3 := New(FunctionName("function_3"))
	err3.Pkg(err2)

	if err3.StackTrace() != "function_3 -> function_2 -> function_1.investment-lost" {
		t.Error("the stacktrace is not produced correctly")
	}
}
