package script

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type (
	Pipe2 interface {
		Exec(cmdS string) Pipe2
		OrPanic() Pipe2
		Match(s string) Pipe2
		ReplaceRegexp(rx *regexp.Regexp, s string) Pipe2
		EachLine(func(string, *strings.Builder)) Pipe2

		ExitStatus() int
		MustStringTrim() string
		Error() error
		SetError(err error)
		ToSlice() ([]string, error)
		WriteFile(string) (int64, error)
		AppendFile(string) (int64, error)
		Stdout() (int, error)
		String() (string, error)
		MustString() string
	}

	Source interface {
		NewPipe() Pipe2
		Echo(msg string) Pipe2
		Exec(cmdS string) Pipe2
		File(fpath string) Pipe2
		Slice([]string) Pipe2
		Slice2(...string) Pipe2
	}
	sourceT struct {
		*Pipe
	}
	sourceMockT struct {
		// *Pipe
		sourceT
	}
)

func NewSource() Source {
	return &sourceT{}
}

func (src *sourceT) NewPipe() Pipe2 {
	// fmt.Fprintf(os.Stderr, "# NewPipe()\n")
	src.Pipe = NewPipe()
	return src
}

func (src *sourceT) Echo(msg string) Pipe2 {
	// fmt.Fprintf(os.Stderr, "# %s\n", cmdS)
	// return Echo(msg)
	src.Pipe = Echo(msg)
	return src
}

func (src *sourceT) EachLine(f func(string, *strings.Builder)) Pipe2 {
	src.Pipe = src.Pipe.EachLine(f)
	return src
}

func (src *sourceT) Exec(cmdS string) Pipe2 {
	// fmt.Fprintf(os.Stderr, "# %s\n", cmdS)
	// return Exec(cmdS)
	src.Pipe = Exec(cmdS)
	return src
}

func (src *sourceT) File(fpath string) Pipe2 {
	// return File(fpath)
	src.Pipe = File(fpath)
	return src
}

func (src *sourceT) ReplaceRegexp(rx *regexp.Regexp, s string) Pipe2 {
	// return Slice(stringL)
	src.Pipe = src.Pipe.ReplaceRegexp(rx, s)
	return src
}

func (src *sourceT) Slice(stringL []string) Pipe2 {
	// return Slice(stringL)
	src.Pipe = Slice(stringL)
	return src
}

func (src *sourceT) ToSlice() ([]string, error) {
	return src.Pipe.Slice()
}

func (src *sourceT) Match(s string) Pipe2 {
	src.Pipe = src.Pipe.Match(s)
	return src
}

func (src *sourceT) OrPanic() Pipe2 {
	src.Pipe = src.Pipe.OrPanic()
	return src
}

func (src *sourceT) Slice2(stringL ...string) Pipe2 {
	return src.Slice(stringL)
	// // return Slice(stringL)
	// src.Pipe = Slice(stringL)
	// return src
}

func NewSourceMock() Source {
	return &sourceMockT{}
}

// func (src *sourceMockT) NewPipe() Pipe2 {
// 	fmt.Fprintf(os.Stderr, "# NewPipe()\n")
// 	// return NewPipe()
// 	src.Pipe = NewPipe()
// 	return src
// }

// func (src *sourceMockT) Echo(msg string) Pipe2 {
// 	// fmt.Fprintf(os.Stderr, "# %s\n", cmdS)
// 	return Echo(msg)
// }

func (src *sourceMockT) Exec(cmdS string) Pipe2 {
	fmt.Fprintf(os.Stderr, "# Exec %s\n", cmdS)
	return src.NewPipe()
}

func (src *sourceMockT) WriteFile(fname string) (int64, error) {
	fmt.Fprintf(os.Stderr, "# WriteFile %q\n", fname)
	return 0, nil
}

func (src *sourceMockT) AppendFile(fname string) (int64, error) {
	fmt.Fprintf(os.Stderr, "# AppendFile %q\n", fname)
	return 0, nil
}

// func (src *sourceMockT) File(fpath string) Pipe2 {
// 	return File(fpath)
// }

// func (src *sourceMockT) Slice(stringL []string) Pipe2 {
// 	return Slice(stringL)
// }

// func (src *sourceMockT) Slice2(stringL ...string) Pipe2 {
// 	return Slice(stringL)
// }

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

func Slice2(stringL ...string) *Pipe {
	return Slice(stringL)
}

// func Echo2f(fmtS string, args ...interface{}) *Pipe {
// 	fmt.Fprintf(os.Stderr, fmtS+"\n", args...)
// 	return NewPipe()
// }
