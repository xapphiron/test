package main

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"strings"
	"strconv"
	"time"
)

const MAX_N = 1000
const ThreadCount = 5

type resultItem struct {
	data []int   // data
	value int    // Result of compute func
}

func main() {
	// Read file
	file, err := os.Open("data.txt")
	if err != nil {
		fmt.Println("Failed to open file.")
		return
	}

	// Create bufio.Reader
	reader := bufio.NewReader(file)

	// Read from file
	data := make([][]int, 0)
	for {
		// Read a line
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line);

		if line == "" {
			data = append(data, make([]int, 0))
			continue
		}

		// Split line
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

	// Timing
	start := time.Now()

	// Base
	base := data[98]

	// Start compute
	var ch = make(chan resultItem, ThreadCount * 5)

	startIndex := 0
	for i := 0; i < ThreadCount; i++ {
		endIndex := startIndex + len(data) / ThreadCount - 1
		if endIndex > len(data) {
			endIndex = len(data)
		}
		//fmt.Println(startIndex, endIndex)
		go computePart(data[:], base[:], startIndex, endIndex, ch)
		startIndex = endIndex + 1
	}
	
	var results = make([]resultItem, 0);
	for i := 0; i < 5 * ThreadCount; i++ {
		item := <-ch

		fmt.Printf("%d ", item.value)
		if (len(results) < 5) {
			results = append(results, item)
		} else {
			minItem, minIndex := findMinResult(results)
			if minItem.value < item.value {
				results[minIndex] = item
			}
		}
	}

	// Output
	fmt.Println("")
	fmt.Println("----------------------------------------")
	for i := 0; i < len(results); i++ {
		fmt.Printf("%d\n",  results[i].value)
	}

	// Timing
	duration := time.Since(start)
	fmt.Println("Time used: ", duration.String())
}

func computePart(data       [][]int, 
				 base       []int, 
				 firstIndex int, 
				 lastIndex  int,
				 resultChan chan <- resultItem) {
	var results = make([]resultItem, 0);

	for i := firstIndex; i < lastIndex; i++ {
		v := compute(base, data[i])

		if (len(results) < 5) {
			results = append(results, resultItem{data[i], v})
		} else {
			minItem, minIndex := findMinResult(results)
			if minItem.value < v {
				results[minIndex] = resultItem{data[i], v}
			}
		}
	}

	for i := 0; i < len(results); i++ {
		resultChan <- results[i]
	}
}

func findMinResult(items []resultItem) (*resultItem, int) {
	var minItem *resultItem = nil;
	var minIndex int;
	for i := 0; i < len(items); i++ {
		if (minItem == nil || items[i].value < minItem.value) {
			minItem = &items[i]
			minIndex = i
		}
	}
	return minItem, minIndex 
}

func compute(s, t []int) int {
    p := make([]int, MAX_N + 1)
    d := make([]int, MAX_N + 1)
    n := len(s)
    m := len(t)

    if n == 0 {
        return m
    } else if (m == 0) {
        return n
    }

    if n > MAX_N {
        n = MAX_N
    }
    if m > MAX_N {
        m = MAX_N
    }


    var t_j int = 0

    var cost int

    for i := 0; i <= n; i++ {
        p[i] = i
    }

    for j := 1; j <= m; j++ {
        t_j = t[j - 1];
        d[0] = j;

        var s_i int = 0;
        for i := 1; i <= n; i++ {
            s_i = s[i - 1];

            if s_i == t_j {
            	cost = 0
            } else {
            	cost = 1
            }

            d[i] = min(d[i - 1] + 1, p[i] + 1, p[i - 1] + cost);
        }

        swap := p;
        p = d;
        d = swap;
    }

    var similarity int = (100 * (max(len(s), len(t)) - p[n])) / max(len(s), len(t));
    return similarity;
}

func min(a, b, c int) int {
	var m int
	if a < b {
		m = a
	} else {
		m = b
	}

	if m < c {
		return m
	} else {
		return c
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
