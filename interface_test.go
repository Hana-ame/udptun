package main

import "testing"

func f(h FrameHandler) {

}

func TestInterface(t *testing.T) {
	i := FrameHandlerInterface{}

	f(i)

	b := FrameHandlerInterface{
		push:  i.push,
		poll:  i.poll,
		close: i.close,
	}

	f(b)
}
