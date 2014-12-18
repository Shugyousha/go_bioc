package BioC

// for now, writing InfonStruct, eventually, write Infons

import (
	"encoding/xml"
	"fmt"
)

type Node struct {
	Refid string `xml:"refid,attr"`
	Role string `xml:"role,attr"`
}

func ( node Node ) Write() {
	fmt.Println( "refid: ", node.Refid,
		" role: ", node.Role )
}

type InfonStruct struct {
//	XMLName xml.Name `xml:"infon"`
	Key string `xml:"key,attr"`
	Value string `xml:",chardata"`
}

func ( infonStruct InfonStruct ) Write() {
	fmt.Println( "infon key: ", infonStruct.Key,
		" value: ", infonStruct.Value )
}

func writeInfonStructs( infonStructs []InfonStruct  ) {
	for _, infonStruct := range infonStructs {
		infonStruct.Write()
	}
}


func  write( infons map[string]string ) {
	for key, value := range infons {
//		fmt.Println( "infon key: ", key, " value: ", value )
		fmt.Print( key, ": ", value )
		fmt.Println()
	}
}


type Relation struct {
	Id     string `xml:"id,attr"`
	Infons map[string] string `xml:"-"`
	InfonStructs []InfonStruct `xml:"infon"`
	Nodes []Node `xml:"node"`
}

func ( relate Relation ) Write() {
	fmt.Println( "id:", relate.Id )
	write(relate.Infons)
//	for _, infonStruct := range relate.InfonStructs {
//		infonStruct.Write();
//	}
	for _, node := range relate.Nodes {
		node.Write();
	}
}

func ( relate *Relation ) Map() {
	relate.Infons = make(map[string]string, len(relate.InfonStructs) )
	for _, infon := range relate.InfonStructs {
		relate.Infons[infon.Key] = infon.Value
	}
}

func ( relate *Relation ) Unmap() {
	relate.InfonStructs = nil
	for key, value := range relate.Infons {
		relate.InfonStructs =
			append( relate.InfonStructs, InfonStruct{key,value} )
	}
}


type Location struct {
	Offset int `xml:"offset,attr"`
	Length int `xml:"length,attr"`
}

func ( location Location ) Write() {
	fmt.Println( "offset: ", location.Offset,
		" length: ", location.Length )
}


type Annotation struct {
	Id     string     `xml:"id,attr"`
	Infons map[string] string `xml:"-"`
	InfonStructs []InfonStruct `xml:"infon"`
	Locations []Location `xml:"location"`
	Text   string `xml:"text,omitempty"`
}

func ( note Annotation ) Write() {
	fmt.Println( "id:", note.Id )
	write(note.Infons)
// 	for _, infonStruct := range note.InfonStructs {
// 		infonStruct.Write();
// 	}
	for _, location := range note.Locations {
		location.Write();
	}
	if len(note.Text) > 0 { fmt.Println( "text:", note.Text ) }
}
	
func ( note *Annotation ) Map() {
	note.Infons = make(map[string]string, len(note.InfonStructs) )
	for _, infon := range note.InfonStructs {
		note.Infons[infon.Key] = infon.Value
	}
}

func ( note *Annotation ) Unmap() {
	note.InfonStructs = nil
	for key, value := range note.Infons {
		note.InfonStructs =
			append( note.InfonStructs, InfonStruct{key,value} )
	}
}


type Sentence struct {
	Infons map[string] string `xml:"-"`
	InfonStructs []InfonStruct `xml:"infon"`
	Offset      int `xml:"offset"`
	Text        string `xml:"text,omitempty"`
	Annotations []Annotation `xml:"annotation"`
	Relations   []Relation `xml:"relation"`
}

func ( sent Sentence ) Write() {
	write(sent.Infons)
// 	for _, infonStruct := range sent.InfonStructs {
// 		infonStruct.Write();
// 	}
	fmt.Println( "offset:", sent.Offset )
	if len(sent.Text) > 0 { fmt.Println( "text:", sent.Text ) }
	for _, note := range sent.Annotations {
		note.Write()
	}
	for _, relate := range sent.Relations {
		relate.Write()
	}
}
	
func ( sent *Sentence ) Map() {
	sent.Infons = make(map[string]string, len(sent.InfonStructs) )
	for _, infon := range sent.InfonStructs {
		sent.Infons[infon.Key] = infon.Value
	}
	for i := range sent.Annotations {
		sent.Annotations[i].Map()
	}
	for i := range sent.Relations {
		sent.Relations[i].Map()
	}
}

func ( sent *Sentence ) Unmap() {
	sent.InfonStructs = nil
	for key, value := range sent.Infons {
		sent.InfonStructs =
			append( sent.InfonStructs, InfonStruct{key,value} )
	}
	for i := range sent.Annotations {
		sent.Annotations[i].Unmap()
	}
	for i := range sent.Relations {
		sent.Relations[i].Unmap()
	}
}


type Passage struct {
	Infons map[string]string `xml:"-"`
	InfonStructs []InfonStruct `xml:"infon"`
	Offset      int `xml:"offset"`
	Text        string `xml:"text,omitempty"`
	Sentences   []Sentence `xml:"sentence"`
	Annotations []Annotation `xml:"annotation"`
	Relations   []Relation `xml:"relation"`
}

func ( psg Passage ) Write() {
//	fmt.Println("size of psg.infons: ", len(psg.Infons))
	write(psg.Infons)
// 	for _, infonStruct := range psg.InfonStructs {
// 		infonStruct.Write();
// 	}
//	writeInfonStructs(psg.InfonStructs)
	fmt.Println( "offset:", psg.Offset )
	if len(psg.Text) > 0 { fmt.Println( "text:", psg.Text ) }
	for _, sent := range psg.Sentences {
		sent.Write()
	}
	for _, note := range psg.Annotations {
		note.Write()
	}
	for _, relate := range psg.Relations {
		relate.Write()
	}
}
	
func ( psg *Passage ) Map() {
//	fmt.Println("in psg.Map()")
	psg.Infons = make(map[string]string, len(psg.InfonStructs) )
	for _, infon := range psg.InfonStructs {
		psg.Infons[infon.Key] = infon.Value
	}
//	fmt.Println("Map: size of psg.Infons: ", len(psg.Infons))
	for i := range psg.Sentences {
		psg.Sentences[i].Map()
	}
	for i := range psg.Annotations {
		psg.Annotations[i].Map()
	}
	for i := range psg.Relations {
		psg.Relations[i].Map()
	}
}

func ( psg *Passage ) Unmap() {
	psg.InfonStructs = nil
	for key, value := range psg.Infons {
		psg.InfonStructs =
			append( psg.InfonStructs, InfonStruct{key,value} )
	}
	for i := range psg.Sentences {
		psg.Sentences[i].Unmap()
	}
	for i := range psg.Annotations {
		psg.Annotations[i].Unmap()
	}
	for i := range psg.Relations {
		psg.Relations[i].Unmap()
	}
}


type Document struct {
	XMLName xml.Name `xml:"document"`
	Id string `xml:"id"`
	Infons map[string] string `xml:"-"`
	InfonStructs []InfonStruct `xml:"infon"`
	Passages []Passage `xml:"passage"`
	Relations   []Relation `xml:"relation"`
}

func (doc Document) Write() {
	fmt.Println( "id:", doc.Id )
//	fmt.Println("size of doc.Infons: ", len(doc.Infons))
	write(doc.Infons)
// 	for _, infonStruct := range doc.InfonStructs {
// 		infonStruct.Write();
// 	}
	for _, psg := range doc.Passages {
		psg.Write()
	}
	for _, relate := range doc.Relations {
		relate.Write()
	}
}

	
func ( doc *Document ) Map() {
//	fmt.Println("in doc.Map()")
	doc.Infons = make(map[string]string, len(doc.InfonStructs) )
	for _, infon := range doc.InfonStructs {
		doc.Infons[infon.Key] = infon.Value
	}
//	fmt.Println("Map: size of doc.Infons: ", len(doc.Infons))

	// fmt.Println("before Map: size of doc.psg.Infons: ")
	// for _, psg := range doc.Passages {
	// 	fmt.Println( len(psg.Infons) )
	// }

	for i := range doc.Passages {
		doc.Passages[i].Map()
	}

	// fmt.Println("after Map: size of doc.psg.Infons: ")
	// for _, psg := range doc.Passages {
	// 	fmt.Println( len(psg.Infons) )
	// }

	for i := range doc.Relations {
		doc.Relations[i].Map()
	}
}

func ( doc *Document ) Unmap() {
	doc.InfonStructs = nil
	for key, value := range doc.Infons {
		doc.InfonStructs =
			append( doc.InfonStructs, InfonStruct{key,value} )
	}
	for i := range doc.Passages {
		doc.Passages[i].Unmap()
	}
	for i := range doc.Relations {
		doc.Relations[i].Unmap()
	}
}

 
type Collection struct {
	XMLName xml.Name `xml:"collection"`
	Source  string `xml:"source"`
	Date     string `xml:"date"`
	Key      string `xml:"key"`
	Infons map[string] string `xml:"-"`
	InfonStructs []InfonStruct `xml:"infon"`
	Documents []Document `xml:"document"`
}

/*
func ( infons map[string] string ) Write() {
	for key, value := range infons {
		fmt.Println("infon ", key, ": ", value );
	}
}
*/

func (col Collection) Write() {
	fmt.Println( "source:", col.Source )
	fmt.Println( "date:", col.Date )
	fmt.Println( "key:", col.Key )
//	writeInfonStructs(col.InfonStructs)
//	fmt.Println("size of col.Infons: ", len(col.Infons))
	write(col.Infons)
// 	for _, infonStruct := range col.InfonStructs {
// 		infonStruct.Write();
// 	}
	for _, doc := range col.Documents {
		doc.Write()
	}
}

	
func ( col *Collection ) Map() {
//	fmt.Println("in col.Map()")
	col.Infons = make(map[string]string, len(col.InfonStructs) )
	for _, infon := range col.InfonStructs {
		col.Infons[infon.Key] = infon.Value
	}
	for i := range col.Documents {
		col.Documents[i].Map()
	}
}

func ( col *Collection ) Unmap() {
	col.InfonStructs = nil
	for key, value := range col.Infons {
		col.InfonStructs =
			append( col.InfonStructs, InfonStruct{key,value} )
	}
	for i := range col.Documents {
		col.Documents[i].Unmap()
	}
}
