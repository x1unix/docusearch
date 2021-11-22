package collections

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStringsSet(t *testing.T) {
	want := StringsSet{
		"foo": void{},
		"bar": void{},
	}

	got := NewStringsSet("foo", "bar", "foo")
	require.Equal(t, want, got)
}

func TestStringsSet_Append(t *testing.T) {
	m := NewStringsSet("foo")
	m.Append("foo", "bar")
	require.ElementsMatch(t, []string{"foo", "bar"}, m.ToArray())
}

func TestStringsSet_Has(t *testing.T) {
	m := NewStringsSet("foo")
	require.True(t, m.Has("foo"))
}
