package main

import (
	"fmt"
	"time"

	"github.com/Feralthedogg/go-functional/pkg/functional"

	"github.com/asynkron/protoactor-go/actor"
)

func benchmarkMap() {
	fmt.Println("=== Benchmark: Map Function ===")
	start := time.Now()
	input := make([]int, 1000000)
	for i := 0; i < len(input); i++ {
		input[i] = i
	}
	result := functional.Map(func(x int) int {
		return x * x
	}, input)
	fmt.Printf("Map completed: len(result) = %d\n", len(result))
	fmt.Println("Time taken:", time.Since(start))
}

func benchmarkReduce() {
	fmt.Println("=== Benchmark: Reduce Function ===")
	start := time.Now()
	input := make([]int, 1000000)
	for i := 0; i < len(input); i++ {
		input[i] = i
	}
	sum := functional.Reduce(func(acc, x int) int {
		return acc + x
	}, 0, input)
	fmt.Printf("Reduce result: %d\n", sum)
	fmt.Println("Time taken:", time.Since(start))
}

func benchmarkCurryAndCompose() {
	fmt.Println("=== Benchmark: Curry & Compose ===")
	start := time.Now()
	add := functional.Curry(func(a, b int) int {
		return a + b
	})
	addFive := add(5)
	result := addFive(10)
	fmt.Printf("Curry result (5 + 10): %d\n", result)

	increment := func(x int) int { return x + 1 }
	multiply := func(x int) int { return x * 2 }
	composed := functional.Compose(multiply, increment)
	result = composed(3)
	fmt.Printf("Compose result ( (3+1)*2 ): %d\n", result)
	fmt.Println("Time taken:", time.Since(start))
}

func benchmarkMaybeMonad() {
	fmt.Println("=== Benchmark: Maybe Monad ===")
	start := time.Now()
	maybeVal := functional.Some(100)
	result := maybeVal.Map(func(x int) int { return x * 3 }).GetOrElse(0)
	fmt.Printf("Maybe Monad result: %d\n", result)
	fmt.Println("Time taken:", time.Since(start))
}

func benchmarkImmutableList() {
	fmt.Println("=== Benchmark: Immutable List ===")
	start := time.Now()
	list := functional.NewImmutableList([]int{1, 2, 3, 4, 5})
	newList := list.Append(6)
	fmt.Printf("Immutable List result: %v\n", newList.GetItems())
	fmt.Println("Time taken:", time.Since(start))
}

func benchmarkLazyEvaluation() {
	fmt.Println("=== Benchmark: Lazy Evaluation ===")
	start := time.Now()
	count := 0
	lazyGen := functional.Generate(func() int {
		count++
		return count
	})
	for i := 0; i < 10; i++ {
		fmt.Printf("Lazy value %d: %d\n", i+1, <-lazyGen)
	}
	fmt.Println("Time taken:", time.Since(start))
}

func benchmarkActorSystem() {
	fmt.Println("=== Benchmark: Actor System with Supervisor ===")
	system := actor.NewActorSystem()

	transform := func(x int) int { return (x + 1) * 2 }
	supervisedProps := functional.NewSupervisedFunctionalActor(transform, functional.OneForOneResume)
	childSpec := functional.ChildSpec{
		Name:  "Transformer",
		Props: supervisedProps,
	}
	supervisorPID := functional.NewSupervisor(system, []functional.ChildSpec{childSpec, childSpec})

	start := time.Now()
	for i := 0; i < 100; i++ {
		system.Root.Send(supervisorPID, 10)
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("Actor system benchmark completed")
	fmt.Println("Time taken:", time.Since(start))

	time.Sleep(2 * time.Second)
}

func main() {
	fmt.Println("=== Benchmarking All Features ===")
	benchmarkMap()
	fmt.Println()
	benchmarkReduce()
	fmt.Println()
	benchmarkCurryAndCompose()
	fmt.Println()
	benchmarkMaybeMonad()
	fmt.Println()
	benchmarkImmutableList()
	fmt.Println()
	benchmarkLazyEvaluation()
	fmt.Println()
	benchmarkActorSystem()
	fmt.Println("=== All Benchmarks Completed ===")
}
