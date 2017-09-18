package main

import "fmt"

var (
	sp = fmt.Sprintf
	sl = fmt.Sprintln
)

// swap group with flag,
//func (v typeinfo) String() string {
//	return sp("Type:	group=%-14d n=%x kind=%-8s	nameat=%-8x n=%-8x nvar=%-8x nfunc=%-8x gidat=%-8x nvirt=%-8x nimport=%-8x descat=%-8x ",
//		 v.group,v.n, ComType(v.kind), v.nameat , v.n, v.nvar, v.nfunc, v.gidat, v.nvirt, v.nimport, v.descat)
//}
// swap group with flag
func (v typeinfo) String() string { return fmt.Sprintf("%#v", v) }

func (v imp) String() string {
	return sp("imp:	flag=%d	fileat=%d	gidat=%s", v.flag, v.fileat, v.gidat)
}
func (v impfile) String() string {
	return sp("impfile:	data=%d", v.data)
}
func (v ref) bstr() string {
	return sp("ref:	kind=%v	flag=%d	customat=%d	nextat=%d", VarKind(v.kind), v.flag, v.customat, v.nextat)
}
func (v gid) String() string {
	return sp("gid:	gid=%d	kind=%v	nextat=%d", v.manure, VarKind(v.kind), v.nextat)
}
func (v gidmap) String() string {
	return sp("gidmap:	a=%d	b=%d	", v.a, v.b)
}
func (v typedesc) String() string {
	return sp("typedesc:	group=%-14d  u1=%s	b=%s	c=%s", v.group, VarKind(v.u1&VarScalar), VarKind(v.b&VarScalar), VarKind(v.c&VarScalar))
}
func (v arraydesc) String() string {
	return sp("arraydesc:	a=%d	b=%d	c=%d", v.a, v.b, v.c)
}
func (v bstr) String() string {
	return sp("bstring:	n=%d	string=%s", v.n, v.data)
}
func (v name) String() string {
	return sp("name: group=%-14d flag=%x nextat=%d string=%s ", v.group, v.flag, v.nextat, v.data)
}
func (v namemap) String() string {
	return sp("namemap:	a=%d	b=%d	c=%d", v.a, v.b, v.c)
}
func (v custom) String() string {
	return sp("custom:	a=%d	b=%d	c=%d", v.a, v.b, v.c)
}
func (v customGID) String() string {
	return sp("customGID: gid=%d	kind=%v	nextat=%d", v.manure, VarKind(v.kind), v.nextat)
}

//go:generate stringer -type ComType
