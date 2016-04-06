package docindexing

import (
	"fmt"
	"bytes"
	"os"
	"log"
	"strings"
)

// Writes serialized index to a specified file
func WriteOutput(filePath string, ind *InvertedIndex) {
	f, err := os.Create(filePath)
    if err != nil {
    	log.Printf("error writing file: %s", err)
    } else {
    	defer f.Close()
    	f.WriteString(toJson(ind))
    }

}

// Serializes index to JSON format
func toJson(ind *InvertedIndex) string {
	var buffer bytes.Buffer
	buffer.WriteString("{\n")
	for term, termEntry := range ind.Terms {
		buffer.WriteString(fmt.Sprintf("\t%q : {\n", term))
		buffer.WriteString(fmt.Sprintf("\t\t%q : %d,\n", "frequency", termEntry.Frequency))
		buffer.WriteString(fmt.Sprintf("\t\t%q : [\n", "documents"))
		for numDoc, docEntry := range termEntry.Documents {
			buffer.WriteString("\t\t\t{\n")
			buffer.WriteString(fmt.Sprintf("\t\t\t\t%q : %d,\n", "id", docEntry.ID))
			buffer.WriteString(fmt.Sprintf("\t\t\t\t%q : %q,\n", "name", docEntry.Name))
			buffer.WriteString(fmt.Sprintf("\t\t\t\t%q : %q,\n", "path", docEntry.Path))
			buffer.WriteString(fmt.Sprintf("\t\t\t\t%q : %d,\n", "frequency", docEntry.Frequency))
			buffer.WriteString(fmt.Sprintf("\t\t\t\t%q : [\n", "positions"))
			for index, pos := range docEntry.Positions {
				if index != len(docEntry.Positions) - 1 {
					buffer.WriteString(fmt.Sprintf("\t\t\t\t\t%d,\n", pos))
				} else {
					buffer.WriteString(fmt.Sprintf("\t\t\t\t\t%d\n", pos))
				}
			}
			buffer.WriteString(fmt.Sprintf("\t\t\t\t]\n"))
			
			if numDoc != len(termEntry.Documents) - 1 {
				buffer.WriteString("\t\t\t},\n")
			} else {
				buffer.WriteString("\t\t\t}\n")
			}
		}
		buffer.WriteString(fmt.Sprintf("\t\t]\n"))
		buffer.WriteString("\t},\n")
	}
	result := strings.TrimSuffix(buffer.String(), ",\n")
	result += "\n}\n"
	
	return result

}