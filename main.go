package main
import (
  "github.com/mleef/lpic/docindexing"
  "fmt"
)
func main() {
	ind := docindexing.NewIndex()
	ind.AddTerm("marc")
	ind.AddTerm("john")
	ind.AddTerm("aaron")
	
	terms := docindexing.ReadFile("/Users/marcleef/Desktop/sample.txt")
	for key, val := range terms {
		fmt.Printf("%s : %d\n", key, val.Frequency)
	}
}
