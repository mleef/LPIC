package docindexing

import (
	"os"
	"path/filepath"
	"log"
	"sync"
)

// Data ingested by workers
type Data struct {
	Document string
	ID       int64
}

// Crawl the file system and add data to the work queue
func CrawlFileSystem(workQueue chan *Data, root string, wg *sync.WaitGroup, verbose bool) {
	wg.Add(1)
	defer wg.Done()
	ID := int64(0)

	// Walk file system
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		workQueue <- &Data{path, ID}
		ID++
		return err
	})

	if err != nil && verbose {
		log.Printf("crawl error: %s", err)
	}

	close(workQueue)
}
