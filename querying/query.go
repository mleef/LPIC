package querying

import (
	"container/heap"
	"sort"
	"github.com/mleef/lpic/docindexing"
	"log"

)


// To store results of queries
type QueryResult struct {
	Score	float64
	Document	*docindexing.DocumentEntry
}

// An Item is something we manage in a priority queue.
type Item struct {
	value    *docindexing.DocumentEntry 
	sourceList	int
	priority int64    // The priority of the item in the queue.
	index int // The index of the item in the heap.
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
func Query(ind *docindexing.InvertedIndex, queryTerms []string, numResults int) QueryResults {
	log.Printf("Executing with %d term query...", len(queryTerms))
	docList := make([]*docindexing.DocumentEntry, 0)
	numDocs := 0
	i := 0
	
	encountered := map[int64]bool{}
	
	// Collect all document lists from query terms
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

	// Calculate scores for each document
	result := make(QueryResults, numDocs)
	i = 0
	for _, documentEntry := range docList {
		result[i] = scoreDocument(ind, queryTerms, documentEntry)
		i++
	}
	// Sort results by score and return given number of results
	sort.Sort(result)
	if(len(result) > numResults) {
		return result[:numResults]	
	}
	return result
	
}

// Calculate the score for a given document given list of query terms
func scoreDocument(ind *docindexing.InvertedIndex, queryTerms []string, documentEntry *docindexing.DocumentEntry) *QueryResult {
	score := float64(0)
	
	// Check if each term appeared in given document and if so add score
	for _, term := range queryTerms {
		if termEntry, found := ind.TermInDocument(term, documentEntry.ID); found {
			score += TFIDF(ind.DocCount, termEntry, documentEntry)
		}
	}
	
	return &QueryResult{Score: score, Document: documentEntry}
	
}
