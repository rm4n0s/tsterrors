package tsterrors

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

// ProduceValidationCode will produce a code, based on Error's object, that will validate the stack trace of that specific error.
func ProduceValidationCode(terr *Error) (string, error) {
	if terr == nil {
		return "", errors.New("error is nil")
	}

	if terr.State == StateEmpty {
		return "", errors.New("can't produce code in a state of empty error")
	}

	fns := []string{string(terr.CurrentFunction)}
	et := ""
	prvErr := terr.PrvErr
	for prvErr != nil {
		curr := prvErr
		fns = append(fns, string(curr.CurrentFunction))
		prvErr = curr.PrvErr
		if prvErr == nil {
			et = string(curr.ErrorTag)
		}
	}

	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	str := ""
	vrbfns := []string{}
	for _, fn := range fns {
		vrbFn := nonAlphanumericRegex.ReplaceAllString(fn, "")
		str += fmt.Sprintf("var %s tsterrors.FunctionName = \"%s\"\n", vrbFn, fn)
		vrbfns = append(vrbfns, vrbFn)
	}

	vrbEt := nonAlphanumericRegex.ReplaceAllString(et, "")
	str += fmt.Sprintf("var %s tsterrors.ErrorTag = \"%s\"\n", vrbEt, et)

	arg1 := vrbEt
	arg2 := strings.Join(vrbfns, ", ")
	str += fmt.Sprintf("if terr.IsRoute(%s, %s) { // todo your validation }", arg1, arg2)
	return str, nil
}

func callerName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return ""
	}
	return f.Name()
}
