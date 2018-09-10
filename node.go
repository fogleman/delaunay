package delaunay

type node struct {
	p       Point
	i       int
	t       int
	removed bool
	prev    *node
	next    *node
}

func newNode(nodes []node, p Point, i int, prev *node) *node {
	n := &nodes[i]
	n.p = p
	n.i = i
	n.removed = false
	if prev == nil {
		n.prev = n
		n.next = n
	} else {
		n.next = prev.next
		n.prev = prev
		prev.next.prev = n
		prev.next = n
	}
	return n
}

func (n *node) remove() *node {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.removed = true
	return n.prev
}
