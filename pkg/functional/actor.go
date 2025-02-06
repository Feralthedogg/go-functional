package functional

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
)

type FunctionalActor[T any, U any] struct {
	transform func(T) U
}

func (fa *FunctionalActor[T, U]) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Println("FunctionalActor - Started successfully")

	case T:
		result := fa.transform(msg)
		fmt.Printf("FunctionalActor - Received: %v, Transformed: %v\n", msg, result)

	default:
		fmt.Printf("FunctionalActor - Unhandled message type: %T\n", msg)
	}
}

func NewFunctionalActor[T any, U any](transform func(T) U) *actor.Props {
	return actor.PropsFromProducer(func() actor.Actor {
		return &FunctionalActor[T, U]{transform: transform}
	})
}
