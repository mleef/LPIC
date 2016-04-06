package worker

import (
	//"fmt"
	"github.com/mleef/lpic/docindexing"
	"runtime"
	"sync"
)

// Determines degree of multi-threading and creates workers
func SpawnWorkers(workQueue chan *docindexing.Data, ind *docindexing.InvertedIndex, wg *sync.WaitGroup) {
	ncpu := runtime.NumCPU()
	runtime.GOMAXPROCS(ncpu) //max # of threads that can be running

	// spawn workers
	for i := 0; i < ncpu; i++ {
		wg.Add(1)
		//fmt.Printf("Spawning worker %d...\n", i)
		go Worker(i, workQueue, ind, wg)
	}
}

// Parses documents and updates index
func Worker(id int, workQueue chan *docindexing.Data, ind *docindexing.InvertedIndex, wg *sync.WaitGroup) {
	defer wg.Done()

	// While channel is open consume data
	for data := range workQueue {
		//fmt.Printf("worker #%d: item %v\n", id, *data)

		// Parse file and check for errors
		result, err := docindexing.ReadFile(data.Document, data.ID)
		if err != nil {
			//fmt.Printf("worker #%d error: %s\n", id, err)
			continue
		}

		// File parsing was successful so update index
		for term, doc := range result {
			ind.AddDocument(term, doc)
		}
	}
}
