package mymux

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ERR_CLOSED Error = ("closed")
)
