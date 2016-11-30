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

func (ec *ErrorCollector) Add(err error) {
	if err != nil {
		ec.errors = append(ec.errors, err)
	}
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

func (ec *ErrorCollector) Error() string {
	msg := ""
	for _, err := range ec.errors {
		msg += err.Error() + "\n"
	}
	return msg[:len(msg)-1]
}
