package main
import (
  "github.com/mleef/lpic/docindexing"
  "fmt"
  "os"
)
func main() {
  ind := docindexing.NewIndex()
  ind.AddTerm("marc")
  ind.AddTerm("john")
  ind.AddTerm("aaron")
  
  terms := docindexing.ReadFile(os.Args[1])
  for key, val := range terms {
    fmt.Printf("%s : %d\n", key, val.Frequency)
  }
}
