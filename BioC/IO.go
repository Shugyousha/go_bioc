package BioC

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"os"
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

type DocumentReader struct {
	inDocument   bool
	inCollection bool
	token        xml.Token
	decoder      *xml.Decoder
}

func (dr *DocumentReader) Start(reader io.Reader) (Collection, error) {
	var col Collection
	var err error

	dr.decoder = xml.NewDecoder(reader)

	dr.inCollection = false
	dr.inDocument = false

	for !dr.inDocument {
		dr.token, err = dr.decoder.Token()
		if dr.token == nil {
			return col, fmt.Errorf("no collection")
		}
		if err != nil {
			return col, fmt.Errorf("Error when decoding token: %s", err)
		}

		switch se := dr.token.(type) {
		case xml.StartElement:

			switch se.Name.Local {
			case "collection":
				dr.inCollection = true
			case "source":
				if dr.inCollection {
					dr.token, err = dr.decoder.Token()
					if err != nil {
						return col, fmt.Errorf("Error when decoding source token: %s", err)
					}
					col.Source = string(dr.token.(xml.CharData))
					err = dr.decoder.Skip()
				}
			case "date":
				if dr.inCollection {
					dr.token, err = dr.decoder.Token()
					if err != nil {
						return col, fmt.Errorf("Error when decoding date token: %s", err)
					}
					col.Date = string(dr.token.(xml.CharData))
					err = dr.decoder.Skip()

				}
			case "key":
				if dr.inCollection {
					dr.token, err = dr.decoder.Token()
					if err != nil {
						return col, fmt.Errorf("Error when decoding key token: %s", err)
					}
					col.Key = string(dr.token.(xml.CharData))
					err = dr.decoder.Skip()
				}
			case "infon":
				if dr.inCollection {
					key := ""
					for i := range se.Attr {
						if se.Attr[i].Name.Local == "key" {
							key = se.Attr[i].Value
							break
						}
					}
					if key == "" {
						return col, fmt.Errorf("infon without key")
					}
					dr.token, err = dr.decoder.Token()
					if err != nil {
						return col, fmt.Errorf("Error when decoding key token: %s", err)
					}
					value := string(dr.token.(xml.CharData))
					err = dr.decoder.Skip()
					col.InfonStructs = append(col.InfonStructs, InfonStruct{key, value})
				}

			case "document":
				if dr.inCollection {
					dr.inDocument = true
					col.Map()
					return col, nil
				}
			}
		}
	}
	col.Map()
	return col, nil
}

func (dr *DocumentReader) Next() (Document, error) {
	var doc Document

	if !dr.inCollection {
		return doc, fmt.Errorf("not in collection")
	}

	if !dr.inDocument {
		return doc, fmt.Errorf("not in document")
	}

	switch se := dr.token.(type) {
	case xml.StartElement:
		if se.Name.Local == "document" {

			se, _ := dr.token.(xml.StartElement)
			dr.decoder.DecodeElement(&doc, &se)

			token, err := dr.decoder.Token()
			if err != nil {
				panic(err)
			}
			dr.token = token

			doc.Map()
			return doc, nil
		}

	case xml.EndElement:
		if se.Name.Local == "collection" {
			dr.inDocument = false
			return doc, fmt.Errorf("eof")
		}
	}
	return doc, nil
}

type WriteDocument struct {
	XMLFile *os.File
}

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

	wr.writeString("<collection>")

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
