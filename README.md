# go-functional
**A Powerful Functional Programming and Actor Model Library For Go**
`go-functional` is a Go library that combines **Functional Programming (FP) and the Actor Model** to enable efficient concurrent programming.
it leverages **Proto.Actor** to provide **Supervisor Trees and Restart Strategies** while supporting **functional utilities (Map, Reduce, Curry, Compose, etc.). Monads, Immutable Data Structures, and Lazy Evaluation.**

---

## Features
- `Map, Reduce` -> Transform and reduce slice data efficiently
- `Curry & Compse` -> High-order functions, currying, and function composition
- `Maybe Monad` -> Safe value handling (`Some`, `None`)
- `Immutale List` -> Immutable list implementation
- `Lazy Evaluation` -> Generate values only when needed

## Actor Model-Based Concurrency
- **Leverages Proto.Actor for actor-based concurrency**
- **FunctionalActor** -> Applies functional programming inside the Actor system
- **Supervisor Trees & Restart Strategies** -> Supports `OneForOne`, `AllForOne`, `Resume`, `Stop`, `Escalate`
- **High-Performance Message Processing** -> Handles thousands of concurrent actors efficiently

## Benchmark & Performance Optimization
- **Benchmarking Map/Reduce operations**
- **Evaluating Actor System and Supervisor performance**
- **Comparing Lazy Evaluation execution speeds**

---

## Directory Structure
```
go-functional/
├── go.mod
├── cmd/
|   └── main.go           # Benchmark and execution example
|
└── pkg/
    └── functional/
        ├── actor.go      # Actor-relaed logic
        ├── functional.go # Functional programming utilities
        └── supervisor.go # Supervisor trees & restart strategies
```

## Installation
**1. Install the package**
```bash
go get github.com/Feralthedogg/go-functional
```
```bash
go get github.com/asynkron/protoactor-go
```
**2. Run benchmarks**
```bash
go run cmd/main.go
```

## Usage Examples
### **Functional Programming (Map, Reduce, Curry)**
```go
import "github.com/Feralthedogg/go-functional/pkg/functional"

numbers := []int{1, 2, 3, 4, 5}

// Map: Square each element
squares := functional.Map(func(x int) int { return x * x }, numbers)

// Reduce: Sum all elements
sum := functional.Reduce(func(acc, x int) int { return acc + x }, 0, numbers)

// Curry: Partial application
add := functional.Curry(func(a, b int) int { return a + b })
addFive := add(5)
result := addFive(10) // 5+ 10 = 15
```

---

### **Lazy Evaluation**
```go
count := 0
gen := functional.Generate(func() int {
    count++
    return count
})

// Retrieve fist 5 values
for i := 0; i < 5; i++ {
    fmt.Println(<-gen)
}
```

---

### **Actor Model + Supervisor Example**
```go
import (
	"go-functional/pkg/functional"
	"github.com/asynkron/protoactor-go/actor"
)

system := actor.NewActorSystem()

// Transformation function: (x+1) * 2
transform := func(x int) int { return (x + 1) * 2 }

// Create an Actor with OneForOne restart strategy
props := functional.NewSupervisedFunctionalActor(transform, functional.OnforOne)

// Create a Supervisor managing two child actors
supervisor := functional.NewSupervisor(system, []functional.ChildSpec{
    {Name: "Worker1", Props: props},
    {Name: "Worker2", Props: props},
})

// Send message to actor
system.Root.Send(supervisor, 10)
```

---

## Benchmark Results (Example)
```bash
=== Benchmarking All Features ===
=== Benchmark: Map Function ===     
Map completed: len(result) = 1000000
Time taken: 5.769ms

=== Benchmark: Reduce Function ===  
Reduce result: 499999500000
Time taken: 2.7015ms

=== Benchmark: Curry & Compose ===  
Curry result (5 + 10): 15
Compose result ( (3+1)*2 ): 8       
Time taken: 0s

=== Benchmark: Maybe Monad ===      
Maybe Monad result: 300
Time taken: 506.7µs

=== Benchmark: Immutable List ===   
Immutable List result: [1 2 3 4 5 6]
Time taken: 0s

=== Benchmark: Lazy Evaluation ===  
Lazy value 1: 1
Lazy value 2: 2
Lazy value 3: 3
Lazy value 4: 4
Lazy value 5: 5
Lazy value 6: 6
Lazy value 7: 7
Lazy value 8: 8
Lazy value 9: 9
Lazy value 10: 10
Time taken: 1.0457ms

=== Benchmark: Actor System with Supervisor ===
9:17PM INF actor system started lib=Proto.Actor system=aM6ZbRt44HpKuLYa4QMTAp id=aM6ZbRt44HpKuLYa4QMTAp
Supervisor started, spawning children...
FunctionalActor - Started successfully
Spawned child 'Transformer': Address:"nonhost"  Id:"$1/$2"
Spawned child 'Transformer': Address:"nonhost"  Id:"$1/$3"
FunctionalActor - Started successfully
Actor system benchmark completed
Time taken: 1.0638256s
=== All Benchmarks Completed ===
```
**Fast Map/Reduce operations, optimized currying & functional composition, and stable Actor system exection using Supervisor.**

---

## License
MIT License
For details, refer to the LICENSE file.