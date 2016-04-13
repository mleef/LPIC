package querying

import (
	"sort"
	"github.com/mleef/lpic/docindexing"
	"fmt"
)


// To store results of queries
type QueryResult struct {
	Score	float64
	Document	*docindexing.DocumentEntry
}

// Query result slice type for sorting
type QueryResults []*QueryResult

// Implement sorting interface for query result structs
func (slice QueryResults) Len() int {
    return len(slice)
}

func (slice QueryResults) Less(i, j int) bool {
    return slice[i].Score > slice[j].Score;
}

func (slice QueryResults) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

// Queries index with given query terms
func Query(ind *docindexing.InvertedIndex, queryTerms []string, numWorkers int, numResults int, rawTF bool) QueryResults {
	docList := make([]*docindexing.DocumentEntry, 0)
	numDocs := 0
	
	encountered := map[int64]bool{}
	
	// Collect all document lists from query terms and filter out duplicates
	for _, term := range queryTerms {
		newDocs := GetDocuments(ind, term)
		if len(newDocs) > 0 {
			for _, entry := range newDocs {
				if !encountered[entry.ID] {
					docList = append(docList, entry)
					encountered[entry.ID] = true
					numDocs++
				}
			}
		}		
	}
	
	fmt.Printf("Scoring %d matching documents...\n", numDocs)

	// To dispatch work
	low := 0
	step := numDocs/numWorkers
	high := low + step
	resultChan := make(chan *QueryResult, numDocs)
	results := make(QueryResults, numDocs)
	
	// Spawn workers for concurrent scoring
	for i := 0; i < numWorkers; i++ {
		if i == numWorkers - 1 {
			high = len(docList)
		}
		go scoreDocument(ind, queryTerms, resultChan, docList[low:high], rawTF)
		low = high
		high += step
	}

	// Get results from workers
	for i := 0; i < numDocs; i++ {
		results[i] = <- resultChan
	}

	// Sort results by score and return given number of results
	sort.Sort(results)
	if(len(results) > numResults) {
		return results[:numResults]	
	}
	
	return results
	
}

// Calculate the score for a given document given list of query terms
func scoreDocument(ind *docindexing.InvertedIndex, queryTerms []string, resultChan chan *QueryResult, documentEntries []*docindexing.DocumentEntry, rawTF bool) {	
	for _, documentEntry := range documentEntries {
		score := float64(0)
		// Check if each term appeared in given document and if so add score
		for _, term := range queryTerms {
			if termEntry, found := ind.TermInDocument(term, documentEntry.ID); found {
				score += TFIDF(ind.DocCount, termEntry, documentEntry, rawTF)
			}
		}
		resultChan <- &QueryResult{Score: score, Document: documentEntry}
	}
}
