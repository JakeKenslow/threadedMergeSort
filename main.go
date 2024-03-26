package main

import (
	"fmt"
	"math/rand"

	"github.com/loov/hrtime"
)

type sliceSort func([]int) []int

func mergeSort(s []int) []int {
	// base cases
	switch {
	case len(s) == 1:
		return s
	case len(s) == 2 && s[0] < s[1]:
		return s
	case len(s) == 2 && s[0] >= s[1]:
		return []int{s[1], s[0]}
	}

	// sort two halves
	mid := len(s) / 2
	a := mergeSort(s[:mid])
	b := mergeSort(s[mid:])

	// merge two halves
	i := 0
	j := 0
	sorted := make([]int, len(s))
	for k := 0; k < len(s); k++ {
		switch {
		case i >= len(a):
			sorted[k] = b[j]
			j++
		case j >= len(b):
			sorted[k] = a[i]
			i++
		case b[j] <= a[i]:
			sorted[k] = b[j]
			j++
		default: // when a[i] < b[j]
			sorted[k] = a[i]
			i++
		}
	}
	return sorted
}

func threadedMergeSort(s []int) []int {
	ret := make(chan []int)
	// routines := int(math.Log2(float64(len(s))))
	routines := 2
	go threadedMergeSortRecursive(s, ret, routines)
	return <-ret
}

func threadedMergeSortAlt(s []int) []int {
	ret := make(chan []int)
	// routines := int(math.Log2(float64(len(s))))
	routines := 6
	go threadedMergeSortRecursive(s, ret, routines)
	return <-ret
}

func threadedMergeSortRecursive(s []int, ret chan []int, routines int) {
	// base cases
	switch {
	case len(s) == 1:
		ret <- s
		return
	case len(s) == 2 && s[0] < s[1]:
		ret <- s
		return
	case len(s) == 2 && s[0] >= s[1]:
		ret <- []int{s[1], s[0]}
		return
	}

	// sort two halves
	mid := len(s) / 2
	var a []int
	var b []int
	if routines > 0 {
		aChan := make(chan []int)
		go threadedMergeSortRecursive(s[:mid], aChan, routines-1) // sort first half on another go routine
		b = mergeSort(s[mid:])
		a = <-aChan
	} else {
		a = mergeSort(s[:mid])
		b = mergeSort(s[mid:])
	}

	// merge two halves
	i, j := 0, 0
	sorted := make([]int, len(s))
	for k := 0; k < len(s); k++ {
		switch {
		case i >= len(a):
			sorted[k] = b[j]
			j++
		case j >= len(b):
			sorted[k] = a[i]
			i++
		case b[j] <= a[i]:
			sorted[k] = b[j]
			j++
		default: // when a[i] < b[j]
			sorted[k] = a[i]
			i++
		}
	}
	ret <- sorted
}

func createRandomSlice(size int) []int {
	r := rand.New(rand.NewSource(327014))
	toSort := make([]int, size)
	for i := 0; i < size; i++ {
		toSort[i] = int(r.Int31n(int32(size * 2)))
	}
	return toSort
}

func printLongSlice(s []int) {
	size := len(s)
	if size >= 20 {
		fmt.Println(s[:10], s[size-10:])
	} else {
		fmt.Println(s)
	}
}

func testSortSpeed(sort sliceSort, toSort []int) int64 {
	startTime := hrtime.Now()
	sort(toSort)
	duration := hrtime.Since(startTime)
	return duration.Nanoseconds()
}

func runAllTests(sorts []sliceSort, toSort []int, runCount int) {
	results := make([]int64, len(sorts))
	for i := 0; i < runCount; i++ {
		for j, sort := range sorts {
			results[j] += testSortSpeed(sort, toSort)
		}
	}
	for _, res := range results {
		fmt.Println(res)
	}
}

func main() {
	// input slice size
	size := 1000000

	toSort := createRandomSlice(size)
	printLongSlice(toSort)

	sorts := []sliceSort{threadedMergeSort, threadedMergeSortAlt}

	runAllTests(sorts, toSort, 20)
}
