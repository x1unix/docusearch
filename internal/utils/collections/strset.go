package collections

// StringsSet is a list of unique strings.
type StringsSet map[string]struct{}

// Append appends value to a set.
func (s StringsSet) Append(vals ...string) {
	for _, v := range vals {
		s[v] = struct{}{}
	}
}

// Slice returns all values as slice of strings.
func (s StringsSet) Slice() []string {
	out := make([]string, 0, len(s))

	for str := range s {
		out = append(out, str)
	}

	return out
}
