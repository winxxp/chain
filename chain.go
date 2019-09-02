package chain

type Handle func() error

func Run(h ...Handle) error {
	return New().Handles(h...).Run()
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
