package xmls

import "encoding/xml"

// CDataText struct for xml marshal with cdata
type CDataText struct {
	Text string `xml:",cdata"`
}

// CharData string like "<![CDATA[" + v + "]]>"
type CharData string

// MarshalXML marshal string with cdata
func (c CharData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(CDataText{string(c)}, start)
}
