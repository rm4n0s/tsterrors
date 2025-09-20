package tsterrors

// FunctionName is the type for the name of the function where the Error object initialized.
type FunctionName string

// ErrorTag is the type for the name of the error, preferably in camel case
type ErrorTag string

// State is the type of the Error's state
type State int

const (
	// StateEmpty is the state where the Error object doesn't contain any error
	StateEmpty State = 0
	// StateFirst is the state where the Error object contains the original error
	StateFirst State = 1
	// StatePackage is the state where the Error object contains other Error object
	StatePackage State = 2
)

func (e ErrorTag) Error() string {
	return string(e)
}

// Error is an object to carry the original error or become a linked list of errors to trace the call path of the original error
type Error struct {
	ErrorString     string       // string from the original error
	CurrentFunction FunctionName // name of the function where this error initialized
	State           State        // the state of the error
	FirstErr        error        // the first error
	ErrorTag        ErrorTag     // the tag of the error
	PrvErr          *Error       // packaged Error
}

// New initializes a Error object with the name of the function.
// The function panics when the name is empty.
func New(fn FunctionName) *Error {
	if len(fn) == 0 {
		panic("function's name can't be empty")
	}
	return &Error{
		CurrentFunction: fn,
		State:           StateEmpty,
	}
}

// NewAuto initializes a Error object with the name of the function automatically given on runtime.
func NewAuto() *Error {
	return &Error{
		CurrentFunction: FunctionName(callerName(1)),
	}
}

// Route will return the Error object from a linked list of Errors based on the order of function names as arguments.
// The order of function name should be the same of function calls in a stack trace.
func (e *Error) Route(fns ...FunctionName) *Error {
	if len(fns) == 0 {
		return nil
	}

	curr := e
	for i, fn := range fns {
		if curr != nil {
			if curr.CurrentFunction == fn {
				if curr.PrvErr != nil {
					if i == len(fns)-1 {
						return curr
					}
					curr = curr.PrvErr
				}
			} else {
				return nil
			}
		}
	}

	return curr
}

// IsRoute will confirm if the order of function names and error's tag exist in the Error's linked list.
// The order of function name should be the same of function calls in a stack trace.
func (e *Error) IsRoute(et ErrorTag, fns ...FunctionName) bool {
	if len(fns) == 0 {
		return e.ErrorTag == et && e.State == StateFirst
	}
	r := e.Route(fns...)
	if r == nil {
		return false
	}
	return r.ErrorTag == et
}

func (e *Error) Error() string {
	return e.ErrorString
}

// Set is the method that is needed to insert the original error inside the Error's object.
// The method will panic if the tsterrors.Error object is not empty.
func (e *Error) Set(et ErrorTag, err error) error {
	if err == nil {
		return nil
	}
	if e.State != StateEmpty {
		panic("don't reset errors in Error object")
	}

	e.ErrorTag = et
	e.State = StateFirst
	e.FirstErr = err
	e.ErrorString = err.Error()
	return e
}

// Pkg is the method that will insert another Error's object.
// The method will panic if the error is not a tsterrors.Error or it is not empty.
func (e *Error) Pkg(err error) error {
	if err == nil {
		return nil
	}
	terr, ok := err.(*Error)
	if !ok {
		panic("error is not created by 'tsterrors'")
	}
	if e.State != StateEmpty {
		panic("don't reset errors in Error object")
	}
	e.PrvErr = terr
	e.State = StatePackage
	e.ErrorString = terr.ErrorString
	e.ErrorTag = terr.ErrorTag
	e.FirstErr = terr.FirstErr
	return e
}

// StackTrace will return a list of functions, in the order they called and seperated by "->", and at the end will be the tag of the original error.
func (e *Error) StackTrace() string {
	str := string(e.CurrentFunction)
	prvErr := e.PrvErr
	for prvErr != nil {
		curr := prvErr
		str += " -> " + string(curr.CurrentFunction)
		prvErr = curr.PrvErr
		if prvErr == nil {
			str += "." + string(curr.ErrorTag)
		}
	}
	return str
}
