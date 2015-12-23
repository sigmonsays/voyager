package crawl

type Crawler struct {
}

func DefaultOptions() *Options {
	return &Options{
		Paths: make([]string, 0),
	}
}

type Options struct {
	Paths []string
}

func (o *Options) AddPath(path string) {
	o.Paths = append(o.Paths, path)
}

func NewCrawler(opts *Options) *Crawler {
	if opts == nil {
		opts = DefaultOptions()
	}

	c := &Crawler{}

	return c
}
