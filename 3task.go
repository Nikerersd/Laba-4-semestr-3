package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numPhilosophers = 5
)

type Fork struct {
	sync.Mutex
}

type Philosopher struct {
	id                int
	leftFork, rightFork *Fork
}

func (p Philosopher) dine(wg *sync.WaitGroup, done chan struct{}) {
	defer wg.Done()

	for {
		select {
		case <-done:
			fmt.Printf("Философ %d закончил обедать.\n", p.id)
			return
		default:
			p.think()
			p.eat()
		}
	}
}

func (p Philosopher) think() {
	fmt.Printf("Филосов %d размышляет о великом.\n", p.id)
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func (p Philosopher) eat() {
	if p.id%2 == 0 {
		p.leftFork.Lock()
		p.rightFork.Lock()
	} else {
		p.rightFork.Lock()
		p.leftFork.Lock()
	}

	fmt.Printf("Философ %d ест спагетти.\n", p.id)
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	p.leftFork.Unlock()
	p.rightFork.Unlock()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	forks := make([]*Fork, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		forks[i] = &Fork{}
	}

	philosophers := make([]*Philosopher, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = &Philosopher{
			id:        i,
			leftFork:  forks[i],
			rightFork: forks[(i+1)%numPhilosophers],
		}
	}

	var wg sync.WaitGroup
	done := make(chan struct{})

	for _, philosopher := range philosophers {
		wg.Add(1)
		go philosopher.dine(&wg, done)
	}

	// Философы едят 5 секунд
	time.Sleep(5 * time.Second)
	close(done)

	wg.Wait()
	fmt.Println("Все философы закончили обедать.")
}