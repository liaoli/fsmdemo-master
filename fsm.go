// gofsm is a simple, featured FSM implementation that has some different features with other FSM implementation.
// One feature of gofsm is it doesn't persist/keep states of objects. When it processes transitions, you must pass current states to Id, so you can look gofsm as a "stateless" state machine. This benefit is one gofsm instance can be used to handle transitions of a lot of object instances, instead of creating a lot of FSM instances. Object instances maintain their states themselves.
// Another feature is it provides a common interface for Moore and Mealy FSM. You can implement corresponding methods (OnExit, PreAction, OnEnter) for those two FSM.
// The third interesting feature is you can export configured transitions into a state diagram. A picture is worth a thousand words.

// Style of gofsm refers to implementation of https://github.com/elimisteve/fsm.

package fsm

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Transition is a state transition and all data are literal values that simplifies FSM usage and make it generic.
type Transition struct {
	From State
	Event
	To     State
	PreAction Action
	NextAction Action
}

type State struct {
	Id          int64
	Name        string
	Description string
}

type Event struct {
	Id   int64
	Name string
}

type Action struct {
	Id int64
	Name string
}

// Delegate is used to process actions. Because gofsm uses literal values as event, state and action, you need to handle them with corresponding functions. DefaultDelegate is the default delegate implementation that splits the processing into three actions: OnExit Action, Action and OnEnter Action. you can implement different delegates.
type Delegate interface {
	// HandleEvent handles transitions
	HandleEvent(preAction Action, NextAction Action,fromState State, toState State, args []interface{}) error
}

// StateMachine is a FSM that can handle transitions of a lot of objects. delegate and transitions are configured before use them.
type StateMachine struct {
	delegate    Delegate
	transitions []Transition
}

// Error is an error when processing event and state changing.
type Error interface {
	error
	BadEvent() Event
	CurrentState() State
}

type smError struct {
	badEvent     Event
	currentState State
}

func (e smError) Error() string {
	return fmt.Sprintf("state machine error: cannot find transition for event [%v] when in state [%v]\n", e.badEvent, e.currentState)
}

	func (e smError) BadEvent() Event {
	return e.badEvent
}

func (e smError) CurrentState() State {
	return e.currentState
}

// NewStateMachine creates a new state machine.
func NewStateMachine(delegate Delegate, transitions ...Transition) *StateMachine {
	return &StateMachine{delegate: delegate, transitions: transitions}
}

// Trigger fires a event. You must pass current state of the processing object, other info about this object can be passed with args.
func (m *StateMachine) Trigger(currentState State, event Event, args ...interface{}) error {
	trans := m.findTransMatching(currentState, event)
	if trans == nil {
		return smError{event, currentState}
	}

	var err error
	//if trans.PreAction != "" {
	err = m.delegate.HandleEvent(trans.PreAction,trans.NextAction, currentState, trans.To, args)
	//}
	return err
}

// findTransMatching gets corresponding transition according to current state and event.
func (m *StateMachine) findTransMatching(fromState State, event Event) *Transition {
	for _, v := range m.transitions {
		//if v.From == fromState && v.Event == event {
		//	return &v
		//}

		if v.From.Name == fromState.Name && v.From.Id == fromState.Id && v.Event.Id == event.Id {
			return &v
		}
	}
	return nil
}

// Export exports the state diagram into a file.
func (m *StateMachine) Export(outfile string) error {
	return m.ExportWithDetails(outfile, "png", "dot", "72", "-Gsize=10,5 -Gdpi=200")
}

// ExportWithDetails  exports the state diagram with more graphviz options.
func (m *StateMachine) ExportWithDetails(outfile string, format string, layout string, scale string, more string) error {
	dot := `digraph StateMachine {

	rankdir=LR
	node[width=1 fixedsize=true shape=circle style=filled fillcolor="darkorchid1" ]
	
	`

	for _, t := range m.transitions {
		link := fmt.Sprintf(`%s -> %s [label="%s | %s"]`, t.From.Name, t.To.Name, t.Event.Name, t.PreAction.Name)
		dot = dot + "\r\n" + link
	}

	dot = dot + "\r\n}"
	cmd := fmt.Sprintf("dot -o%s -T%s -K%s -s%s %s", outfile, format, layout, scale, more)

	return system(cmd, dot)
}

func system(c string, dot string) error {

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(`cmd`, `/C`, c)
	} else {
		cmd = exec.Command(`/bin/sh`, `-c`, c)
	}
	cmd.Stdin = strings.NewReader(dot)
	return cmd.Run()

}
