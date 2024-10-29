package param

type Param struct {
	Name     string
	Values   []int
	Iterator int
}

func New(name string, values []int) *Param {
	return &Param{
		Name:     name,
		Values:   values,
		Iterator: 0,
	}
}

func (p *Param) Next() int {
	value := p.Current()
	p.Iterator++
	return value
}

func (p *Param) Current() int {
	return p.Values[p.Iterator%len(p.Values)]
}
