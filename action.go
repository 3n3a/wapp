package wapp

// Type for the function in Action
type ActionFunc = func(*ActionCtx) error

// Main Action Container
type Action struct {
	// function that is executed when you know
	f ActionFunc
}

// NewAction creates and initializes a new Action
//
// f: Expects a function of type ActionFunc
func NewAction(f ActionFunc) Action {
	action := Action{}
	action.f = f
	return action
}