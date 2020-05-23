package chain

type Handle func(ctx *Context) error

func Run(h ...Handle) error {
	return New().Handles(h...).Run()
}

//HandleChain 按循序调用函数
type HandleChain struct {
	err        error
	ctx        *Context
	handles    []Handle
	lastHandle func(ctx *Context, err error)
	eachBefore func(*Context)
	eachAfter  func(*Context, error)
}

func New() *HandleChain {
	return &HandleChain{
		handles: make([]Handle, 0, 3),
		ctx:     &Context{},
	}
}

func (c *HandleChain) Last(h func(ctx *Context, err error)) *HandleChain {
	c.lastHandle = h
	return c
}

func (c *HandleChain) EachBefore(h func(ctx *Context)) *HandleChain {
	c.eachBefore = h
	return c
}

func (c *HandleChain) EachAfter(h func(ctx *Context, err error)) *HandleChain {
	c.eachAfter = h
	return c
}

func (c *HandleChain) Handles(h ...Handle) *HandleChain {
	c.handles = append(c.handles, h...)
	return c
}

func (c *HandleChain) Run() error {
	for _, h := range c.handles {
		if c.err == nil {
			if c.eachBefore != nil {
				c.eachBefore(c.ctx)
			}

			c.err = h(c.ctx)

			if c.eachAfter != nil {
				c.eachAfter(c.ctx, c.err)
			}
		}
	}

	if c.lastHandle != nil {
		c.lastHandle(c.ctx, c.err)
	}

	return c.err
}
