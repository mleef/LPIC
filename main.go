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
  
  terms, err := docindexing.ReadFile(os.Args[1])
  if err != nil {
    fmt.Println(err)
  } else {
    for key, val := range terms {
      fmt.Printf("%s : %d\n", key, val.Frequency)
    }
  }
}
