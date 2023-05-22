package workflow

type Workflow struct {
}

func (w *Workflow) Processes() *Processes {
	return &Processes{}
}
