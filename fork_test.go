package fork_test

import (
	"fork"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForkSliceToSlice(t *testing.T) {
	input := []int{0, 1, 2, 3, 4, 5}
	output := []int{0, 1, 4, 9, 16, 25}
	squares := fork.Slice[int, int](input).
		Parallelism(3).
		ToSlice(squareFn)

	assert.Equal(t, output, sorted(squares))
}

func TestForkSliceToChan(t *testing.T) {
	input := []int{0, 1, 2, 3, 4, 5}
	output := []int{0, 1, 4, 9, 16, 25}
	squares := fork.Slice[int, int](input).
		Parallelism(3).
		ToChan(squareFn)
	var result []int
	for {
		square, ok := <-squares
		if !ok {
			break
		}
		result = append(result, square)
	}

	assert.Equal(t, output, sorted(result))
}

func TestForkFromChanToSlice(t *testing.T) {
	inputSize := 5
	output := []int{0, 1, 4, 9, 16}
	input := make(chan int, inputSize+1)
	for i := 0; i < inputSize; i++ {
		input <- i
	}
	close(input)

	squares := fork.Chan[int, int](input).
		Parallelism(2).
		ToSlice(squareFn)

	assert.Equal(t, output, sorted(squares))
}

func TestForkFromChanToChan(t *testing.T) {
	inputSize := 5
	output := []int{0, 1, 4, 9, 16}
	input := make(chan int, inputSize+1)
	for i := 0; i < inputSize; i++ {
		input <- i
	}
	close(input)

	squares := fork.Chan[int, int](input).
		Parallelism(2).
		ToChan(squareFn)

	var result []int
	for {
		square, ok := <-squares
		if !ok {
			break
		}
		result = append(result, square)
	}

	assert.Equal(t, output, sorted(result))
}

func TestForkEarlyExitOnSlice(t *testing.T) {
	input := []int{0, 1, 2, 3, 4, 5}
	squares := fork.Slice[int, int](input).
		Parallelism(3).
		ToSlice(func(_ int) (int, bool) {
			return 0, true
		})

	assert.True(t, len(squares) < len(input))
}

func TestForkEarlyExitOnChan(t *testing.T) {
	inputSize := 5
	input := make(chan int, inputSize+1)
	for i := 0; i < inputSize; i++ {
		input <- i
	}
	close(input)
	squares := fork.Chan[int, int](input).
		Parallelism(3).
		ToChan(func(_ int) (int, bool) {
			return 0, true
		})

	assert.True(t, len(squares) < inputSize)
}

func TestForkMapKeysToSlice(t *testing.T) {
	input := map[int]string{
		0: "zero",
		1: "one",
		2: "two",
	}
	output := []int{0, 1, 4}
	keySquares := fork.Keys[int, string, int](input).
		ToSlice(squareFn)

	assert.Equal(t, output, sorted(keySquares))
}

func TestForkMapValuesToSlice(t *testing.T) {
	input := map[string]int{
		"zero": 0,
		"one":  1,
		"two":  2,
	}
	output := []int{0, 1, 4}
	valueSquares := fork.Values[string, int, int](input).
		ToSlice(squareFn)

	assert.Equal(t, output, sorted(valueSquares))
}

func TestForkDoesntAllowNoParallelism(t *testing.T) {
	input := []int{0, 1, 2, 3, 4, 5}
	output := []int{0, 1, 4, 9, 16, 25}
	squares := fork.Slice[int, int](input).
		Parallelism(0).
		ToSlice(squareFn)

	assert.Equal(t, output, squares)
}

func squareFn(input int) (int, bool) {
	return input * input, false
}

func sorted(values []int) []int {
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	return values
}
