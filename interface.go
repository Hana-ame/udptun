package main

type FrameHandler interface {
	Push(f Frame) error
	Poll() (f Frame, err error)

	Close() error
}

type FrameHandlerInterface struct {
	push func(f Frame) error
	poll func() (f Frame, err error)

	close func() error
}

func (i FrameHandlerInterface) Push(f Frame) error {
	return i.push(f)
}
func (i FrameHandlerInterface) Poll() (Frame, error) {
	return i.poll()
}
func (i FrameHandlerInterface) Close() error {
	return i.close()
}
