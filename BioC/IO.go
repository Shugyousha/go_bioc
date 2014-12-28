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
		fmt.Fprintf(os.Stderr, "Error when opening file: %q\n", err)
	}

	breader := bufio.NewReader(file)
	decoder := xml.NewDecoder(breader)

	err = decoder.Decode(&col)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when decoding file: %q\n", err)
	}

	return col
}

func WriteCollection(col Collection, filename string) error {
	perm := os.FileMode(0664)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	bwriter := bufio.NewWriter(f)
	_, err = bwriter.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	_, err = bwriter.Write([]byte("<!DOCTYPE collection SYSTEM \"BioC.dtd\">\n"))
	if err != nil {
		return err
	}

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
				if !dr.inCollection {
					break
				}
				dr.token, err = dr.decoder.Token()
				if err != nil {
					return col, fmt.Errorf("Error when decoding source token: %s", err)
				}
				col.Source = string(dr.token.(xml.CharData))
				err = dr.decoder.Skip()

			case "date":
				if !dr.inCollection {
					break
				}
				dr.token, err = dr.decoder.Token()
				if err != nil {
					return col, fmt.Errorf("Error when decoding date token: %s", err)
				}
				col.Date = string(dr.token.(xml.CharData))
				err = dr.decoder.Skip()

			case "key":
				if !dr.inCollection {
					break
				}
				dr.token, err = dr.decoder.Token()
				if err != nil {
					return col, fmt.Errorf("Error when decoding key token: %s", err)
				}
				col.Key = string(dr.token.(xml.CharData))
				err = dr.decoder.Skip()

			case "infon":
				if !dr.inCollection {
					break
				}
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
				col.Infons[key] = value

			case "document":
				if dr.inCollection {
					dr.inDocument = true
					return col, nil
				}
			}
		}
	}
	return col, nil
}

func (dr *DocumentReader) Next() (*Document, error) {
	var doc Document

	if !dr.inCollection {
		return nil, fmt.Errorf("not in collection")
	}

	if !dr.inDocument {
		return nil, fmt.Errorf("not in document")
	}

	for {
		switch se := dr.token.(type) {
		case xml.StartElement:
			if se.Name.Local == "document" {

				se, _ := dr.token.(xml.StartElement)
				err := dr.decoder.DecodeElement(&doc, &se)
				if err != nil {
					return nil, fmt.Errorf("Error when decoding element: %s", err)
				}

				token, err := dr.decoder.Token()
				if err != nil {
					return nil, fmt.Errorf("Error when getting next token: %s", err)
				}
				dr.token = token

				return &doc, nil
			}

		case xml.EndElement:
			if se.Name.Local == "collection" {
				dr.inDocument = false
				return &doc, fmt.Errorf("EOF")
			}

		default:
			token, err := dr.decoder.Token()
			if err != nil {
				return nil, fmt.Errorf("Error when getting next token: %s", err)
			}
			dr.token = token
		}
	}
	return nil, fmt.Errorf("End of function case")
}

type DocumentWriter struct {
	Writer io.Writer
}

func (dw *DocumentWriter) Start(writer io.Writer, col Collection) error {
	var err error
	dw.Writer = writer

	_, err = dw.Writer.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	_, err = dw.Writer.Write([]byte("<!DOCTYPE collection SYSTEM \"BioC.dtd\">"))
	if err != nil {
		return err
	}

	dw.Writer.Write([]byte(fmt.Sprintf("<collection><source>%s</source><date>%s</date><key>%s</key>", col.Source, col.Date, col.Key)))

	for key, value := range col.Infons {
		dw.Writer.Write([]byte(fmt.Sprintf("<infon key=\"%s\">%s</infon>", key, value)))
	}

	return err
}

func (dw *DocumentWriter) Next(doc Document) error {
	data, err := xml.Marshal(doc)
	if err != nil {
		return err
	}
	dw.Writer.Write(data)
	return nil
}

func (dw *DocumentWriter) Close() {
	dw.Writer.Write([]byte("</collection>"))
}
