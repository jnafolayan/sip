package ezw

import (
	"github.com/jnafolayan/sip/pkg/signal"
)

type Zerotree struct {
	Level    int
	Children []*Zerotree
}

func SignificancePass(s signal.Signal2D, level int) *Zerotree {
	root := &Zerotree{}

	return root
}
