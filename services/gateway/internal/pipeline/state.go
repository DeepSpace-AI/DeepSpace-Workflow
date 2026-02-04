package pipeline

type State struct {
	RequestBody  []byte
	ResponseBody []byte
	TraceID      string
	OrgID        int64
	Model        string
	CostAmount   float64
	RefID        string
	StatusCode   int
	Error        error
	Halted       bool
	Meta         map[string]any
}

func NewState() *State {
	return &State{Meta: map[string]any{}}
}
