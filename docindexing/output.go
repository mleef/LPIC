package docindexing

import (
	"fmt"
	"bytes"
	"os"
	"log"
	"strings"
	"strconv"
	"runtime"
	"bufio"
)

// Writes serialized index to a specified file
func WriteOutput(filePath string, ind *InvertedIndex, jsonFormat bool) {
	f, err := os.Create(filePath)
    if err != nil {
    	log.Printf("error writing file: %s", err)
    } else {
    	defer f.Close()
    	log.Printf("Writing output file %s", filePath)
    	if jsonFormat {
    		f.WriteString(toJSON(ind))
    	} else {
    		f.WriteString(toLPIC(ind))
    	}
    }

}

// Reads serialized index from a specified file
func ReadOutput(filePath string) *InvertedIndex {
	f, err := os.Open(filePath)
    if err != nil {
    	log.Printf("error reading file: %s", err)
    	return nil
    } else {
    	defer f.Close()
    	log.Printf("Reading output file %s", filePath)
    	return fromLPIC(f)
    }

}

// Serializes index to LPIC format
func toLPIC(ind *InvertedIndex) string {
	ind.IndexLock.Lock()
	defer ind.IndexLock.Unlock()
	defer runtime.Gosched()
	var buffer bytes.Buffer

	for term, termEntry := range ind.Terms {
		buffer.WriteString(fmt.Sprintf("%s,%d,", term, termEntry.Frequency))
		for index, docEntry := range termEntry.Documents {
			buffer.WriteString(fmt.Sprintf("%d %s %d ", docEntry.ID, docEntry.Path, docEntry.Frequency))
			for posInd, pos := range docEntry.Positions {
				if posInd == len(docEntry.Positions) - 1 {
					buffer.WriteString(fmt.Sprintf("%d", pos))
				} else {
					buffer.WriteString(fmt.Sprintf("%d-", pos))
				}
			}
			if index == len(termEntry.Documents) - 1 {
				buffer.WriteString("\n")
			} else {
				buffer.WriteString(",")
			}
		}
	}

	
	return buffer.String()
}

// Serializes index to JSON format
func toJSON(ind *InvertedIndex) string {
	ind.IndexLock.Lock()
	defer ind.IndexLock.Unlock()
	defer runtime.Gosched()
	var buffer bytes.Buffer
	buffer.WriteString("{\n")
	for term, termEntry := range ind.Terms {
		buffer.WriteString(fmt.Sprintf("\t%q : {\n", term))
		buffer.WriteString(fmt.Sprintf("\t\t%q : %d,\n", "frequency", termEntry.Frequency))
		buffer.WriteString(fmt.Sprintf("\t\t%q : [\n", "documents"))
		for numDoc, docEntry := range termEntry.Documents {
			buffer.WriteString("\t\t\t{\n")
			buffer.WriteString(fmt.Sprintf("\t\t\t\t%q : %d,\n", "id", docEntry.ID))
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

// Read in LPIC format and return constructed index
func fromLPIC(file *os.File) *InvertedIndex {
	ind := NewIndex()
	reader := bufio.NewReader(file)
	rawLine, err := reader.ReadString('\n')
    for err == nil {
        line := strings.Split(rawLine, ",")
        term := line[0]
        ind.AddTerm(term, false)
        for _, document := range line[2:] {
        	docSplit := strings.Split(document, " ")
        	positions := make([]int, 0)
        	for _, pos := range strings.Split(docSplit[3], "-") {
        		conv, err := strconv.Atoi(pos)
        		if err == nil {
        			positions = append(positions, conv)
        		}
        	}
        	
        	path := docSplit[1]
        	id, err1 := strconv.Atoi(docSplit[0])
        	frequency, err2 := strconv.Atoi(docSplit[2])
        	
        	if err1 == nil && err2 == nil {
        		ind.AddDocument(term, &DocumentEntry{path, int64(id), frequency, positions}, false)
        	}
        }
        rawLine, err = reader.ReadString('\n')
    }
	
	return ind
}