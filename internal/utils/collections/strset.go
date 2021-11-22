package collections

// StringsSet is a list of unique strings.
type StringsSet map[string]struct{}

// NewStringsSet constructs StringsSet from slice.
func NewStringsSet(vals ...string) StringsSet {
	set := make(StringsSet, len(vals))
	set.Append(vals...)
	return set
}

// Append appends value to a set.
func (s StringsSet) Append(vals ...string) {
	for _, v := range vals {
		s[v] = struct{}{}
	}
}

// Has checks if set contains a value.
func (s StringsSet) Has(val string) bool {
	if len(s) == 0 {
		return false
	}

	_, ok := s[val]
	return ok
}

// ToArray returns all values as list of strings.
func (s StringsSet) ToArray() []string {
	out := make([]string, 0, len(s))

	for str := range s {
		out = append(out, str)
	}

	return out
}
