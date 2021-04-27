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

// SendSamples sends the sample stream over the channel ch using the driver from ctx.
func (ctx *Context) SendSamples(ch int, in <-chan Sample) (n int64, err error) {
	return ctx.Driver.SendSamples(ch, in)
}
