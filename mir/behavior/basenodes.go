package behavior

import "time"

type Node struct {
	children []INode
	status   Status
}

func (n *Node) Status() Status {
	return n.status
}

func (n *Node) Reset() {
	if n.status != READY {
		n.status = READY
		if n.children != nil {
			for _, child := range n.children {
				child.Reset()
			}
		}
	}
}

func (n *Node) Step() {
	if n.status != RUNNING {
		n.Reset()
	} else if n.children != nil {
		for _, child := range n.children {
			child.Step()
		}
	}
}

type PriorityNode struct {
	Node
	idx      int
	period   time.Duration
	lasttime time.Duration
}

func (n *PriorityNode) Reset() {
	n.idx = -1
}

func (n *PriorityNode) Visit(c *BT) {

	now := c.GetTime()

	eval := n.lasttime == 0 || n.period == 0 || n.lasttime+n.period < now

	if eval {
		n.lasttime = now

		found := false
		for idx, child := range n.children {

			if !found {

				if child.Status() == FAILED || child.Status() == SUCCESS {
					child.Reset()
				}

				child.Visit(c)

				if child.Status() == SUCCESS || child.Status() == RUNNING {
					n.status = child.Status()
					found = true
					n.idx = idx
				}
			} else {
				child.Reset()
			}
		}
		if !found {
			n.status = FAILED
		}
	} else {
		if n.idx != -1 {
			child := n.children[n.idx]
			if child.Status() != RUNNING {
				child.Visit(c)
				n.status = child.Status()
				if n.status != RUNNING {
					n.lasttime = 0
				}
			}
		}
	}
}

type SequenceNode struct {
	Node
	idx int
}

func (n *SequenceNode) Reset() {
	n.Node.Reset()
	n.idx = 0
}

func (n *SequenceNode) Visit(c *BT) {
	if n.Status() != RUNNING {
		n.idx = 0
	}

	for i := n.idx; i < len(n.children); i++ {
		child := n.children[i]
		child.Visit(c)
		if child.Status() == RUNNING || child.Status() == FAILED {
			n.status = child.Status()
			return
		}
	}

	n.status = SUCCESS
}

type ParallelNode struct {
	Node
	stoponanycomplete bool
}

func (n *ParallelNode) Visit(c *BT) {
	done, anydone := true, false
	for _, child := range n.children {
		switch child.(type) {
		case *ConditionNode:
			child.Reset()
		}
		if child.Status() != SUCCESS {
			child.Visit(c)
			if child.Status() == FAILED {
				n.status = FAILED
				return
			}
		}

		if child.Status() == RUNNING {
			done = false
		} else {
			anydone = true
		}

	}

	if done || (n.stoponanycomplete && anydone) {
		n.status = SUCCESS
	} else {
		n.status = RUNNING
	}
}

func (n *ParallelNode) Step() {
	if n.status != RUNNING {
		n.Reset()
	} else if n.children != nil {
		for _, child := range n.children {
			if child.Status() == SUCCESS {
				switch child.(type) {
				case *ConditionNode:
					child.Reset()
				}
			}
		}
	}
}

type ConditionFunc func(*BT) bool
type ConditionNode struct {
	Node
	fn ConditionFunc
}

func (n *ConditionNode) Visit(c *BT) {
	if n.fn(c) {
		n.status = SUCCESS
	} else {
		n.status = FAILED
	}
}

func Priority(period time.Duration, children ...INode) INode {
	return &PriorityNode{
		Node: Node{
			children: children,
		},
		period: period,
		idx:    -1,
	}
}

func Parallel(children ...INode) INode {
	return &ParallelNode{
		Node: Node{
			children: children,
		},
	}
}

func Sequence(children ...INode) INode {
	return &SequenceNode{
		Node: Node{
			children: children,
		},
		idx: 0,
	}
}

func Condition(cond ConditionFunc) INode {
	return &ConditionNode{
		fn: cond,
	}
}

func If(cond ConditionFunc, node INode) INode {
	return Sequence(
		Condition(cond),
		node,
	)
}

func While(cond ConditionFunc, node INode) INode {
	return Parallel(
		Condition(cond),
		node,
	)
}
