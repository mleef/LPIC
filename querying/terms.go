package querying

import (
	"github.com/mleef/lpic/docindexing"
	"math"
	"runtime"
)


// Gets document list associated with term
func GetDocuments(ind *docindexing.InvertedIndex, term string) docindexing.DocumentEntries {
	ind.IndexLock.Lock()
	defer ind.IndexLock.Unlock()
	defer runtime.Gosched()
	
	// Return document list of term if its in the index
	if termEntry, found := ind.Terms[term]; found {
		return termEntry.Documents
	} else {
		return make([]*docindexing.DocumentEntry, 0)
	}
}

// Calculates term frequency inverse document frequency of given term
func TFIDF(docCount int64, termEntry *docindexing.TermEntry, documentEntry *docindexing.DocumentEntry) float64 {
	return rawTermFrequency(documentEntry)*inverseDocumentFrequency(docCount, termEntry)
}

// Inverse document frequency of given term
func inverseDocumentFrequency(docCount int64, termEntry *docindexing.TermEntry) float64 {
	return math.Log(float64(docCount)/float64(len(termEntry.Documents)))
}

// Gets frequency of given term
func rawTermFrequency(documentEntry *docindexing.DocumentEntry) float64 {
		return float64(documentEntry.Frequency)
}

// Binary frequency
func binaryTermFrequency(documentEntry *docindexing.DocumentEntry) int {
	freq := rawTermFrequency(documentEntry)
	if freq > 0 {
		return 1
	} else {
		return 0
	}
}

// Log normalized frequency
func logNormalizedTermFrequency(documentEntry *docindexing.DocumentEntry) float64 {
	freq := rawTermFrequency(documentEntry)
	if freq != 0 {
		return math.Log(freq)
	} else {
		return 0
	}
}