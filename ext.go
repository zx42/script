package script

import (
	"fmt"
	"strings"
)

func (p *Pipe) OrPanic() *Pipe {
	status := p.ExitStatus()
	if status == 0 {
		return p
	}
	p.SetError(nil)
	errOutS, err := p.String()
	if err != nil {
		msg := fmt.Errorf("cmd %q: %w", p.lastCmdLine, err)
		panic(msg)
	}
	msg := fmt.Sprintf("cmd %q failed: %d\n\n%s", p.lastCmdLine, status, errOutS)
	panic(msg)
}

func (p *Pipe) StringWithDefault(defS string) string {
	s, err := p.String()
	if err != nil {
		s = defS
	}
	return s
}

func (p *Pipe) MustString() string {
	return MustS(p.String())
}

func (p *Pipe) MustStringTrim() string {
	return strings.TrimSpace(MustS(p.String()))
}

func MustS(s string, err error) string {
	if err != nil {
		fmt.Println(s)
		panic(err)
	}
	return s
}
