package modular

type Options struct {
	SampleRate       int
	BufferSize       int
	DriverBufferSize int
	SampleSize       int
}

type Context struct {
	Options
	Driver Driver
}

// Send sends the sample stream over the channel ch using the driver from ctx.
func (ctx *Context) Send(ch int, in <-chan V) (n int64, err error) {
	return ctx.Driver.Send(ch, in)
}

// SendReader sends the sample stream over the channel ch using the driver from ctx.
func (ctx *Context) SendReader(ch int, r Reader) (n int64, err error) {
	return ctx.Driver.SendReader(ch, r)
}
