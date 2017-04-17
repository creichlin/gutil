package gutil

import (
	"log"
)

//// utils for making error handling easier

/*
Shortcut for pattern:
if err != nil {
  log.Fatalf("bla %v, %v", x, err)
}
gutil.FatalIf("bla %v, %v", x, err)

if one of the parameters after the format string is of type error (implicitly not nil)
log.Fatal will be called
*/

func FatalIf(message string, values ...interface{}) {
	for _, value := range values {
		_, isErr := value.(error)
		if isErr {
			log.Fatalf(message, values...)
			return
		}
	}
}

/*
instead of using repeatedly

_, err := callSomethingFoo()
if err != nil {
  // do something
}
err := callSomethingBar()
if err != nil {
  // do something
}
...

errors := gutil.NewErrorCollector()
_, err := callSomethingFoo()
errors.Add(err)
_, err := callSomethingBar()
errors.Add(err)
if err.Has() {
  // do something
}

or it can be a parameter to functions so functions collect the errors or a member of a struct where
code can deposit errors
*/

type ErrorCollector struct {
	errors []error
}

func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{}
}

func (ec *ErrorCollector) Add(err error) bool {
	if err != nil {
		ec.errors = append(ec.errors, err)
		return true
	}
	return false
}

func (ec *ErrorCollector) Has() bool {
	return len(ec.errors) > 0
}

func (ec *ErrorCollector) First() error {
	if ec.Has() {
		return ec.errors[0]
	}
	return nil
}

func (ec *ErrorCollector) Last() error {
	if ec.Has() {
		return ec.errors[len(ec.errors)-1]
	}
	return nil
}

func (ec *ErrorCollector) All() error {
	if ec.Has() {
		return ec
	}
	return nil
}

// ThisOrNil returns the instance if there are errors
// inside, otherwise nil. This is usefull to mimik the
// default error handling where nil means no error
func (ec *ErrorCollector) ThisOrNil() *ErrorCollector {
	if ec.Has() {
		return ec
	}
	return nil
}

func (ec *ErrorCollector) StringList() []string {
	errs := []string{}
	for _, err := range ec.errors {
		errs = append(errs, err.Error())
	}
	return errs
}

// Merge
// creates a new errors collector with the errors from both
// parameters
func Merge(a, b *ErrorCollector) *ErrorCollector {
	m := NewErrorCollector()
	m.errors = append(m.errors, a.errors...)
	m.errors = append(m.errors, b.errors...)
	return m
}

func (ec *ErrorCollector) Error() string {
	msg := ""
	for _, err := range ec.errors {
		msg += err.Error() + "\n"
	}
	if len(msg) == 0 {
		return "Empty error, check for errors on gutil.ErrorCollector using errors.Has() or errors.All() not err != nil"
	}
	return msg[:len(msg)-1]
}
