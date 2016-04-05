package docindexing

import (
    "bufio"
    "log"
    "errors"
    "os"
    "strings"
    "regexp"
    "container/list"
    "unicode"
)

var supportedExtensions = map[string]int{
  "txt" : 0,
  "py" : 0,
  "c" : 0,
  "cpp" : 0,
  "go" : 0,
  "java" : 0,
}

// List of stop words to ignore when parsing
var stopWords = map[string]int{
  "a" : 0,
  "about" : 0,
  "above" : 0,
  "after" : 0,
  "again" : 0,
  "against" : 0,
  "all" : 0,
  "am" : 0,
  "an" : 0,
  "and" : 0,
  "any" : 0,
  "are" : 0,
  "arent" : 0,
  "as" : 0,
  "at" : 0,
  "be" : 0,
  "because" : 0,
  "been" : 0,
  "before" : 0,
  "being" : 0,
  "below" : 0,
  "between" : 0,
  "both" : 0,
  "but" : 0,
  "by" : 0,
  "cant" : 0,
  "cannot" : 0,
  "could" : 0,
  "couldnt" : 0,
  "did" : 0,
  "didnt" : 0,
  "do" : 0,
  "does" : 0,
  "doesnt" : 0,
  "doing" : 0,
  "dont" : 0,
  "down" : 0,
  "during" : 0,
  "each" : 0,
  "few" : 0,
  "for" : 0,
  "from" : 0,
  "further" : 0,
  "had" : 0,
  "hadnt" : 0,
  "has" : 0,
  "hasnt" : 0,
  "have" : 0,
  "havent" : 0,
  "having" : 0,
  "he" : 0,
  "hed" : 0,
  "hell" : 0,
  "hes" : 0,
  "her" : 0,
  "here" : 0,
  "heres" : 0,
  "hers" : 0,
  "herself" : 0,
  "him" : 0,
  "himself" : 0,
  "his" : 0,
  "how" : 0,
  "hows" : 0,
  "i" : 0,
  "id" : 0,
  "ill" : 0,
  "im" : 0,
  "ive" : 0,
  "if" : 0,
  "in" : 0,
  "into" : 0,
  "is" : 0,
  "isnt" : 0,
  "it" : 0,
  "its" : 0,
  "itself" : 0,
  "lets" : 0,
  "me" : 0,
  "more" : 0,
  "most" : 0,
  "mustnt" : 0,
  "my" : 0,
  "myself" : 0,
  "no" : 0,
  "nor" : 0,
  "not" : 0,
  "of" : 0,
  "off" : 0,
  "on" : 0,
  "once" : 0,
  "only" : 0,
  "or" : 0,
  "other" : 0,
  "ought" : 0,
  "our" : 0,
  "ours" : 0,
  "ourselves" : 0,
  "out" : 0,
  "over" : 0,
  "own" : 0,
  "same" : 0,
  "shant" : 0,
  "she" : 0,
  "shed" : 0,
  "shell" : 0,
  "shes" : 0,
  "should" : 0,
  "shouldnt" : 0,
  "so" : 0,
  "some" : 0,
  "such" : 0,
  "than" : 0,
  "that" : 0,
  "thats" : 0,
  "the" : 0,
  "their" : 0,
  "theirs" : 0,
  "them" : 0,
  "themselves" : 0,
  "then" : 0,
  "there" : 0,
  "theres" : 0,
  "these" : 0,
  "they" : 0,
  "theyd" : 0,
  "theyll" : 0,
  "theyre" : 0,
  "theyve" : 0,
  "this" : 0,
  "those" : 0,
  "through" : 0,
  "to" : 0,
  "too" : 0,
  "under" : 0,
  "until" : 0,
  "up" : 0,
  "very" : 0,
  "was" : 0,
  "wasnt" : 0,
  "we" : 0,
  "wed" : 0,
  "well" : 0,
  "weve" : 0,
  "were" : 0,
  "werent" : 0,
  "what" : 0,
  "whats" : 0,
  "when" : 0,
  "whens" : 0,
  "where" : 0,
  "wheres" : 0,
  "which" : 0,
  "while" : 0,
  "who" : 0,
  "whos" : 0,
  "whom" : 0,
  "why" : 0,
  "whys" : 0,
  "with" : 0,
  "wont" : 0,
  "would" : 0,
  "wouldnt" : 0,
  "you" : 0,
  "youd" : 0,
  "youll" : 0,
  "youre" : 0,
  "youve" : 0,
  "your" : 0,
  "yours" : 0,
  "yourself" : 0,
  "yourselves" : 0,
}

// Validates and opens file, returns a scanner to read file
func OpenFile(path string) (*os.File, os.FileInfo, *bufio.Scanner, error) {
  // Open file
  file, err := os.Open(path)
  if err != nil {
    return nil, nil, nil, err
  }

  // Read file line by line
  scanner := bufio.NewScanner(file)

  // Get file metadata 
  fileInfo, err := file.Stat()
  if err != nil {
    return nil, nil, nil, err
  }

  // Validate file type
  if valid, err := validFile(path, fileInfo); !valid {
    return nil, nil, nil, err
  }
  
  return file, fileInfo, scanner, nil
}

// Reads file line by line
func ReadFile(path string) (map[string]*DocumentEntry, error) {
  log.Printf("Working on file %s", path)
  // Open file for reading
  file, fileInfo, scanner, err := OpenFile(path)
  defer file.Close()
  
  // Something went wrong in opening the file
  if err != nil {
    return nil, err
  }
  
  // Get file metadata
  fileName := fileInfo.Name()
  fileSize := fileInfo.Size()
  
  // To store resulting term frequencies
  termCounts := make(map[string]*DocumentEntry)
  position := 0

  // Format each line and then update frequencies
  for scanner.Scan() {
    terms := formatTerms(strings.Fields(scanner.Text()))
    for _, term := range terms {
      if _, present := termCounts[term]; !present {
        termCounts[term] = &DocumentEntry{fileName, path, fileSize, 0, list.New()}
      } else {
        termCounts[term].Frequency++
        termCounts[term].Positions.PushBack(position)        
      }
      position++
    }
  }
  
  // Make sure the scanner didn't fail
  if err := scanner.Err(); err != nil {
    return termCounts, err
  }

  return termCounts, nil
}

// Validates that file is parseable
func validFile(path string, fileInfo os.FileInfo) (bool, error) {
  
  // Check for regular file
  if fileInfo.IsDir() {
      return false, errors.New("Cannot parse directory")
  }

  // Get file extension
  pathSplit := strings.Split(path, ".")
  if len(pathSplit) == 0 {
    return false, errors.New("File has no extension")
  }
  
  // Check that file type is supported
  extension := pathSplit[len(pathSplit) - 1]
  if _,present := supportedExtensions[extension]; !present {
    return false, errors.New("Unsupported file type")
  }
  
  return true, nil
}

// Formats words for term construction
func formatTerms(words []string) []string {
    reg, err := regexp.Compile("[^A-Za-z]+")
    if err != nil {
        panic(err)
    }
    
  // To store formatted strings
  newWords := make([]string, 0)
  
  for _, word := range words {
    // Check if a stop word
    if _, inMap := stopWords[strings.ToLower(word)]; !inMap {
      // Strip non alphabetic characters and spaces
      newWord := clearSpaces(reg.ReplaceAllString(word, ""))
      
      // Make sure there are remaining characters
      if(len(newWord) > 0) {
        newWords = append(newWords, newWord)
      }
    }
  }  

  return newWords
}

// Remove all spaces from string
func clearSpaces(str string) string {
    return strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) {
            return -1
        }
        return r
    }, str)
}
