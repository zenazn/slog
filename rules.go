package slog

type rule struct {
	selector string
	level    Level
}

type rules []rule

func (r rules) Len() int {
	return len(r)
}

func (r rules) Less(i, j int) bool {
	// We sort largest to smallest
	return r[i].selector > r[j].selector
}

func (r rules) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
