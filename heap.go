package main

type Heap struct {
	m map[uint32]Node
}

func (h *Heap) Alloc(n Node, at uint32) {
	h.m[at] = n
}

func (h *Heap) At(at uint32) Node {
	return h.m[at]
}

func NewHeap() Heap {
	return Heap{
		make(map[uint32]Node),
	}
}
