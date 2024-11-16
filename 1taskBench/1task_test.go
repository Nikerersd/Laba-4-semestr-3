package main

import (
	"sync"
	"testing"
	"time"
)

// Бенчмарк для Mutex
func BenchmarkMutex(b *testing.B) {
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	for i := 0; i < b.N; i++ {
		wg.Add(numGoroutines)
		for j := 0; j < numGoroutines; j++ {
			go testMutex(&wg, mu)
		}
		wg.Wait()
	}
}

// Бенчмарк для Semaphore
func BenchmarkSemaphore(b *testing.B) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 3) // Ограничение на 3 горутины
	for i := 0; i < b.N; i++ {
		wg.Add(numGoroutines)
		for j := 0; j < numGoroutines; j++ {
			go testSemaphore(&wg, sem)
		}
		wg.Wait()
	}
}

// Бенчмарк для SemaphoreSlim
func BenchmarkSemaphoreSlim(b *testing.B) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 3)
	retries := 5
	for i := 0; i < b.N; i++ {
		wg.Add(numGoroutines)
		for j := 0; j < numGoroutines; j++ {
			go testSemaphoreSlim(&wg, sem, retries)
		}
		wg.Wait()
	}
}

// Бенчмарк для Barrier
func BenchmarkBarrier(b *testing.B) {
	var wg sync.WaitGroup
	barrier := &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		barrier.Add(numGoroutines)
		wg.Add(numGoroutines)
		for j := 0; j < numGoroutines; j++ {
			go testBarrier(&wg, barrier)
		}
		wg.Wait()
	}
}

// Бенчмарк для SpinLock
func BenchmarkSpinLock(b *testing.B) {
	var wg sync.WaitGroup
	var counter int32
	for i := 0; i < b.N; i++ {
		wg.Add(numGoroutines)
		for j := 0; j < numGoroutines; j++ {
			go func() {
				defer wg.Done()
				testSpinLock(&counter)
			}()
		}
		wg.Wait()
	}
}

// Бенчмарк для SpinWait
func BenchmarkSpinWait(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(numGoroutines)
		for j := 0; j < numGoroutines; j++ {
			go func() {
				defer wg.Done()
				testSpinWait()
			}()
		}
		wg.Wait()
	}
}

// Бенчмарк для Monitor
func BenchmarkMonitor(b *testing.B) {
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	cond := sync.NewCond(mu)
	for i := 0; i < b.N; i++ {
		wg.Add(numGoroutines)
		for j := 0; j < numGoroutines; j++ {
			go testMonitor(&wg, mu, cond)
		}
		time.Sleep(time.Microsecond * 1000) // даём время горутинам заблокироваться
		cond.Broadcast()
		wg.Wait()
	}
}
