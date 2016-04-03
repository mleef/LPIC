package docindexing

import (
	"container/list"
)

// Highest level inverted index structure
type InvertedIndex struct {
	Terms map[string]*TermEntry
	TermCount int
}

// Stores terms and their respective document entries
type TermEntry struct {
	Term string
	Frequency int
	Documents *list.List
}

// Stores document metadata as well as additional term information
type DocumentEntry struct {
	Name string
	Path string
	Size int
	Frequency int
	Positions *list.List
}

// Constructs new inverted index
func NewIndex() *InvertedIndex {
	return &InvertedIndex{Terms: make(map[string]*TermEntry), TermCount: 0}
}

// Adds new term entry to the index
func (ind *InvertedIndex) AddTerm(term string) {
	ind.Terms[term] = &TermEntry{Term: term, Frequency: 0, Documents: list.New()}
	ind.TermCount++;
}

// Adds new document entry to given term entry
func (ind *InvertedIndex) AddDocument(term string, document *DocumentEntry) {
	if _, found := ind.Terms[term]; found {
		ind.Terms[term].Frequency += document.Frequency
		ind.Terms[term].Documents.PushBack(document)
	}
}