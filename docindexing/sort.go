package docindexing

import (
	"sort"
)

// Implement sorting interface for document entry structs
func (slice DocumentEntries) Len() int {
    return len(slice)
}

func (slice DocumentEntries) Less(i, j int) bool {
    return slice[i].ID < slice[j].ID;
}

func (slice DocumentEntries) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}


// Sorts every document list in the index by ID for easy querying
func (ind *InvertedIndex) SortDocumentLists() {
	ind.IndexLock.Lock()
	defer ind.IndexLock.Unlock()
	
	for _, termEntry := range ind.Terms {
		sort.Sort(termEntry.Documents)
	}
}