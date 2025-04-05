package main

import "io"

type FramePushHandler interface {
	Push(f Frame) error
}
type FramePollHandler interface {
	Poll() (f Frame, err error)
}

type FramePushCloserHandler interface {
	FramePushHandler
	io.Closer
}

type FramePollCloserHandler interface {
	FramePollHandler
	io.Closer
}

type FrameHandler interface {
	FramePushHandler
	FramePollHandler
	io.Closer
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
