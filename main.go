package main
import (
  "github.com/mleef/lpic/docindexing"
)
func main() {
	ind := docindexing.NewIndex()
	ind.AddTerm("marc")
	ind.AddTerm("john")
	ind.AddTerm("aaron")
}
