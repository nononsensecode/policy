package policy

type Filter struct {
	like   []string
	equals []string
}

func (f Filter) Like() []string {
	return f.like
}

func (f Filter) Equals() []string {
	return f.equals
}
