# Locally Parallel Index Construction

Leveraging Go's concurrency features for efficient inverted index construction and interactive querying.

## usage

Simply make to build:
```bash
make
```
### Flags
```bash
Usage of lpic:
  -json
    	generate additional JSON formatted index file
  -norm-tf
    	log normalize raw term frequency
  -num-results int
    	number of query results to show (default 5)
  -num-workers int
    	number of worker threads (default 4)
  -out-dir string
    	destination directory for constructed index (default "./")
  -out-file string
    	file name of constructed index (default "index")
  -verbose
    	print verbose progress
```

### Index Construction
```bash
./lpic [optional flags] build [required target root directory to crawl for documents]
```

### Index Querying
```bash
./lpic [optional flags] query [generated index file in .lpic format]
```

## license 
MIT (see [LICENSE](https://github.com/mleef/LPIC/blob/master/LICENSE) file)