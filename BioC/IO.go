package BioC

// for now, writing InfonStruct, eventually, write Infons

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"reflect"
)

func ReadCollection(filename string) Collection {
	var col Collection

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when opening file.\n")
	}

	breader := bufio.NewReader(file)
	decoder := xml.NewDecoder(breader)

	err = decoder.Decode(&col)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when decoding file.\n")
	}

	col.Map()
	return col
}

func WriteCollection(col Collection, filename string) error {
	col.Unmap()

	perm := os.FileMode(0664)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	bwriter := bufio.NewWriter(f)
	encoder := xml.NewEncoder(bwriter)
	encoder.Indent(" ", "")

	err = encoder.Encode(col)
	if err != nil {
		return err
	}

	err = bwriter.Flush()
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return err
}

type ReadDocument struct {
	inDocument   bool
	inCollection bool
	token        xml.Token
	decoder      *xml.Decoder
	err          string
}

func (rd *ReadDocument) Start(file string) (Collection, string) {

	var col Collection

	xmlFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	rd.decoder = xml.NewDecoder(xmlFile)

	rd.inCollection = false
	rd.inDocument = false

	for !rd.inDocument {

		rd.token, err = rd.decoder.Token()
		if rd.token == nil {
			return col, "no collection"
		}
		if err != nil {
			panic(err)
		}

		switch se := rd.token.(type) {
		case xml.StartElement:
			//			fmt.Println( se.Name.Local )

			switch se.Name.Local {
			case "collection":
				rd.inCollection = true
				//				fmt.Println( "found collection" )
			case "source":
				if rd.inCollection {
					rd.token, err = rd.decoder.Token()
					col.Source = string(rd.token.(xml.CharData))
					err = rd.decoder.Skip()
					//				fmt.Println( "found source" )
				}
			case "date":
				if rd.inCollection {
					rd.token, err = rd.decoder.Token()
					col.Date = string(rd.token.(xml.CharData))
					err = rd.decoder.Skip()
					//				fmt.Println( "found date" )

				}
			case "key":
				if rd.inCollection {
					rd.token, err = rd.decoder.Token()
					col.Key = string(rd.token.(xml.CharData))
					err = rd.decoder.Skip()
					//				fmt.Println( "found key" )
				}
			case "infon":
				if rd.inCollection {
					key := ""
					for i := range se.Attr {
						if se.Attr[i].Name.Local == "key" {
							key = se.Attr[i].Value
							break
						}
					}
					if key == "" {
						return col, "infon without key"
					}
					rd.token, err = rd.decoder.Token()
					value := string(rd.token.(xml.CharData))
					err = rd.decoder.Skip()
					col.InfonStructs = append(col.InfonStructs, InfonStruct{key, value})
				}

			case "document":
				if rd.inCollection {
					rd.inDocument = true
					//					fmt.Println( "found document" )
					col.Map()
					return col, ""
				}
			}
		}
	}
	col.Map()
	return col, "end of Start"
}

func (rd *ReadDocument) Next() (Document, string) {

	var doc Document

	if !rd.inCollection {
		return doc, "not in collection"
	}

	if !rd.inDocument {
		return doc, "not in document"
	}

	//doc_loop:	for {

	switch se := rd.token.(type) {
	case xml.StartElement:
		if se.Name.Local == "document" {

			se, _ := rd.token.(xml.StartElement)
			rd.decoder.DecodeElement(&doc, &se)

			//			doc.Write()

			token, err := rd.decoder.Token()
			if err != nil {
				panic(err)
			}
			rd.token = token
			if false {
				fmt.Println(rd.token)
				fmt.Println(reflect.TypeOf(rd.token))
			}

			doc.Map()
			return doc, ""

		}
	case xml.EndElement:
		if se.Name.Local == "collection" {
			rd.inDocument = false
			return doc, "eof"
		}
	}
	return doc, "end of next"
}

type WriteDocument struct {
	XMLFile *os.File
}

//  based on MarshalIndent
// needs more digging getting type info from BioC struct
// 	/*
// 	 func writeElement(v interface{} ) ([]byte, error) {
// 	 var b bytes.Buffer
// 	 enc := NewEncoder(&b)
// 	 err := enc.marshalValue(reflect.ValueOf(v), nil)
// 	 enc.Flush()
// 	 if err != nil {
// 	 return nil, err
// 	 }
// 	 return b.Bytes(), nil
// 	 }
// 	 */

func (wr *WriteDocument) writeString(s string) {
	wr.XMLFile.Write([]byte(s))
}

func (wr *WriteDocument) writeElement(val, name string) {

	wr.writeString("<")
	wr.writeString(name)
	wr.writeString(">")

	wr.writeString(val)

	wr.writeString("<")
	wr.writeString("/")
	wr.writeString(name)
	wr.writeString(">")

}

func (wr *WriteDocument) writeElementAttr(val, name string, attrs []xml.Attr) {

	wr.writeString("<")
	wr.writeString(name)

	for i := range attrs {
		wr.writeString(" ")
		wr.writeString(attrs[i].Name.Local)
		wr.writeString("=\"")
		xml.EscapeText(wr.XMLFile, []byte(attrs[i].Value))
		wr.writeString("\"")
	}

	wr.writeString(">")

	wr.writeString(val)

	wr.writeString("<")
	wr.writeString("/")
	wr.writeString(name)
	wr.writeString(">")

}

func (wr *WriteDocument) Start(file string, col Collection) error {

	var err error
	wr.XMLFile, err = os.Create(file)
	if err != nil {
		panic(err)
	}

	// next 4 lines in WriteCollection also

	_, err = wr.XMLFile.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	_, err = wr.XMLFile.Write([]byte("<!DOCTYPE collection SYSTEM 'BioC.dtd'>"))
	if err != nil {
		return err
	}

	//	startElement(wr.XMLFile, col, "collection" )
	//	wr.XMLFile.Write( []byte( "<collection>" ) )
	wr.writeString("<collection>")

	/*
		data, err := xml.Marshal( col.Source )
		wr.XMLFile.Write( data )
		data, err = xml.Marshal( col.Date )
		wr.XMLFile.Write( data )
		data, err = xml.Marshal( col.Key )
		wr.XMLFile.Write( data )
	*/

	wr.writeElement(col.Source, "source")
	wr.writeElement(col.Date, "date")
	wr.writeElement(col.Key, "key")

	//	col.Unmap() -- not needed because using Infons directly here
	for key, value := range col.Infons {
		attrs := []xml.Attr{xml.Attr{xml.Name{"", "key"}, key}}
		wr.writeElementAttr(value, "infon", attrs)
	}

	return err
}

func (wr *WriteDocument) Next(doc Document) {
	doc.Unmap()
	data, err := xml.Marshal(doc)
	if err != nil {
		panic(err)
	}
	wr.XMLFile.Write(data)
}

func (wr *WriteDocument) Close() {
	wr.XMLFile.Write([]byte("</collection>"))
	wr.XMLFile.Close()
}
