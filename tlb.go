package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// BinaryReader implements methods required for initializing
// itself from a binary stream. This interface accepts an io.Reader
// instead of []byte, compare it to binary/encoding BinaryMarshaler.
type BinaryReader interface {
	ReadBinary(io.Reader) error
}

// Node represents a meaningful datum in a tlb
type Node interface {
	BinaryReader
	//process() error
}

// SecNum enumerates the possible sections found in a tlb file. A
// microsoft tlb contains 13 relevant sections and 15 total
type SecNum int

const (
	SecTypeInfo SecNum = iota
	SecImp
	SecImpFile
	SecRef
	SecGIDMap
	SecGID
	SecNameMap
	SecName
	SecBstr
	SecTypeDesc
	SecArrayDesc
	SecCustom
	SecCustomGID
	SecTrash1
	SecTrash2
	NSections = 15
)

// Section provides structured access to the data in a tlb section.
// Data from the section may be read directly using ReaderAt
type Section struct {
	ID SecNum
	io.ReaderAt
	src   *io.SectionReader
	head  *secthead
	words uint32
	data  BinaryReader
}

type Typelib struct {
	head       *head
	sectlist   *sectlist
	fd         *os.File
	Section    []*Section
	idmap      map[uint32][]Node
	am, bm, cm map[uint32][]Node
	hashmap    map[uint32][]Node
	namemap    map[uint32]string
	group      map[uint32][]Node
	pos        uint32
	Node       []Node
	Heap
}

//go:generate wire9 -f tlb_wire9.go tlb.go

//wire9 head     msft[4] majmin[4] gidat[4] lang[8] flag[4] ver[4] flag2[4] ntype[4] he1[4]  he2[4] context[4] nname[4] nnamebyte[4] nameat[4] importat[4] customat[4] r3[8] dispatchat[4] nimport[4]
//wire9 secthead at[4] n[4] fu[4] u[4]
//wire9 sectlist s[15,[]secthead]
//wire9 arraydesc  a[4] b[4] c[4]

//wire9 typeinfo kind[4] at[4] u1[16] nvar[2] nfunc[2] u2[20] gidat[4] group[4] nameat[4] ver[4] he1[4] he2[4] cdat[4] nimport[2] nvirt[2] n[4] descat[4] r1[4] pad[8]
//wire9 imp        flag[4] fileat[4] gidat[2]
//wire9 impfile    data[1]
//wire9 ref        kind[4] flag[4] customat[4] nextat[4]
//wire9 gid        manure[16] kind[4] nextat[4]
//wire9 gidmap     a[4] b[4]
//wire9 name       group[4] nextat[4] n[1] flag[1] hash[2] data[n]
//wire9 namemap    a[4] b[4] c[4]
//wire9 bstring    n[2] data[n]
//wire9 typedesc   group[2] u1[2] b[4] c[4]
//wire9 custom     a[4] b[4] c[4]
//wire9 customGID  manure[16] kind[4] nextat[4]

//wire9 secTypeinfo   n[4] list[n,[]typeinfo]
//wire9 secImp        n[4] list[n,[]imp]
//wire9 secImpfile    n[4] list[n,[]impfile]
//wire9 secRef        n[4] list[n,[]ref]
//wire9 secGid        n[4] list[n,[]gid]
//wire9 secGidmap     n[4] list[n,[]gidmap]
//wire9 secName       n[4] list[n,[]name]
//wire9 secNamemap    n[4] list[n,[]namemap]
//wire9 secBstring    n[4] list[n,[]bstring]
//wire9 secArraydesc  n[4] list[n,[]arraydesc]
//wire9 secTypedesc   n[4] list[n,[]typedesc]
//wire9 secCustom     n[4] list[n,[]custom]
//wire9 secCustomGID  n[4] list[n,[]customGID]

//wire9 typedata n[4] data[n-1,[]byte]
//wire9 funcdesc n[2] id[2] t1[2] t2[2] u1[1] u2[2] vt1[2] size[2] id2[1] numxconv[1]  defxcust[1] invxkind[1]
//wire9 Func     info[4] kind[4] flag[4] virtat[2]  n[2] flag2[4] nparam[2] nopt[2] he1[4] he2[4] entryat[4] he3[4] cdat[4] cdparam[4]
//wire9 Var      info[4] kind[4] flag[4] varkind[2] n[2] value[4] he1[4] he2[4] u1[4] customat[4] he3[4]
//wire9 Param    kind[4] nameat[4] flag[4]

type wcReader struct {
	nread int
	nbyte int
	r     io.Reader
}

// Open returns a new ReadSeeker reading the tlb section.
func (s *Section) Open() io.ReadSeeker { return io.NewSectionReader(s.src, 0, 1<<63-1) }

// :/func.*wordsizeof/

// New returns a pointer to an initialized struct representing
// the enumerated section.
func (s SecNum) Make() (sec BinaryReader) {
	switch s {
	case SecTypeInfo:
		return new(secTypeinfo)
	case SecImp:
		return new(secImp)
	case SecImpFile:
		return new(secImpfile)
	case SecRef:
		return new(secRef)
	case SecGIDMap:
		return new(secGidmap)
	case SecGID:
		return new(secGid)
	case SecNameMap:
		return new(secNamemap)
	case SecName:
		return new(secName)
	case SecBstr:
		return new(secBstring)
	case SecTypeDesc:
		return new(secTypedesc)
	case SecArrayDesc:
		return new(secArraydesc)
	case SecCustom:
		return new(secCustom)
	case SecCustomGID:
		return new(secCustomGID)
	}
	return nil
}

func (t tab) Println(i ...interface{}) {
	indent := []interface{}{strings.Repeat("	", int(tracer))}
	log.Println(append(indent, i...)...)
}
func trace(s string) string {
	tracer.Println("Enter", s)
	tracer++
	return s
}
func un(s string) {
	tracer--
	tracer.Println("Exits", s, "\n")
}

func (t *Typelib) ReadSections() (err error) {
	defer un(trace("ReadSections"))

	buf := new(bytes.Buffer)
	injectread := func(n uint32, section BinaryReader, src io.Reader) (err error) {
		if err = binary.Write(buf, binary.LittleEndian, n); err != nil {
			return err
		}
		if err = section.ReadBinary(io.MultiReader(buf, src)); err != nil {
			if err != io.EOF {
				return err
			}
		}
		return nil
	}
	hex := func(i interface{}) string {
		return fmt.Sprintf("%x", i)
	}

	for i, v := range t.sectlist.s {
		if i == 13 {
			break
		}
		words := t.wordsizeof(SecNum(i))
		if words == 0 {
			tracer.Println("no words for section", i)
			continue
		}
		tracer.Println(v)
		s := &Section{
			ID:    SecNum(i),
			src:   io.NewSectionReader(t.fd, int64(v.at), int64(v.n)),
			head:  &v, // creates a copy, capturing v
			words: words,
		}
		sectiondata := s.ID.Make()
		s1 := v.at
		e1 := v.at + v.n
		xs := (((s1 / 16) + 1) * 3 * 3) + (s1 * 3)
		xe := (((e1 / 16) + 1) * 3 * 3) + (e1 * 3)

		func() {

			tracer.Println("SecNum(i)",
				fmt.Sprintf("/usr/as/tlb/14.hex:#%d,#%d",
					xs, xe,
				),
				"start", hex(s1),
				"size", hex(v.n),
				"end", hex(e1),
			)

			if err = injectread(s.words, sectiondata, s.src); err != nil {
				return
			}
			s.data = sectiondata
			t.Section = append(t.Section, s)
			tracer.Println("Added section:", t.Section[SecNum(i)].data, "\n\n")
		}()
		if err != nil {
			return
		}
	}
	// The following GUID is for the ID of the typelib if this project is exposed to COM
	//[assembly: Guid("7fad1542-e92e-4402-8e9a-b852bddc45a2")]
	for _, v := range t.Section {
		if err = sectionErr(v); err != nil {
			return
		}
	}
	return t.checkSections()
}
func (t *Typelib) parseTypeInfo() (err error) {
	defer un(trace("parseTypeInfo"))

	return
}

func (t *Typelib) checkSections() (err error) {
	defer un(trace("checkSections"))
	if err = t.checkTypeInfo(); err != nil {
		return
	}
	if err = t.checkTypeDesc(); err != nil {
		return
	}
	if err = t.checkNames(); err != nil {
		return
	}
	if err = t.checkNameMap(); err != nil {
		return
	}
	if err = t.checkGIDMap(); err != nil {
		return
	}
	sec := t.Section[SecTypeInfo].data.(*secTypeinfo)
	for _, v := range sec.list {
		func() {
			defer un(trace(fmt.Sprintf("%v", v)))
			for j, w := range t.hashmap[v.group] {
				tracer.Println(j, w)
			}
		}()
	}
	return
}

func (t *Typelib) checkTypeInfo() (err error) {
	defer un(trace("checkTypeInfo"))
	sec := t.Section[SecTypeInfo].data.(*secTypeinfo)
	for i, v := range sec.list {
		sec.list[i].kind &= 0x0f
		tracer.Println("add", v.group, "as parent")
		td := &typedata{}
		err = td.ReadBinary(io.NewSectionReader(t.fd, int64(v.at), 1<<62))
		if err != nil {
			return err
		}
		tracer.Println("typeinfo data:",
			fmt.Sprintf("%d %d %x",
				td.n,
				td.data[0],
				td.data,
			),
		)

		t.hashmap[v.group] = append(t.hashmap[v.group], &sec.list[i])
	}
	return nil
}

func (t *Typelib) checkTypeDesc() (err error) {
	defer un(trace("checkTypeDesc"))
	sec := t.Section[SecTypeDesc].data.(*secTypedesc)
	for i := range sec.list {
		v := &sec.list[i]
		func() {
			p := t.findParent(uint32(v.group))
			t.hashmap[p] = append(t.hashmap[p], &sec.list[i])
		}()
	}
	return nil
}

func (t *Typelib) findParent(x uint32) (n uint32) {
	defer un(trace("findParent"))
	for _, w := range t.Section[SecTypeInfo].data.(*secTypeinfo).list {
		g := w.group
		if x < g {
			break
		}
		n = g
	}
	tracer.Println("the parent of", x, "is", n)
	return n
}

func (t *Typelib) checkNames() (err error) {
	defer un(trace("checkNames"))
	sec := t.Section[SecName].data.(*secName)
	for i, v := range sec.list {
		p := t.findParent(v.group)
		t.hashmap[p] = append(t.hashmap[p], &sec.list[i])
	}
	return nil
}
func (t *Typelib) checkGIDMap() (err error) {
	defer un(trace("checkGIDMap"))
	sec := t.Section[SecGIDMap].data.(*secGidmap)
	for i := range sec.list {
		v := &sec.list[i]
		v.a >>= 16
		v.b >>= 16
		tracer.Println("checkGIDMap", v)
	}
	return nil
}
func (t *Typelib) checkGID() (err error) {
	defer un(trace("checkGID"))
	sec := t.Section[SecGID].data.(*secGid)
	for i := range sec.list {
		v := &sec.list[i]
		tracer.Println("checkGIDMap", v)
	}
	return nil
}

func (t *Typelib) checkNameMap() (err error) {
	defer un(trace("checkNameMap"))
	return nil
}

func (t *Typelib) checkEnumDesc(ti *typeinfo) (err error) {
	//subtype := ti.iface
	//err = bind(subtype, name, hash, flag, ppTInfo, kind, bindptr)
	//
	//if err != nil && desckind == nil {}
	// check 4 msismatch
	return fmt.Errorf("dont call this")
}
func sectionErr(s *Section) error {
	if s == nil || s.data == nil {
		if s == nil {
			return fmt.Errorf("bad section: nil")
		}
		return fmt.Errorf("bad section: data list nil")
	}
	return nil
}

func (t *Typelib) Itemize(n Node, foundat uint32) error {
	tracer.Println("Element found at address", foundat, n)
	t.Heap.Alloc(n, foundat)
	return nil
}

func (t *Typelib) Lookup(k Kind) []Node {
	return t.idmap[uint32(k)]
}

// :/type.SecNum/

func (t *Typelib) wordsizeof(id SecNum) (n uint32) {
	defer func() {
	}()
	size := t.sectlist.s[id].n
	switch id {
	case SecTypeInfo:
		n = t.head.ntype
	case SecImp:
		n = t.head.nimport
	case SecImpFile:
		n = size
	case SecGID:
		n = size / 24
	case SecName:
		n = t.head.nname
	case SecBstr:
		n = size / 2
	case SecTypeDesc:
		n = t.head.ntype
	case SecRef:
		n = size
	case SecArrayDesc:
		n = size
	case SecCustom:
		n = size
	case SecCustomGID:
		n = size
	case SecGIDMap:
		n = 2
	case SecNameMap:
		n = size
	}
	tracer.Println("Typelib: wordsize of:", id, n)
	return n
}

func NewTypelib(name string) (*Typelib, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	t := &Typelib{
		fd:       fd,
		head:     &head{},
		sectlist: &sectlist{},
		idmap:    make(map[uint32][]Node),
		hashmap:  make(map[uint32][]Node),
		Heap:     NewHeap(),
	}
	if err = t.head.ReadBinary(fd); err != nil {
		return nil, err
	}
	size := int64(t.head.ntype*4 + 80 + 4)
	if err = t.sectlist.ReadBinary(io.NewSectionReader(fd, size, 0xfffff)); err != nil {
		return nil, err
	}
	return t, err
}

func newWCReader(r io.Reader) *wcReader {
	return &wcReader{r: r}
}

func (wc *wcReader) Read(p []byte) (n int, err error) {
	if n, err = wc.r.Read(p); err != nil {
		return
	}
	wc.nread++
	wc.nbyte += n
	return
}

func (s Section) Size() uint32 {
	return uint32(s.src.Size())
}
func (s Section) Words() (n uint32) {
	return s.words
}

func (s secthead) bstr() string {
	return fmt.Sprintf("seg: at=%d n=%d", s.at, s.n)
}

/*
	>awk -F'[: ]' '{ gsub("Sect","",$0); printf "case Sect%s: return &%s{}\n", $2, $2}'
	>awk -F'[: ]' '{ gsub("Sect","",$0); printf "case %s: return &%s{}\n", $2, $2}'
*/

// Post-processing is done after each section element is read
// in check for corruption and inconsistency.

func (z *head) process() error     { return nil }
func (z *secthead) process() error { return nil }
func (z *sectlist) process() error { return nil }

func (z *typeinfo) process() (err error) {
	z.kind &= 0xf

	/*
		var td *TypeDesc
		seekto := func(addr int) io.ReaderAt {
			return z.NewReaderAt(addr)
		}
		if z.kind == ComAlias {
			if err = z.up.parseTypeDesc(seekto(z.descat)); err != nil{
				return
			}
		}

		if err = z.up.parseName(seekto(z.nameat)); err != nil{
			return
		}
		if err = z.up.parseHref(seekto(1)); err != nil{
			return
		}
		if err = z.up.parseDoc(seekto(1)); err != nil{
			return
		}
		if err = z.up.parseHelp(seekto(1)); err != nil{
			return
		}
		if err = z.up.parseGID(seekto(z.gidat)); err != nil{
			return
		}
		if err = td.parseVar(seekto(z.aux)); err != nil{
			return
		}
		switch z.kind {
			case VarCoClass:
				err = z.up.parseCoClass(seekto(descat))
			case Dispatch:
				if z.descat != 0xffffffff {
				}
			default:
		}
		if err != nil {
			return err
		}
		if err = parseCustomData(seekto(cdat)); err != nil{
			return err
		}
		z.TypeDesc = td
	*/
	return nil
}

func (z *TypeDesc) process() error {
	return nil
}

func (z *Func) process() error      { return nil }
func (z *Var) process() error       { return nil }
func (z *Param) process() error     { return nil }
func (z *imp) process() error       { return nil }
func (z *impfile) process() error   { return nil }
func (z *ref) process() error       { return nil }
func (z *gid) process() error       { return nil }
func (z *gidmap) process() error    { return nil }
func (z *name) process() error      { return nil }
func (z *namemap) process() error   { return nil }
func (z *bstr) process() error      { return nil }
func (z *arraydesc) process() error { return nil }
func (z *custom) process() error    { return nil }
func (z *customGID) process() error { return nil }

type tab int

var tracer tab
