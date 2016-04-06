package main

import (

		"fmt"
		"runtime"

		)

type Data struct {
	Document string
}

func main() {
	var data = []Data{
			  		Data{
			  				"doc1",
			  		},
			  		Data{
			  				"doc2",
			  		},
			  		Data{
			  				"doc3",
			  		},
			  		Data{
			  				"doc4",
			  		},
			  		Data{
			  				"doc5",
			  		},
			  		Data{
			  				"doc6",
			  		},
				}

	Master(data)
}

func Master(work []Data) {
    ncpu := runtime.NumCPU()
    if len(work) < ncpu {
        ncpu = len(work)
    }
    runtime.GOMAXPROCS(ncpu) //max # of threads that can be running

    queue := make(chan *Data)

    // spawn workers
    for i := 0; i < ncpu; i++ {
        go Worker(i, queue)
    }

    // master: give work
    for i, item := range(work) {
        fmt.Printf("master: give work %v\n", item)
        queue <- &work[i] 
    }

    // all work is done
    // signal workers there is no more work
    for n := 0; n < ncpu; n++ {
        queue <- nil
    }
 
    close(queue)

}

func Worker(id int, queue chan *Data) {
    var data *Data
    for {
        data = <-queue
        if data == nil {
            break
        }
        fmt.Printf("worker #%d: item %v\n", id, *data)

        //processData(data)
    }
}