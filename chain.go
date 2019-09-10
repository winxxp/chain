package chain

import (
	"errors"
	"strconv"
	"strings"
)

type Handle func() error

func Run(h ...Handle) error {
	return New().Handles(h...).Run()
}

func ParallelRun(h ...Handle) error {
	return New().Handles(h...).ParallelRun()
}

type HandleChain struct {
	err        error
	handles    []Handle
	lastHandle func(error)
}

func New() *HandleChain {
	return &HandleChain{
		handles: make([]Handle, 0, 3),
	}
}

func (r *HandleChain) Last(h func(error)) *HandleChain {
	r.lastHandle = h
	return r
}

func (r *HandleChain) Handles(h ...Handle) *HandleChain {
	r.handles = append(r.handles, h...)
	return r
}

func (c *HandleChain) Run() error {
	for _, h := range c.handles {
		if c.err == nil {
			c.err = h()
		}
	}

	if c.lastHandle != nil {
		c.lastHandle(c.err)
	}

	return c.err
}

func (c *HandleChain) ParallelRun() error {
	var (
		errs []error
		ch   = make(chan error, len(c.handles))
	)

	for _, h := range c.handles {
		go func(h Handle) {
			ch <- h()
		}(h)
	}

	for range c.handles {
		errs = append(errs, <-ch)
	}

	sb := &strings.Builder{}
	i := 0
	for _, err := range errs {
		if err != nil {
			sb.WriteString(strconv.Itoa(i))
			sb.WriteString(": ")
			sb.WriteString(err.Error())
			sb.WriteString("\n")
			i++
		}
	}

	if sb.String() != "" {
		c.err = errors.New(sb.String())
	}

	return c.err
}
