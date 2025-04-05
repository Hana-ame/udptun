package main

import "fmt"

type doubleError struct {
	tags   [2]string
	errors [2]error
}

func (e *doubleError) Error() string {
	return e.tags[0] + ": " + e.errors[0].Error() + ", " +
		e.tags[1] + ": " + e.errors[1].Error()
}

func DoubleError(error1, error2 error, tag1, tag2 string) error {
	if error1 == nil {
		return error2
	}
	if error2 == nil {
		return error1
	}
	return &doubleError{
		tags:   [2]string{tag1, tag2},
		errors: [2]error{error1, error2},
	}
}

type errorList struct {
	tags   []string
	errors []error
}

func ErrorList() *errorList {
	return &errorList{
		tags:   make([]string, 0),
		errors: make([]error, 0),
	}
}

func (e *errorList) Append(tag string, error error) {
	if error == nil {
		return
	}
	e.tags = append(e.tags, tag)
	e.errors = append(e.errors, error)
}

func (e *errorList) Error() string {
	s := ("error list:\n")
	for i, tag := range e.tags {
		s += fmt.Sprintf("%s: %s\n", tag, e.errors[i].Error())
	}
	return s
}
