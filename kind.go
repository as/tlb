package main

type Kind int
type Kinder interface {
	Kind() Kind
}

func (k Kind) L() Kind { return k & 0xf }

func (v imp) Kind() Kind       { return Kind(0) }
func (v name) Kind() Kind      { return Kind(0) }
func (v ref) Kind() Kind       { return Kind(0) }
func (v gid) Kind() Kind       { return Kind(0) }
func (v namemap) Kind() Kind   { return Kind(0) }
func (v gidmap) Kind() Kind    { return Kind(0) }
func (v typeinfo) Kind() Kind  { return Kind(0) }
func (v typedesc) Kind() Kind  { return Kind(0) }
func (v arraydesc) Kind() Kind { return Kind(0) }
func (v custom) Kind() Kind    { return Kind(0) }
func (v customGID) Kind() Kind { return Kind(0) }
