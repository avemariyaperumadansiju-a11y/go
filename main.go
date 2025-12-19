// Name: Avemariya Perumadan Siju
// Student ID: 241ADB033
// Course: Go Programming Assignment – Concurrent Chunk Sorting (gosort)

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// --------------------------------------------------
// Entry point
// --------------------------------------------------

func main() {
	rFlag := flag.Int("r", -1, "generate N random integers (N >= 10)")
	iFlag := flag.String("i", "", "input file with integers")
	dFlag := flag.String("d", "", "directory containing .txt files")

	flag.Parse()

	if *rFlag != -1 {
		if err := runRandom(*rFlag); err != nil {
			log.Fatal(err)
		}
		return
	}

	if *iFlag != "" {
		if err := runInputFile(*iFlag); err != nil {
			log.Fatal(err)
		}
		return
	}

	if *dFlag != "" {
		if err := runDirectory(*dFlag); err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Fatal("Usage: gosort -r N | -i input.txt | -d directory")
}

// --------------------------------------------------
// Mode -r : Random numbers
// --------------------------------------------------

func runRandom(n int) error {
	if n < 10 {
		return errors.New("N must be >= 10")
	}

	numbers := generateRandomNumbers(n)

	fmt.Println("Original numbers:")
	fmt.Println(numbers)

	processAndPrint(numbers)
	return nil
}

// --------------------------------------------------
// Mode -i : Input file
// --------------------------------------------------

func runInputFile(filename string) error {
	numbers, err := readNumbersFromFile(filename)
	if err != nil {
		return err
	}

	if len(numbers) < 10 {
		return errors.New("input file must contain at least 10 valid integers")
	}

	fmt.Println("Original numbers:")
	fmt.Println(numbers)

	processAndPrint(numbers)
	return nil
}

// --------------------------------------------------
// Mode -d : Directory
// --------------------------------------------------

func runDirectory(dir string) error {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return errors.New("invalid directory")
	}

	outDir := dir + "_sorted_YOUR_FIRSTNAME_YOUR_SURNAME_YOUR_STUDENT_ID"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".txt" {
			continue
		}

		inPath := filepath.Join(dir, e.Name())
		numbers, err := readNumbersFromFile(inPath)
		if err != nil || len(numbers) < 10 {
			continue
		}

		sorted := process(numbers)

		outPath := filepath.Join(outDir, e.Name())
		if err := writeNumbersToFile(outPath, sorted); err != nil {
			return err
		}
	}

	return nil
}

// --------------------------------------------------
// Core processing
// --------------------------------------------------

func processAndPrint(numbers []int) {
	chunks := splitIntoChunks(numbers)

	fmt.Println("\nChunks before sorting:")
	printChunks(chunks)

	sortedChunks := sortChunksConcurrently(chunks)

	fmt.Println("\nChunks after sorting:")
	printChunks(sortedChunks)

	result := mergeSortedChunks(sortedChunks)

	fmt.Println("\nFinal sorted result:")
	fmt.Println(result)
}

func process(numbers []int) []int {
	chunks := splitIntoChunks(numbers)
	sortedChunks := sortChunksConcurrently(chunks)
	return mergeSortedChunks(sortedChunks)
}

// --------------------------------------------------
// Chunking logic
// --------------------------------------------------

func splitIntoChunks(numbers []int) [][]int {
	n := len(numbers)

	numChunks := int(math.Ceil(math.Sqrt(float64(n))))
	if numChunks < 4 {
		numChunks = 4
	}

	chunks := make([][]int, 0, numChunks)

	baseSize := n / numChunks
	extra := n % numChunks

	start := 0
	for i := 0; i < numChunks; i++ {
		size := baseSize
		if i < extra {
			size++
		}
		end := start + size
		if start < n {
			chunks = append(chunks, numbers[start:end])
		}
		start = end
	}

	return chunks
}

// --------------------------------------------------
// Concurrent sorting
// --------------------------------------------------

func sortChunksConcurrently(chunks [][]int) [][]int {
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for i := range chunks {
		go func(idx int) {
			defer wg.Done()
			sort.Ints(chunks[idx])
		}(i)
	}

	wg.Wait()
	return chunks
}

// --------------------------------------------------
// Merge logic (k-way merge)
// --------------------------------------------------

func mergeSortedChunks(chunks [][]int) []int {
	result := make([]int, 0)

	indices := make([]int, len(chunks))

	for {
		minVal := 0
		minChunk := -1

		for i, chunk := range chunks {
			if indices[i] < len(chunk) {
				val := chunk[indices[i]]
				if minChunk == -1 || val < minVal {
					minVal = val
					minChunk = i
				}
			}
		}

		if minChunk == -1 {
			break
		}

		result = append(result, minVal)
		indices[minChunk]++
	}

	return result
}

// --------------------------------------------------
// Helpers
// --------------------------------------------------

func generateRandomNumbers(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rand.Intn(1000) // range: 0–999
	}
	return nums
}

func readNumbersFromFile(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var numbers []int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		val, err := strconv.Atoi(line)
		if err != nil {
			return nil, errors.New("invalid integer in file")
		}
		numbers = append(numbers, val)
	}

	return numbers, nil
}

func writeNumbersToFile(filename string, numbers []int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, n := range numbers {
		fmt.Fprintln(writer, n)
	}
	return writer.Flush()
}

func printChunks(chunks [][]int) {
	for i, c := range chunks {
		fmt.Printf("Chunk %d: %v\n", i, c)
	}
}
