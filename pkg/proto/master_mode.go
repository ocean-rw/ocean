package proto

import "fmt"

var DefaultMode = NewMode(0, 1, 4, 2)

type Mode uint32

func NewMode(v, r, n, m uint8) Mode {
	mode := uint32(v)<<24 + uint32(r)<<16 + uint32(n)<<8 + uint32(m)
	return Mode(mode)
}

func (m Mode) R() int {
	return int((m & 0xff0000) >> 16)
}

func (m Mode) N() int {
	return int((m & 0xff00) >> 8)
}

func (m Mode) M() int {
	return int(m & 0xff)
}

func (m Mode) String() string {
	return fmt.Sprintf("R%dN%dM%d", m.R(), m.N(), m.M())
}
