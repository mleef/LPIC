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
	DocCount  int64
	IndexLock *sync.Mutex
}

// Stores document metadata as well as additional term information
type DocumentEntry struct {
	Path      string
	ID        int64
	Frequency int
	Positions []int
}

// Document slice type for sorting
type DocumentEntries []*DocumentEntry

// Stores terms and their respective document entries
type TermEntry struct {
	Frequency int
	Documents DocumentEntries
}

// Constructs new inverted index
func NewIndex() *InvertedIndex {
	return &InvertedIndex{Terms: make(map[string]*TermEntry), TermCount: int64(0), IndexLock: &sync.Mutex{}}
}

// Adds new term entry to the index
func (ind *InvertedIndex) addTerm(term string, verbose bool) {
	// Make sure we aren't overwriting existing term
	if _, found := ind.Terms[term]; !found {
		if verbose {
			log.Printf("Adding new term %s to index", term)
		}
		ind.Terms[term] = &TermEntry{Frequency: 0, Documents: make(DocumentEntries, 0)}
		ind.TermCount++
	}
}

// Adds new document entry to given term entry
func (ind *InvertedIndex) addDocument(term string, document *DocumentEntry, verbose bool) {
	// Make sure term is in index
	ind.addTerm(term, verbose)
	if verbose {
		log.Printf("Adding new document %s to term %s's document list", document.Path, term)
	}
	
	// Safely update values
	ind.Terms[term].Frequency += document.Frequency
	ind.Terms[term].Documents = append(ind.Terms[term].Documents, document)
	ind.DocCount++
}


// Public method for adding terms
func (ind *InvertedIndex) AddDocument(term string, document *DocumentEntry, verbose bool) {
	// Get index lock
	ind.IndexLock.Lock()
	defer ind.IndexLock.Unlock()
	defer runtime.Gosched()
	
	ind.addDocument(term, document, verbose)
}

// Public method for adding documents
func (ind *InvertedIndex) AddTerm(term string, verbose bool) {
	// Get index lock
	ind.IndexLock.Lock()
	defer ind.IndexLock.Unlock()
	defer runtime.Gosched()
	
	ind.addTerm(term, verbose)
}

// Adds new document entry to given term entry
func (ind *InvertedIndex) AddDocuments(result map[string]*DocumentEntry, verbose bool) {
	// Get index lock
	ind.IndexLock.Lock()
	defer ind.IndexLock.Unlock()
	defer runtime.Gosched()
	
	for term, doc := range result {
		ind.addDocument(term, doc, verbose)
	}
}

// Checks if a term occurs was found in a given document
func (ind *InvertedIndex) TermInDocument(term string, id int64) (*TermEntry, bool) {
	ind.IndexLock.Lock()
	defer ind.IndexLock.Unlock()
	defer runtime.Gosched()
	if termEntry, found := ind.Terms[term]; found {
		for _, documentEntry := range termEntry.Documents {
			if documentEntry.ID == id {
				return termEntry, true
			}
		}
	}
	
	return nil, false

}
