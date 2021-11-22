package search

import (
	"strings"
	"unicode"

	"github.com/x1unix/docusearch/internal/utils/collections"
)

// EnglishCommonVerbs is collection of common articles, auxiliary verbs and
// other elements of English language that used only to carry semantic load.
var EnglishCommonVerbs = collections.NewStringsSet(
	"the", "to", "of", "or", "a", "an", "in", "is", "isn",
	"t", "doesn", "ain", "didn", "did", "was", "wasn",
	"were", "weren", "would", "wouldn", "m", "am", "be",
	"of", "that", "in", "on", "are", "aren", "or",
)

func tokenizeText(str string) []string {
	return strings.FieldsFunc(str, func(r rune) bool {
		// TODO: properly handle contractions such as "aren't", "ain't", "doesn't", etc.
		return !unicode.IsNumber(r) && !unicode.IsLetter(r)
	})
}

// WordsFromString returns a list of unique words from string text.
//
// Second parameter allows specifying ignore list to filter common verbs, articles, etc.
func WordsFromString(str string, ignoreList collections.StringsSet) []string {
	allWords := tokenizeText(str)
	uniqueWords := make(collections.StringsSet)
	for _, word := range allWords {
		lowered := strings.ToLower(word)
		if ignoreList.Has(lowered) {
			continue
		}

		uniqueWords.Append(lowered)
	}
	return uniqueWords.ToArray()
}
