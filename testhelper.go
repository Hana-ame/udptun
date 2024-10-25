package main

import "github.com/Hana-ame/udptun/Tools/debug"

type TestHelper string

const testHelper TestHelper = "helper"

func (t TestHelper) ReadChan(rc RecvChannel) {
	c := rc.RecvChan()
	for {
		f := <-c
		s := SprintFrame(f)
		debug.I(t, s)
	}
}
