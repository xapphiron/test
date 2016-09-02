package main

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"strings"
	"strconv"
	"time"
	"runtime"
)

const MAX_N = 1000

type resultItem struct {
	data []int   // data
	value int    // Result of compute func
}

func main() {
	// Load args
	if len(os.Args) < 4 {
		fmt.Println("Invalid args")
		fmt.Println("Usage:")
		fmt.Println("  ", os.Args[0], " <data-file> <base-record-index> <num-of-threads>")
		fmt.Println("")
		return
	}

	dataFile := os.Args[1]
	baseIndex, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid base index")
		return
	}
	threads, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid number of threads")
		return
	}

	// Set max threads
	runtime.GOMAXPROCS(threads)

	// Read file
	file, err := os.Open(dataFile)
	if err != nil {
		fmt.Println("Failed to open file.")
		return
	}

	// Create bufio.Reader
	reader := bufio.NewReader(file)

	// Read from file
	data := make([][]int, 0, 100000)
	for {
		// Read a line
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line);

		if line == "" {
			data = append(data, make([]int, 0, 1))
			continue
		}

		// Split each line
		parts := strings.Split(line, "\t")
		var intParts []int = make([]int, len(parts))
		for i := 0; i < len(parts); i++ {
			intParts[i], err = strconv.Atoi(parts[i])
			if err != nil {
				fmt.Printf("Invalid line: %s\n", line)
			}
		}
		data = append(data, intParts)
	}

	// Timing start
	start := time.Now()

	// Base record
	base := data[baseIndex]

	// Channel to receive results
	var ch = make(chan *resultItem, threads * 5)

	// Start routines
	startIndex := 0
	for i := 0; i < threads; i++ {
		endIndex := startIndex + len(data) / threads
		if endIndex > len(data) {
			endIndex = len(data)
		}
		//fmt.Println(startIndex, endIndex)
		go computePart(data[:], base[:], startIndex, endIndex, ch)
		startIndex = endIndex
	}
	
	// Get results
	var results = make([]*resultItem, 0, 5);
	for i := 0; i < 5 * threads; i++ {
		item := <-ch

		if (len(results) < 5) {
			results = append(results, item)
		} else {
			minItem, minIndex := findMinResult(results)
			if minItem.value < item.value {
				results[minIndex] = item
			}
		}
	}

	// Timing
	duration := time.Since(start)
	fmt.Println("Time used: ", duration.String())

	// Output
	for i := 0; i < len(results); i++ {
		fmt.Printf("%d\n", results[i].value)
	}
}

func computePart(data       [][]int, 
				 base       []int, 
				 firstIndex int, 
				 lastIndex  int,
				 resultChan chan <- *resultItem) {
	runtime.LockOSThread()

	levenshtein := NewLevenshteinDistance()

	var results = make([]*resultItem, 0, 5);
	for i := firstIndex; i < lastIndex; i++ {
		v := levenshtein.compute(base, data[i])

		if (len(results) < 5) {
			results = append(results, &resultItem{data[i], v})
		} else {
			minItem, minIndex := findMinResult(results)
			if minItem.value < v {
				results[minIndex].data = data[i]
				results[minIndex].value = v
			}
		}
	}

	for i := 0; i < len(results); i++ {
		resultChan <- results[i]
	}
}

func findMinResult(items []*resultItem) (*resultItem, int) {
	var minItem *resultItem = nil;
	var minIndex int;
	for i := 0; i < len(items); i++ {
		if (minItem == nil || items[i].value < minItem.value) {
			minItem = items[i]
			minIndex = i
		}
	}
	return minItem, minIndex 
}



//=================================================================
// LevenshteinDistance
//=================================================================

type LevenshteinDistance struct {
	arr1, arr2 []int16
}

func NewLevenshteinDistance() *LevenshteinDistance {
	obj := new(LevenshteinDistance)
    obj.arr1 = make([]int16, MAX_N + 1)
    obj.arr2 = make([]int16, MAX_N + 1)
    return obj
}

func (this LevenshteinDistance) compute(s, t []int) int {
	p := this.arr1[:]
	d := this.arr2[:]
    n := len(s)
    m := len(t)

    if n == 0 {
        return m
    } else if (m == 0) {
        return n
    }

    var maxlen int = n
    if m > maxlen {
    	maxlen = m
    }

    if n > MAX_N {
        n = MAX_N
    }
    if m > MAX_N {
        m = MAX_N
    }

    var t_j int
    var s_i int
    var cost int16

    for i := 0; i <= n; i++ {
        p[i] = int16(i)
    }

    for j := 1; j <= m; j++ {
        t_j = t[j - 1]
        d[0] = int16(j)

        for i := 1; i <= n; i++ {
            s_i = s[i - 1];

            if s_i == t_j {
            	cost = 0
            } else {
            	cost = 1
            }

            m := d[i - 1] + 1;
            if p[i] + 1 < m {
            	m = p[i] + 1
            }
            if p[i - 1] + cost < m {
            	m = p[i - 1] + cost
            }
            d[i] = m
        }

        p, d = d, p
    }

    var similarity int = (100 * (maxlen - int(p[n]))) / maxlen;
    return similarity;
}
