package BioC

import (
	"encoding/xml"
	"fmt"
)

type Node struct {
	Refid string `xml:"refid,attr"`
	Role  string `xml:"role,attr"`
}

func (node Node) Write() {
	fmt.Println("refid: ", node.Refid,
		" role: ", node.Role)
}

type InfonStruct struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",chardata"`
}

func (infonStruct InfonStruct) Write() {
	fmt.Println("infon key: ", infonStruct.Key,
		" value: ", infonStruct.Value)
}

type Relation struct {
	Id           string            `xml:"id,attr"`
	Infons       map[string]string `xml:"-"`
	InfonStructs []InfonStruct     `xml:"infon"`
	Nodes        []Node            `xml:"node"`
}

type relationTmp struct {
	Id           string        `xml:"id,attr"`
	InfonStructs []InfonStruct `xml:"infon"`
	Nodes        []Node        `xml:"node"`
}

func (r *Relation) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var reltmp relationTmp

	err = d.DecodeElement(&reltmp, &start)
	if err != nil {
		return
	}
	r.Id = reltmp.Id
	r.InfonStructs = reltmp.InfonStructs

	r.Infons = make(map[string]string, len(reltmp.InfonStructs))

	for _, s := range reltmp.InfonStructs {
		r.Infons[s.Key] = s.Value
	}

	return
}

func (relate Relation) Write() {
	fmt.Println("id:", relate.Id)
	for _, node := range relate.Nodes {
		node.Write()
	}
}

type Location struct {
	Offset int `xml:"offset,attr"`
	Length int `xml:"length,attr"`
}

func (location Location) Write() {
	fmt.Println("offset: ", location.Offset,
		" length: ", location.Length)
}

type Annotation struct {
	Id           string            `xml:"id,attr"`
	Infons       map[string]string `xml:"-"`
	InfonStructs []InfonStruct     `xml:"infon"`
	Locations    []Location        `xml:"location"`
	Text         string            `xml:"text,omitempty"`
}

type annotationTmp struct {
	Id           string        `xml:"id,attr"`
	InfonStructs []InfonStruct `xml:"infon"`
	Locations    []Location    `xml:"location"`
	Text         string        `xml:"text,omitempty"`
}

func (a *Annotation) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var annotmp annotationTmp

	err = d.DecodeElement(&annotmp, &start)
	if err != nil {
		return
	}

	a.Id = annotmp.Id
	a.Locations = annotmp.Locations
	a.InfonStructs = annotmp.InfonStructs
	a.Infons = make(map[string]string, len(annotmp.InfonStructs))

	for _, s := range annotmp.InfonStructs {
		a.Infons[s.Key] = s.Value
	}

	return
}

func (note Annotation) Write() {
	fmt.Println("id:", note.Id)
	for _, location := range note.Locations {
		location.Write()
	}
	if len(note.Text) > 0 {
		fmt.Println("text:", note.Text)
	}
}

type Sentence struct {
	Infons       map[string]string `xml:"-"`
	InfonStructs []InfonStruct     `xml:"infon"`
	Offset       int               `xml:"offset"`
	Text         string            `xml:"text,omitempty"`
	Annotations  []Annotation      `xml:"annotation"`
	Relations    []Relation        `xml:"relation"`
}

type sentenceTmp struct {
	InfonStructs []InfonStruct `xml:"infon"`
	Offset       int           `xml:"offset"`
	Text         string        `xml:"text,omitempty"`
	Annotations  []Annotation  `xml:"annotation"`
	Relations    []Relation    `xml:"relation"`
}

func (s *Sentence) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var senttmp sentenceTmp

	err = d.DecodeElement(&senttmp, &start)
	if err != nil {
		return
	}

	s.Offset = senttmp.Offset
	s.Text = senttmp.Text
	s.Annotations = senttmp.Annotations
	s.Relations = senttmp.Relations
	s.InfonStructs = senttmp.InfonStructs
	s.Infons = make(map[string]string, len(senttmp.InfonStructs))

	for _, str := range senttmp.InfonStructs {
		s.Infons[str.Key] = str.Value
	}

	return
}

func (sent Sentence) Write() {
	fmt.Println("offset:", sent.Offset)
	if len(sent.Text) > 0 {
		fmt.Println("text:", sent.Text)
	}
	for _, note := range sent.Annotations {
		note.Write()
	}
	for _, relate := range sent.Relations {
		relate.Write()
	}
}

type Passage struct {
	Infons       map[string]string `xml:"-"`
	InfonStructs []InfonStruct     `xml:"infon"`
	Offset       int               `xml:"offset"`
	Text         string            `xml:"text,omitempty"`
	Sentences    []Sentence        `xml:"sentence"`
	Annotations  []Annotation      `xml:"annotation"`
	Relations    []Relation        `xml:"relation"`
}

type passageTmp struct {
	InfonStructs []InfonStruct `xml:"infon"`
	Offset       int           `xml:"offset"`
	Text         string        `xml:"text,omitempty"`
	Sentences    []Sentence    `xml:"sentence"`
	Annotations  []Annotation  `xml:"annotation"`
	Relations    []Relation    `xml:"relation"`
}

func (p *Passage) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var psgtmp passageTmp

	err = d.DecodeElement(&psgtmp, &start)
	if err != nil {
		return
	}

	p.Offset = psgtmp.Offset
	p.Text = psgtmp.Text
	p.Annotations = psgtmp.Annotations
	p.Relations = psgtmp.Relations
	p.Sentences = psgtmp.Sentences
	p.InfonStructs = psgtmp.InfonStructs
	p.Infons = make(map[string]string, len(psgtmp.InfonStructs))

	for _, str := range psgtmp.InfonStructs {
		p.Infons[str.Key] = str.Value
	}

	return
}

func (psg Passage) Write() {
	fmt.Println("offset:", psg.Offset)
	if len(psg.Text) > 0 {
		fmt.Println("text:", psg.Text)
	}
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

type Document struct {
	XMLName      xml.Name          `xml:"document"`
	Id           string            `xml:"id"`
	Infons       map[string]string `xml:"-"`
	InfonStructs []InfonStruct     `xml:"infon"`
	Passages     []Passage         `xml:"passage"`
	Relations    []Relation        `xml:"relation"`
}

type documentTmp struct {
	XMLName      xml.Name      `xml:"document"`
	Id           string        `xml:"id"`
	InfonStructs []InfonStruct `xml:"infon"`
	Passages     []Passage     `xml:"passage"`
	Relations    []Relation    `xml:"relation"`
}

func (doc *Document) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var doctmp documentTmp

	err = d.DecodeElement(&doctmp, &start)
	if err != nil {
		return
	}

	doc.Id = doctmp.Id
	doc.Passages = doctmp.Passages
	doc.Relations = doctmp.Relations
	doc.InfonStructs = doctmp.InfonStructs
	doc.Infons = make(map[string]string, len(doctmp.InfonStructs))

	for _, str := range doctmp.InfonStructs {
		doc.Infons[str.Key] = str.Value
	}

	return
}

func (doc Document) Write() {
	fmt.Println("id:", doc.Id)
	for _, psg := range doc.Passages {
		psg.Write()
	}
	for _, relate := range doc.Relations {
		relate.Write()
	}
}

type Collection struct {
	XMLName      xml.Name          `xml:"collection"`
	Source       string            `xml:"source"`
	Date         string            `xml:"date"`
	Key          string            `xml:"key"`
	Infons       map[string]string `xml:"-"`
	InfonStructs []InfonStruct     `xml:"infon"`
	Documents    []Document        `xml:"document"`
}

type collectionTmp struct {
	XMLName      xml.Name      `xml:"collection"`
	Source       string        `xml:"source"`
	Date         string        `xml:"date"`
	Key          string        `xml:"key"`
	InfonStructs []InfonStruct `xml:"infon"`
	Documents    []Document    `xml:"document"`
}

func (col *Collection) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var coltmp collectionTmp

	err = d.DecodeElement(&coltmp, &start)
	if err != nil {
		return
	}

	col.Source = coltmp.Source
	col.Date = coltmp.Date
	col.Key = coltmp.Key
	col.Documents = coltmp.Documents
	col.InfonStructs = coltmp.InfonStructs
	col.Infons = make(map[string]string, len(coltmp.InfonStructs))

	for _, str := range coltmp.InfonStructs {
		col.Infons[str.Key] = str.Value
	}

	return
}

func (col Collection) Write() {
	fmt.Println("source:", col.Source)
	fmt.Println("date:", col.Date)
	fmt.Println("key:", col.Key)

	for _, doc := range col.Documents {
		doc.Write()
	}
}
