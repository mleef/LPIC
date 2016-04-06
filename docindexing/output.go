package docindexing

import (
	"fmt"
)


func WriteOutput(filePath string, ind *InvertedIndex) {
	// Print index
	fmt.Printf("# Terms: %d\n\n", ind.TermCount)
	for term, termEntry := range ind.Terms {
		fmt.Printf("Term: %s, Collection Frequency: %d\n", term, termEntry.Frequency)
		for _, docEntry := range termEntry.Documents {
			fmt.Printf("Document Name: %s, Document ID: %d, Document Frequency: %d\n", docEntry.Name, docEntry.ID, docEntry.Frequency)
		}
		fmt.Println()
	}

}