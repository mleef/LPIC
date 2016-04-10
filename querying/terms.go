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
	return float64(rawFrequency(ind, term))/float64(ind.DocCount)
}

// Gets frequency of given term
func rawFrequency(ind *docindexing.InvertedIndex, term string) int {
	ind.IndexLock.Lock()
	defer ind.IndexLock.Unlock()
	
	// Return document list of term if its in the index
	if termEntry, found := ind.Terms[term]; found {
		return termEntry.Frequency
	} else {
		return 0
	}
}

// Binary frequency
func binaryFrequency(ind *docindexing.InvertedIndex, term string) int {
	freq := rawFrequency(ind, term)
	if freq > 0 {
		return 1
	} else {
		return 0
	}
}

// Log normalized frequency
func logNormalizedFrequency(ind *docindexing.InvertedIndex, term string) float64 {
	freq := rawFrequency(ind, term)
	if freq != 0 {
		return math.Log(float64(freq))
	} else {
		return 0
	}
}