/*
  fibonacci done different ways as an experiment
*/

package main

import (
	"sync"
)

func NaiveFib(x int) int {
	switch x {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return NaiveFib(x-1) + NaiveFib(x-2)
	}
}

func MemoizedFib(x int) int {
	switch x {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 1
	}
	memo := make([]int, x)
	memo[1] = 1
	memo[2] = 1
	return innerMemoizedFib(x, memo)
}

func innerMemoizedFib(x int, memo []int) int {
	if memo[x-1] == 0 {
		memo[x-1] = innerMemoizedFib(x-1, memo)
	}
	// note that by memoizing x-1 we are also assured x-2 has been calcuated
	return memo[x-1] + memo[x-2]
}

func StacklessFib(x int) int {
	var a = 0
	var b = 1 // b is the fib number following a
	var n = 1 // n is the fib number of b

	for n < x {
		// calcuate the next two fib numbers
		a = a + b
		b = b + a
		n = n + 2
	}

	if n == x {
		return b
	} else {
		return a
	}
}

// taking advantage of the fact that writes to int are atomic, and that the answer is always the same,
// we can let multiple processors calculate and update the same slot in the memo array without danger
// of doing more than a little extra work.
func RaceyConcurrentMemoizedFib(x int) int {
	switch x {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 1
	}
	memo := make([]int, x)
	memo[1] = 1
	memo[2] = 1
	return innerRaceyConcurrentMemoizedFib(x, memo)
}

func innerRaceyConcurrentMemoizedFib(x int, memo []int) int {
	var done sync.Mutex
	var wait bool = false // used to avoid locking the mutex uselessly, which costs more time than checking a local bool variable
	if memo[x-1] == 0 {
		done.Lock()
		wait = true
		go func(y int, memo []int, done *sync.Mutex) {
			memo[y] = innerRaceyConcurrentMemoizedFib(y, memo)
			done.Unlock()
		}(x-1, memo, &done)
	}
	if memo[x-2] == 0 {
		innerRaceyConcurrentMemoizedFib(x-2, memo)
	}
	if wait {
		done.Lock()
	}
	return memo[x-1] + memo[x-2]
}

type LockedMemo struct {
	answer int
	lock   sync.Mutex
}

func LockedConcurrentMemoizedFib(x int) int {
	switch x {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 1
	}
	memo := make([]LockedMemo, x)
	memo[1].answer = 1
	memo[2].answer = 1
	return innerLockedConcurrentMemoizedFib(x, memo)
}

func innerLockedConcurrentMemoizedFib(x int, memo []LockedMemo) int {
	var wait bool = false

	memo[x-1].lock.Lock()
	if memo[x-1].answer == 0 {
		wait = true
		go func(y int, memo []LockedMemo) {
			memo[y].answer = innerLockedConcurrentMemoizedFib(y, memo)
			memo[y].lock.Unlock()
		}(x-1, memo)
	}

	memo[x-2].lock.Lock() // since we all lock in the locks in decending order there is no possibility of deadlock
	if memo[x-2].answer == 0 {
		innerLockedConcurrentMemoizedFib(x-2, memo)
	}
	memo[x-2].lock.Unlock() // we can release x-2 while we [possibly] wait for x-1

	if wait {
		memo[x-1].lock.Lock()
	}
	memo[x-1].lock.Unlock()

	return memo[x-1].answer + memo[x-2].answer
}
