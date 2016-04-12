package main

import (
	"log"
	"github.com/mleef/lpic/docindexing"
	"github.com/mleef/lpic/querying"
	"github.com/mleef/lpic/worker"
	"runtime"
	"sync"
	"time"
	"flag"
)

func main() {
	// Optional flags and defaults for index building
    var numWorkers = flag.Int("num-workers", runtime.GOMAXPROCS(runtime.NumCPU()), "number of worker threads")
    var json = flag.Bool("json", false, "generate additional JSON formatted index file")
    var verboseOutput = flag.Bool("verbose", false, "print verbose progress")
    var outputDir = flag.String("out-dir", "./", "destination directory of constructed index")
    var outputFile = flag.String("out-file", "index", "file name of constructed index")
    
    // Optional flags and defaults for index querying
    var numResults = flag.Int("num-results", 5, "number of query results to show")
    var rawTF = flag.Bool("raw-tf", false, "use raw term frequency instead of log norm")
    

	// Get flags
	flag.Parse()
	
	// Validate argument length
	if len(flag.Args()) != 2 {
		log.Fatal("Incorrect number of arguments (2 expected)")
	}
	// Get search starting point from args
	action := flag.Args()[0]
	path := flag.Args()[1]
	
	if action == "build" {
		BuildIndex(path, *numWorkers, *json, *outputDir, *outputFile, *verboseOutput)
	} else if action == "query" {
		QueryIndex(path, *numResults, *rawTF)
	} else {
		log.Fatal("Unknown command")
	}
}

// Build index using given parameters
func BuildIndex(searchPath string, numWorkers int, json bool, outputDir string, outputFile string, verboseOutput bool) {
	
	// To wait on go routines
	var wg sync.WaitGroup
	
	// Initialize needed structures
	log.Println("Initializing index and document pool...")
	ind := docindexing.NewIndex()
	documentPool := make(chan *docindexing.Data)

	// Start timing
	start := time.Now()

	// Commence crawling and index construction
	log.Println("Beginning crawl and index construction...")
	wg.Add(numWorkers + 1)
	go worker.SpawnWorkers(documentPool, numWorkers, ind, &wg, verboseOutput)
	go docindexing.CrawlFileSystem(documentPool, searchPath, &wg, verboseOutput)

	// Wait until all goroutines finish
	wg.Wait()
	
	// Calculate time elapsed
	log.Printf("Finished building index in %s", time.Since(start))
	
	// Sort document lists for querying
	log.Printf("Sorting document lists for query optimization...")
	start = time.Now()
	ind.SortDocumentLists()
	log.Printf("Finished sorting documents in %s", time.Since(start))
	
	// Write output in .lpic format by default
	start = time.Now()
	docindexing.WriteOutput(outputDir + outputFile + ".lpic", ind, false)
	log.Printf("Finished writing output in %s", time.Since(start))
	
	
	// Write index to file in JSON format
	if(json) {
		start = time.Now()
		docindexing.WriteOutput(outputDir + outputFile + ".json", ind, true)
		log.Printf("Finished writing output in %s", time.Since(start))
	}
}

// Start query session using given parameters
func QueryIndex(filePath string, numResults int, rawTF bool) {
	ind := docindexing.ReadOutput(filePath)
	if ind == nil {
		log.Fatal("Error building index")
	} else {
		querying.InteractiveSearch(ind, numResults, rawTF)
	}
}