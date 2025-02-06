package functional

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

// RestartStrategy 타입 및 상수 정의
type RestartStrategy int

const (
	OneForOne RestartStrategy = iota
	AllForOne

	OneForOneResume
	AllForOneResume

	OneForOneStop
	AllForOneStop

	OneForOneEscalate
	AllForOneEscalate
)

// NewSupervisedFunctionalActor는 주어진 변환 함수와 재시작 전략에 따라
// Supervisor 옵션이 적용된 FunctionalActor의 Props를 생성합니다.
func NewSupervisedFunctionalActor[T any, U any](transform func(T) U, strategy RestartStrategy) *actor.Props {
	var supervisorStrategy actor.SupervisorStrategy

	switch strategy {
	case OneForOne:
		supervisorStrategy = actor.NewOneForOneStrategy(
			3,             // 최대 재시도 횟수
			5*time.Second, // 기간 내 최대 재시도 시간
			func(reason interface{}) actor.Directive {
				fmt.Printf("OneForOne: Restarting due to: %v\n", reason)
				return actor.RestartDirective
			},
		)
	case AllForOne:
		supervisorStrategy = actor.NewAllForOneStrategy(
			3,
			5*time.Second,
			func(reason interface{}) actor.Directive {
				fmt.Printf("AllForOne: Restarting all children due to: %v\n", reason)
				return actor.RestartDirective
			},
		)
	case OneForOneResume:
		supervisorStrategy = actor.NewOneForOneStrategy(
			3,
			5*time.Second,
			func(reason interface{}) actor.Directive {
				fmt.Printf("OneForOne: Resuming due to: %v\n", reason)
				return actor.ResumeDirective
			},
		)
	case AllForOneResume:
		supervisorStrategy = actor.NewAllForOneStrategy(
			3,
			5*time.Second,
			func(reason interface{}) actor.Directive {
				fmt.Printf("AllForOne: Resuming all children due to: %v\n", reason)
				return actor.ResumeDirective
			},
		)
	case OneForOneStop:
		supervisorStrategy = actor.NewOneForOneStrategy(
			3,
			5*time.Second,
			func(reason interface{}) actor.Directive {
				fmt.Printf("OneForOne: Stopping due to: %v\n", reason)
				return actor.StopDirective
			},
		)
	case AllForOneStop:
		supervisorStrategy = actor.NewAllForOneStrategy(
			3,
			5*time.Second,
			func(reason interface{}) actor.Directive {
				fmt.Printf("AllForOne: Stopping all children due to: %v\n", reason)
				return actor.StopDirective
			},
		)
	case OneForOneEscalate:
		supervisorStrategy = actor.NewOneForOneStrategy(
			3,
			5*time.Second,
			func(reason interface{}) actor.Directive {
				fmt.Printf("OneForOne: Escalating due to: %v\n", reason)
				return actor.EscalateDirective
			},
		)
	case AllForOneEscalate:
		supervisorStrategy = actor.NewAllForOneStrategy(
			3,
			5*time.Second,
			func(reason interface{}) actor.Directive {
				fmt.Printf("AllForOne: Escalating all children due to: %v\n", reason)
				return actor.EscalateDirective
			},
		)
	default:
		supervisorStrategy = actor.DefaultSupervisorStrategy()
	}

	// Props 생성 시, WithSupervisorStrategy 옵션을 함께 전달합니다.
	return actor.PropsFromProducer(func() actor.Actor {
		return &FunctionalActor[T, U]{transform: transform}
	}, actor.WithSupervisor(supervisorStrategy))
}

// ChildSpec은 Supervisor가 관리할 자식 Actor의 사양을 정의합니다.
type ChildSpec struct {
	Name  string
	Props *actor.Props
}

// SupervisorActor는 시작 시 자식 Actor들을 생성하고, 자식이 종료되면 재시작하는 Supervisor입니다.
type SupervisorActor struct {
	specs    []ChildSpec
	children map[*actor.PID]ChildSpec
}

// NewSupervisorActor는 주어진 ChildSpec 목록으로 SupervisorActor를 생성합니다.
func NewSupervisorActor(specs []ChildSpec) actor.Actor {
	return &SupervisorActor{
		specs:    specs,
		children: make(map[*actor.PID]ChildSpec),
	}
}

// Receive 메서드에서 Supervisor는 자식 Actor를 생성하고, 종료된 자식을 재시작합니다.
func (s *SupervisorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Println("Supervisor started, spawning children...")
		for _, spec := range s.specs {
			child := ctx.Spawn(spec.Props)
			s.children[child] = spec
			ctx.Watch(child)
			fmt.Printf("Spawned child '%s': %v\n", spec.Name, child)
		}
	case *actor.Terminated:
		terminatedPID := msg.Who
		spec, exists := s.children[terminatedPID]
		if exists {
			fmt.Printf("Child '%s' terminated: %v. Restarting...\n", spec.Name, terminatedPID)
			newChild := ctx.Spawn(spec.Props)
			delete(s.children, terminatedPID)
			s.children[newChild] = spec
			ctx.Watch(newChild)
			fmt.Printf("Restarted child '%s': %v\n", spec.Name, newChild)
		}
	default:
		// 필요에 따라 추가 메시지 처리
	}
}

// NewSupervisor는 주어진 ChildSpec 목록을 관리하는 SupervisorActor의 PID를 반환합니다.
func NewSupervisor(system *actor.ActorSystem, specs []ChildSpec) *actor.PID {
	props := actor.PropsFromProducer(func() actor.Actor {
		return NewSupervisorActor(specs)
	})
	return system.Root.Spawn(props)
}
