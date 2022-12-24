package main

const tagLength = 2

type PortalBuf []byte

func (p PortalBuf) Tag() PortalBuf {
	if len(p) >= tagLength {
		return (PortalBuf)(p)[:tagLength]
	}
	return nil
}

// add tag to PortalBuf
func (p PortalBuf) AddTag(tag any) PortalBuf {
	if len(p) >= tagLength {
		if t, ok := tag.(int); ok {
			(p)[0] = byte(t % 256)
			(p)[1] = byte(t / 256)
			return (PortalBuf)(p)
		}
	}
	return nil
}

// n <= 0 return all Data part
// n > 0 return Data part length=n
func (p PortalBuf) Data(n int) PortalBuf {
	if len(p) >= tagLength {
		if n > 0 {
			return (PortalBuf)(p)[tagLength : tagLength+n]
		} else {
			return (PortalBuf)(p)[tagLength:]
		}
	}
	return nil
}

// length is n
func (p PortalBuf) Raw(n int) PortalBuf {
	return (p)[:n]
}

// length is tagLength+n
func (p PortalBuf) DataAndTag(n int) PortalBuf {
	return (p)[:n+tagLength]
}
