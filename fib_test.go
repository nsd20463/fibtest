/*
  fibonacci done different ways as an experiment
*/

package main

import (
	"flag"
	"testing"
)

func TestCorrectness(t *testing.T) {
	var a [5]int
	var x int
	for x = 0; x < 25; x++ {
		a[0] = NaiveFib(x)
		a[1] = MemoizedFib(x)
		a[2] = StacklessFib(x)
		a[3] = RaceyConcurrentMemoizedFib(x)
		a[4] = LockedConcurrentMemoizedFib(x)

		t.Logf("NaiveFib(%d) = %d\n", x, a[0])
		for i := range a {
			if a[0] != a[i] {
				t.Errorf("fib(%d) answers don't match: %d\n", x, a)
				break
			}
		}
	}
}

var n int = 25

func init() {
	flag.IntVar(&n, "n", n, "fib number to benchmark")
}

func BenchmarkNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NaiveFib(n)
	}
}

func BenchmarkMemoized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MemoizedFib(n)
	}
}

func BenchmarkStackless(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StacklessFib(n)
	}
}

func BenchmarkRaceyConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RaceyConcurrentMemoizedFib(n)
	}
}

func BenchmarkLockedConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LockedConcurrentMemoizedFib(n)
	}
}
