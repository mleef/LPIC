package docindexing

import (
	"os"
	"path/filepath"
	"log"
	"sync"
	"time"
)

// Data ingested by workers
type Data struct {
	Document string
	ID       int64
}

// Crawl the file system and add data to the work queue
func CrawlFileSystem(workQueue chan *Data, root string, wg *sync.WaitGroup, verbose bool) {
	defer wg.Done()
	ID := int64(0)
	
	// Start timing
	start := time.Now()
	
	// Walk file system
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		workQueue <- &Data{path, ID}
		ID++
		return err
	})

	if err != nil && verbose {
		log.Printf("crawl error: %s", err)
	}

	log.Printf("Finished crawling %s in %s", root, time.Since(start))
	close(workQueue)
}
