package handler

type Breadcrumb struct {
	crumbs []*Crumb
}

func NewBreadcrumb() *Breadcrumb {
	b := &Breadcrumb{
		crumbs: make([]*Crumb, 0),
	}
	return b
}

func (b *Breadcrumb) Add(url, name string) {
	c := &Crumb{
		Url:  url,
		Name: name,
	}
	b.crumbs = append(b.crumbs, c)
}

func (b *Breadcrumb) Crumbs() []*Crumb {
	return b.crumbs
}

type Crumb struct {
	Url  string
	Name string
}
