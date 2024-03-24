package utils

import "encoding/xml"

type XMLItemInner struct {
	XMLName xml.Name
	Value   interface{} `xml:",chardata"`
}

type XMLItem struct {
	XMLName xml.Name `xml:"Item"`
	Inner   []XMLItemInner
}

// Create a struct to represent the KeyValues element
type XMLList struct {
	XMLName xml.Name  `xml:"List"`
	Items   []XMLItem `xml:"Item"`
}

// Transforms a utils.Map to XML
func (m *Map) ToXML() (XMLItem, error) {
	inners := []XMLItemInner{}

	for key, val := range *m {
		inners = append(
			inners,
			XMLItemInner{
				XMLName: xml.Name{Local: key},
				Value:   val,
			},
		)
	}

	var item = XMLItem{
		Inner: inners,
	}
	return item, nil
}
