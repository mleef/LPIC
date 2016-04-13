package querying

import (
	"github.com/mleef/lpic/docindexing"
	"fmt"
	"os"
	"bufio"
	"strings"
	"time"
)

// Shell like interaction for querying index
func InteractiveSearch(ind *docindexing.InvertedIndex, numWorkers int, numResults int, rawTF bool) {
	reader := bufio.NewReader(os.Stdin)
	// Read user input and execute given queries
	var results QueryResults
	for {
		fmt.Print("lpic-query$ ")
		query, _ := reader.ReadString('\n')
		query = query[:len(query) - 1]
		if query == "-quit" || query == "-q" || query == "-exit" {
			fmt.Println("Goodbye!")
			os.Exit(0)
		} else {
			querySplit := strings.Split(query, " ")
			// Start timing
			start := time.Now()
			if len(querySplit) == 0 {
				singleQuery := make([]string, 1)
				singleQuery[0] = query
				results = Query(ind, singleQuery, numWorkers, numResults, rawTF)
			} else {
				results = Query(ind, querySplit, numWorkers, numResults, rawTF)
			}
			fmt.Printf("Query time: %s\n", time.Since(start))
		}
		
		// Print results
		if len(results) == 0 {
			fmt.Printf("No matching documents found\n")
		} else {
			i := 1
			for _, result := range results {
				fmt.Printf("(%d) %s: %f\n", i, result.Document.Path, result.Score)
				i++
			}
		}
		
	}
}