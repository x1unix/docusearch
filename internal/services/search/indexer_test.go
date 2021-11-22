package search

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/x1unix/docusearch/internal/utils/collections"
)

func TestWordsFromString(t *testing.T) {
	cases := map[string]struct {
		want       []string
		ignoreList []string
	}{
		"simple": {
			want:       []string{"quick", "brown", "fox", "jumps", "over", "lazy", "dog"},
			ignoreList: []string{"the"},
		},
		"long": {
			want: []string{
				"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit", "praesent",
				"malesuada", "nunc", "non", "purus", "hendrerit", "dictum", "tortor", "vivamus", "commodo",
				"urna", "consequat", "dapibus", "donec", "efficitur", "nisl", "justo", "vulputate", "nam",
				"porttitor", "finibus", "quam", "vel", "suscipit", "pellentesque", "aliquam", "dignissim",
				"scelerisque", "fusce", "bibendum", "rutrum", "sapien", "lacinia", "cursus", "luctus",
				"auctor", "lacus", "viverra", "nulla", "mauris", "massa", "eros", "interdum", "imperdiet",
				"velit", "ullamcorper", "tempor", "risus", "quis", "euismod", "aliquet", "condimentum",
				"augue", "sollicitudin", "nec", "faucibus", "ante",
			},
			ignoreList: []string{"et", "sed", "ut", "a", "in", "eu", "id", "mi"},
		},
		"english-verbs": {
			want: []string{
				"you", "like", "sandwich", "i", "clever", "not", "there", "yet",
				"but", "much", "honest", "work", "belongs", "people", "republic",
			},
			ignoreList: EnglishCommonVerbs.ToArray(),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			data, err := ioutil.ReadFile(filepath.Join("testdata", name+".txt"))
			require.NoError(t, err, "failed to read fixture")

			got := WordsFromString(string(data), collections.NewStringsSet(c.ignoreList...))
			require.ElementsMatch(t, c.want, got)
		})
	}
}
