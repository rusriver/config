package config

type Location []string

func (loc *Location) Push(s string) {
	*loc = append(*loc, s)
}

func (loc *Location) Pop() {
	*loc = (*loc)[:len(*loc)-1]
}
