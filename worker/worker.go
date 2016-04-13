package worker

import (
	"log"
	"github.com/mleef/lpic/docindexing"
	"sync"
)

// Determines degree of multi-threading and creates workers
func SpawnWorkers(workQueue chan *docindexing.Data, numWorkers int, ind *docindexing.InvertedIndex, wg *sync.WaitGroup, verbose bool) {

	// spawn workers
	for i := 0; i < numWorkers; i++ {
		if verbose {
			log.Printf("Spawning worker %d...\n", i)
		}
		go Worker(i, workQueue, ind, wg, verbose)
	}
}

// Parses documents and updates index
func Worker(id int, workQueue chan *docindexing.Data, ind *docindexing.InvertedIndex, wg *sync.WaitGroup, verbose bool) {
	defer wg.Done()
	documentsIndexed := 0
	
	// While channel is open consume data
	for data := range workQueue {
		if verbose {
			log.Printf("Worker #%d: item %s\n", id, data.Document)
		}
		
		// Parse file and check for errors
		result, err := docindexing.ReadFile(data.Document, data.ID, verbose)
		if err != nil {
			if verbose {
				log.Printf("Worker #%d error: %s\n", id, err)
			}
			continue
		}

		// Increment work count
		documentsIndexed++

		// File parsing was successful so update index
		for term, doc := range result {
			ind.AddDocument(term, doc, verbose)
		}
	}
	
	log.Printf("Worker #%d indexed %d documents\n", id, documentsIndexed)
}
