package querying

import (
	"github.com/mleef/lpic/docindexing"
	"math"
)


// Gets document list associated with term
func GetDocuments(ind *docindexing.InvertedIndex, term string) docindexing.DocumentEntries {
	ind.IndexLock.Lock()
	defer ind.IndexLock.Unlock()
	
	// Return document list of term if its in the index
	if termEntry, found := ind.Terms[term]; found {
		return termEntry.Documents
	} else {
		return make([]*docindexing.DocumentEntry, 0)
	}
}

// Calculates term frequency inverse document frequency of given term
func TFIDF(ind *docindexing.InvertedIndex, term string) float64 {
	return rawTermFrequency(ind, term)*inverseDocumentFrequency(ind, term)
}

// Inverse document frequency of given term
func inverseDocumentFrequency(ind *docindexing.InvertedIndex, term string) float64 {
	ind.IndexLock.Lock()
	defer ind.IndexLock.Unlock()
		
	// Return document list of term if its in the index
	if termEntry, found := ind.Terms[term]; found {
		return math.Log(float64(ind.DocCount)/float64(len(termEntry.Documents)))
	} else {
		return 0.0
	}
}

// Gets frequency of given term
func rawTermFrequency(ind *docindexing.InvertedIndex, term string) float64 {
	ind.IndexLock.Lock()
	defer ind.IndexLock.Unlock()
	
	// Return document list of term if its in the index
	if termEntry, found := ind.Terms[term]; found {
		return float64(termEntry.Frequency)
	} else {
		return 0
	}
}

// Binary frequency
func binaryTermFrequency(ind *docindexing.InvertedIndex, term string) int {
	freq := rawTermFrequency(ind, term)
	if freq > 0 {
		return 1
	} else {
		return 0
	}
}

// Log normalized frequency
func logNormalizedTermFrequency(ind *docindexing.InvertedIndex, term string) float64 {
	freq := rawTermFrequency(ind, term)
	if freq != 0 {
		return math.Log(freq)
	} else {
		return 0
	}
}