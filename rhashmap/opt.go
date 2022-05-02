package rhashmap

type Option interface {
	apply(*config)
}

type hashFunc func(str string) uint64

func (h hashFunc) apply(c *config) {
	c.hashFunc = h
}

func WithHashFunc(hfunc func(str string) uint64) Option {
	return hashFunc(hfunc)
}

type withCap int

func (wc withCap) apply(c *config) {
	c.cap = int(wc)
}

func WithCap(cap int) Option {
	return withCap(cap)
}
