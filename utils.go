package main

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
