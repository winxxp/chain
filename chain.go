package chain

type Handle func() error

type C struct {
	err     error
	handles []Handle
}

func New() *C {
	return &C{
		handles: make([]Handle, 0, 3),
	}
}

func (r *C) Handles(h ...Handle) *C {
	r.handles = append(r.handles, h...)
	return r
}

func (c *C) Run() error {
	for _, h := range c.handles {
		if c.err == nil {
			c.err = h()
		}
	}

	return c.err
}
