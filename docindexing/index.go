package docindexing

import (
	"log"
	"runtime"
	"sync"
)

// Highest level inverted index structure
type InvertedIndex struct {
	Terms     map[string]*TermEntry
	TermCount int64
	IndexLock *sync.Mutex
}

// Stores terms and their respective document entries
type TermEntry struct {
	Frequency int
	Documents []*DocumentEntry
	EntryLock *sync.Mutex
}

// Stores document metadata as well as additional term information
type DocumentEntry struct {
	Path      string
	ID        int64
	Frequency int
	Positions []int
}

// Constructs new inverted index
func NewIndex() *InvertedIndex {
	return &InvertedIndex{Terms: make(map[string]*TermEntry), TermCount: int64(0), IndexLock: &sync.Mutex{}}
}

// Adds new term entry to the index
func (ind *InvertedIndex) AddTerm(term string, verbose bool) {
	ind.IndexLock.Lock()
	// Make sure we aren't overwriting existing term
	if _, found := ind.Terms[term]; !found {
		if verbose {
			log.Printf("Adding new term %s to index", term)
		}
		ind.Terms[term] = &TermEntry{Frequency: 0, Documents: make([]*DocumentEntry, 0), EntryLock: &sync.Mutex{}}
		ind.TermCount++
	}
	ind.IndexLock.Unlock()
	runtime.Gosched()
}

// Adds new document entry to given term entry
func (ind *InvertedIndex) AddDocument(term string, document *DocumentEntry, verbose bool) {
	// Make sure term is in index
	ind.AddTerm(term, verbose)
	if verbose {
		log.Printf("Adding new document %s to term %s's document list", document.Path, term)
	}
	// Get index and term locks
	ind.IndexLock.Lock();
	ind.Terms[term].EntryLock.Lock()
	
	// Safely update values
	ind.Terms[term].Frequency += document.Frequency
	ind.Terms[term].Documents = append(ind.Terms[term].Documents, document)
	
	// Release locks
	ind.Terms[term].EntryLock.Unlock()
	ind.IndexLock.Unlock()
	runtime.Gosched()
}
