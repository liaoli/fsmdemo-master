package fsm

// EventProcessor defines OnExit, Action and OnEnter actions.
type EventProcessor interface {
	// OnExit Action handles exiting a state
	OnExit(fromState State, args []interface{})
	// PreAction is used to handle transitions
	PreAction(action Action, fromState State, toState State, args []interface{}) error

	// NextAction is used to handle transitions
	NextAction(action Action, fromState State, toState State, args []interface{}) error

	// OnActionFailure failed to execute the Action
	OnActionFailure(action Action, fromState State, toState State, args []interface{}, err error)

	// OnEnter Action handles entering a state
	OnEnter(toState State, args []interface{})
}

// DefaultDelegate is a default delegate.
// it splits processing of actions into three actions: OnExit, Action and OnEnter.
type DefaultDelegate struct {
	P EventProcessor
}

// HandleEvent implements Delegate interface and split HandleEvent into three actions.
func (dd *DefaultDelegate) HandleEvent(preAction Action, nextAction Action, fromState State, toState State, args []interface{}) error {
	if fromState != toState {
		dd.P.OnExit(fromState, args)
	}

	var err error
	if preAction.Id > 0 {
		err = dd.P.PreAction(preAction, fromState, toState, args)

		if err != nil {
			dd.P.OnActionFailure(preAction, fromState, toState, args, err)
			return err
		}
	}

	if err == nil && fromState != toState  {
		dd.P.OnEnter(toState, args)
		if nextAction.Id > 0 {
			err = dd.P.NextAction(nextAction, fromState, toState, args)
			if err != nil {
				dd.P.OnActionFailure(nextAction, fromState, toState, args, err)
				return err
			}
		}

	}

	return err
}
