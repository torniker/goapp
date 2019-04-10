package app

type path struct {
	path     string
	segments []string
	index    int
}

func (p path) Next() string {
	if len(p.segments) < p.index+2 {
		return ""
	}
	return p.segments[p.index+1]
}

func (p path) Current() string {
	if len(p.segments) < p.index+1 {
		return ""
	}
	return p.segments[p.index]
}
