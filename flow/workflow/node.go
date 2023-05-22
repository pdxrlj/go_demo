package workflow

type Node struct {
	UniqueId        string
	subProcesses    *Processes
	conditionalNode map[string]Processes

	next []*Node
	pre  []*Node

	nodeIndex int
}
