package docindexing

import (
	"container/list"
	"sync"
	"runtime"
)

// Highest level inverted index structure
type InvertedIndex struct {
	Terms map[string]*TermEntry
	TermCount int
	IndexLock *sync.Mutex
	
}

// Stores terms and their respective document entries
type TermEntry struct {
	Term string
	Frequency int
	Documents *list.List
	EntryLock *sync.Mutex
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
	return &InvertedIndex{Terms: make(map[string]*TermEntry), TermCount: 0, IndexLock: &sync.Mutex{}}
}

// Adds new term entry to the index
func (ind *InvertedIndex) AddTerm(term string) {
	ind.IndexLock.Lock()
	// Make sure we aren't overwriting existing term
	if _, found := ind.Terms[term]; !found {
		ind.Terms[term] = &TermEntry{Term: term, Frequency: 0, Documents: list.New(), EntryLock: &sync.Mutex{}}
		ind.TermCount++;
	}
	ind.IndexLock.Unlock()
	runtime.Gosched()
}

// Adds new document entry to given term entry
func (ind *InvertedIndex) AddDocument(term string, document *DocumentEntry) {
	// Make sure term is in index
	ind.AddTerm(term)
	
	// Safely add document to term list
	ind.Terms[term].EntryLock.Lock()
	ind.Terms[term].Frequency += document.Frequency
	ind.Terms[term].Documents.PushBack(document)
	ind.Terms[term].EntryLock.Unlock()
	runtime.Gosched()
}