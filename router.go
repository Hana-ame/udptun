package main

import (
	"sort"
	"time"

	tools "github.com/Hana-ame/udptun/Tools"
)

type AddrSegment struct {
	addr Addr
	mask uint32
}

func (s AddrSegment) Compare(e AddrSegment) int {
	ans := int(s.addr) - int(e.addr)
	// if ans == 0 {
	// 	return int(s.mask) - int(e.mask)
	// }
	return ans
}

type AddrSegments []AddrSegment

func (s AddrSegments) Sort() {
	sort.Slice(s[:], func(i, j int) bool {
		if s[i].addr == s[j].addr {
			return s[i].mask < s[j].mask
		}
		return s[i].addr < s[j].addr
	})

}
func (s AddrSegments) BinarySearch(target Addr) int {
	// Inlining is faster than calling BinarySearchFunc with a lambda.
	n := len(s)
	// Define x[-1] < target and x[n] >= target.
	// Invariant: x[i-1] < target, x[j] >= target.
	i, j := 0, n-1
	for i < j {
		h := int(uint(i+j+1) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if int(s[h].addr)-int(target) > 0 {
			j = h - 1
		} else {
			i = h
		}
	}
	// i == j, x[i-1] < target, and x[j] (= x[i]) >= target  =>  answer is i.
	return j
}

type Router struct {
	Chan
	AddrSegments
	*tools.ConcurrentHashMap[Addr, Chan]
}

// defaultAddr 推荐 0
func NewRouter(defaultAddr Addr, defaultChan Chan) *Router {
	return &Router{
		Chan:              defaultChan,
		AddrSegments:      AddrSegments{{defaultAddr, 0}},
		ConcurrentHashMap: tools.NewConcurrentHashMap[Addr, Chan]().Put(defaultAddr, defaultChan),
	}
}

// TODO
func (r *Router) Query(addr Addr) Chan {
	// 先查找路由表
	seg := r.AddrSegments[r.BinarySearch(addr)]
	if (addr & Addr(seg.mask)) == seg.addr {
		return r.GetOrDefault(seg.addr, r.Chan)
	}
	return r.Chan
}

func (r *Router) Serve(addr Addr, writeCh, readCh Chan) error {
	// r.AddrSegments = append(r.AddrSegments, seg)
	// r.AddrSegments.Sort()
	r.Put(addr, readCh)
	for {
		f, err := writeCh.ReadFrame()
		if err != nil {
			break
		}
		fc := r.GetOrDefault(f.Dst(), r.Chan)
		if fc == nil {
			continue
		}
		if err := fc.WriteFrame(f); err != nil {
			continue
		}
	}

	// go func() {
	time.Sleep(time.Minute / 2)
	if readCh == r.GetOrDefault(addr, readCh) {
		r.Remove(addr)
	}
	// readCh.Close()  // 不会再写了，所以关掉
	// writeCh.Close() // 不会再写了，所以关掉
	// }()
	// 不会再写了，所以关掉
	return DoubleError(writeCh.Close(), readCh.Close(), "write", "read")
}
