package main

import (
	"fmt"
	"github.com/mleef/lpic/docindexing"
	"github.com/mleef/lpic/worker"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	// Initialize needed structures
	fmt.Println("Initializing index and document pool...")
	ind := docindexing.NewIndex()
	documentPool := make(chan *docindexing.Data)

	// Get search starting point from args
	searchPath := os.Args[1]

	// Commence crawling and index construction
	fmt.Println("Beginning crawl and index construction...")
	go worker.SpawnWorkers(documentPool, ind, &wg)
	go docindexing.CrawlFileSystem(documentPool, searchPath, &wg)

	// Allow goroutines to start
	time.Sleep(100 * time.Millisecond)

	// Wait until all goroutines finish
	wg.Wait()

	// Print index
	fmt.Printf("# Terms: %d\n", ind.TermCount)
	for term, termEntry := range ind.Terms {
		fmt.Printf("Term: %s, Collection Frequency: %d\n", term, termEntry.Frequency)
		for _, docEntry := range termEntry.Documents {
			fmt.Printf("Document Name: %s, Document ID: %d, Document Frequency: %d\n", docEntry.Name, docEntry.ID, docEntry.Frequency)
		}
		fmt.Println()
	}

}
