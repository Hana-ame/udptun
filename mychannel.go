package main

type RecvChannel interface {
	RecvChan() chan MyFrame
}

type SendChannel interface {
	SendChan() chan MyFrame
}

type RawChannel chan MyFrame

func (c RawChannel) RecvChan() chan MyFrame {
	return c
}
func (c RawChannel) SendChan() chan MyFrame {
	return c
}
func (c RawChannel) Chan() chan MyFrame {
	return c
}

type MergeChannel RawChannel

func (c MergeChannel) RecvChan() chan MyFrame {
	return c
}

// 从底层开始close
func NewMergeChannel(sc0, sc1 RawChannel) MergeChannel {
	rc := make(MergeChannel)
	go func() {
		defer CloseCh(sc0)
		defer CloseCh(sc1)
		defer CloseCh(rc)
		for {
			select {
			case f, ok := <-sc0:
				if !ok {
					return
				}
				rc <- f
			case f, ok := <-sc1:
				if !ok {
					return
				}
				rc <- f
			}
		}
	}()
	return rc
}

func CloseCh(c chan MyFrame) {
	defer func() {
		recover()
	}()
	close(c)
}
