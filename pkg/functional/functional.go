package functional

// ------------------------
// Higher-Order Functions
// ------------------------

func Map[T any, U any](fn func(T) U, list []T) []U {
	result := make([]U, len(list))
	for i, v := range list {
		result[i] = fn(v)
	}
	return result
}

func Reduce[T any, U any](fn func(U, T) U, init U, list []T) U {
	acc := init
	for _, v := range list {
		acc = fn(acc, v)
	}
	return acc
}

// ------------------------
// Currying
// ------------------------

func Curry[A any, B any, C any](fn func(A, B) C) func(A) func(B) C {
	return func(a A) func(B) C {
		return func(b B) C {
			return fn(a, b)
		}
	}
}

// ------------------------
// Function Composition
// ------------------------

func Compose[A, B, C any](f func(B) C, g func(A) B) func(A) C {
	return func(a A) C {
		return f(g(a))
	}
}

// ------------------------
// Maybe Monad
// ------------------------

type Maybe[T any] struct {
	value *T
}

func Some[T any](v T) Maybe[T] {
	return Maybe[T]{value: &v}
}

func None[T any]() Maybe[T] {
	return Maybe[T]{value: nil}
}

func (m Maybe[T]) Map(fn func(T) T) Maybe[T] {
	if m.value == nil {
		return None[T]()
	}
	newVal := fn(*m.value)
	return Some(newVal)
}

func (m Maybe[T]) Bind(fn func(T) Maybe[T]) Maybe[T] {
	if m.value == nil {
		return None[T]()
	}
	return fn(*m.value)
}

func (m Maybe[T]) GetOrElse(defaultValue T) T {
	if m.value == nil {
		return defaultValue
	}
	return *m.value
}

// ------------------------
// Immutable Data Structure
// ------------------------

type ImmutableList[T any] struct {
	items []T
}

func NewImmutableList[T any](items []T) ImmutableList[T] {
	newItems := make([]T, len(items))
	copy(newItems, items)
	return ImmutableList[T]{items: newItems}
}

func (l ImmutableList[T]) Append(item T) ImmutableList[T] {
	newList := make([]T, len(l.items)+1)
	copy(newList, l.items)
	newList[len(newList)-1] = item
	return ImmutableList[T]{items: newList}
}

func (l ImmutableList[T]) GetItems() []T {
	newItems := make([]T, len(l.items))
	copy(newItems, l.items)
	return newItems
}

// ------------------------
// Lazy Evaluation
// ------------------------

func Generate[T any](fn func() T) <-chan T {
	ch := make(chan T)
	go func() {
		for {
			ch <- fn()
		}
	}()
	return ch
}
