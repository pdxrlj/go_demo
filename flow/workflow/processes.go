package workflow

type Processes struct {
	UniqueId string
	nodes    map[string]*Node
}

func (p *Processes) Node(name string, workload workloadHandler) *Node {
	node := p.GetNode(name)
	if node != nil {

	}
	return nil
}

func (p *Processes) GetNode(name string) *Node {
	return p.nodes[name]
}

func (p *Processes) Add() {

}
