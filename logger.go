package hego

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

func tabbedJoin(s []string) (res string) {
	for _, val := range s {
		res += val + "\t"
	}
	return
}

type logger struct {
	name    string
	verbose int
	maxIter int
	writer  *tabwriter.Writer
	buf     bytes.Buffer
}

func newLogger(name string, cols []string, verbose, maxIter int) *logger {
	l := logger{}
	l.name = name
	l.writer = tabwriter.NewWriter(&l.buf, 0, 0, 3, []byte(" ")[0], tabwriter.AlignRight)
	l.verbose = verbose
	l.maxIter = maxIter
	fmt.Fprintln(l.writer, tabbedJoin(cols))
	return &l
}

func (l *logger) AddLine(i int, cols []string) {
	if i%l.verbose == 0 || i+1 == l.maxIter {
		fmt.Fprintln(l.writer, tabbedJoin(cols))
	}
}

func (l *logger) Flush() {
	l.writer.Flush()
	if l.verbose > 0 {
		fmt.Println(l.buf.String())
	}
}
