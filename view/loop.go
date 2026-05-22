package view

// Loop provides all standard Blade $loop variables inside @foreach.
type Loop struct {
	Index      int  // 0-based
	Iteration  int  // 1-based
	Remaining  int
	First      bool
	Last       bool
	Even       bool
	Odd        bool
	Depth      int
	Parent     *Loop // for nested loops (future)
}

// newLoop creates a Loop instance for the current iteration.
func newLoop(index, total int) *Loop {
	iteration := index + 1
	return &Loop{
		Index:     index,
		Iteration: iteration,
		Remaining: total - iteration,
		First:     index == 0,
		Last:      index == total-1,
		Even:      iteration%2 == 0,
		Odd:       iteration%2 == 1,
		Depth:     1,
	}
}
