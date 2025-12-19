# gosort - Concurrent Chunk Sorting

This is a Go command-line program that sorts integer numbers concurrently using goroutines. It implements a split-sort-merge approach where the input is divided into chunks, each sorted in parallel, and then merged into a single sorted list.

## How to Run

### Prerequisites
- Go 1.25.4 or later

### Building
```bash
go build
```

Or run directly:
```bash
go run .
```

### Modes

#### 1. Random Numbers Mode (-r)
Generates N random integers (range: 0-999) and sorts them.

Usage:
```bash
./gosort -r N
```
Where N >= 10.

Example:
```bash
./gosort -r 20
```

Output:
- Original unsorted numbers
- Chunks before sorting
- Chunks after sorting
- Final merged sorted result

#### 2. Input File Mode (-i)
Reads integers from a plain text file (one integer per line, empty lines ignored) and sorts them.

Usage:
```bash
./gosort -i input.txt
```

The file must contain at least 10 valid integers.

Example:
```bash
./gosort -i numbers.txt
```

Output: Same as -r mode.

#### 3. Directory Mode (-d)
Processes all .txt files in the specified directory, sorting each independently and saving sorted versions to a new directory.

Usage:
```bash
./gosort -d incoming
```

Creates a directory named `incoming_sorted_Avemariya_Perumadan_Siju_241ADB033` containing sorted .txt files with the same filenames.

No console output; sorted files are written to the output directory.

## Design Decisions

- **Chunking**: Number of chunks is max(4, ceil(sqrt(n))) to ensure concurrency even for small n.
- **Concurrency**: Each chunk is sorted in its own goroutine using sync.WaitGroup for synchronization.
- **Merging**: Custom k-way merge to combine sorted chunks without re-sorting the entire list.
- **Error Handling**: Strict validation for inputs; program exits with clear error messages on invalid input.
- **Random Range**: 0-999 for simplicity and readability.
- **Directory Mode**: Processes all .txt files regardless of number of integers (unlike -i which requires >=10).