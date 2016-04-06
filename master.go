package main

import (
	"fmt"
	"github.com/mleef/lpic/docindexing"
	"github.com/mleef/lpic/worker"
	"runtime"
	"sync"
	"time"
	"flag"
)

func main() {
    var numWorkers = flag.Int("num-workers", runtime.GOMAXPROCS(runtime.NumCPU()), "number of worker threads")
    var verboseOutput = flag.Bool("verbose", false, "print verbose progress")
    var outputDir = flag.String("out-dir", ".", "destination directory for constructed index")
    var outputFile = flag.String("out-file", "index.json", "file name of constructed index")


	flag.Parse()
	
	var wg sync.WaitGroup
	// Initialize needed structures
	fmt.Println("Initializing index and document pool...")
	ind := docindexing.NewIndex()
	documentPool := make(chan *docindexing.Data)

	// Get search starting point from args
	searchPath := flag.Args()[0]

	// Commence crawling and index construction
	fmt.Println("Beginning crawl and index construction...")
	go worker.SpawnWorkers(documentPool, *numWorkers, ind, &wg, *verboseOutput)
	go docindexing.CrawlFileSystem(documentPool, searchPath, &wg, *verboseOutput)

	// Allow goroutines to start
	time.Sleep(100 * time.Millisecond)

	// Wait until all goroutines finish
	wg.Wait()
	
	docindexing.WriteOutput(*outputDir + *outputFile, ind)

}
