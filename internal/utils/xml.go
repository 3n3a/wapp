package utils

import "encoding/xml"

type XMLKeyValue struct {
	XMLName xml.Name    `xml:"KeyValue"`
	Name    string      `xml:"key,attr"`
	Value   interface{} `xml:",chardata"`
}

// Create a struct to represent the KeyValues element
type XMLKeyValues struct {
	XMLName   xml.Name      `xml:"KeyValues"`
	KeyValues []XMLKeyValue `xml:"KeyValue"`
}

// Transforms a utils.Map to XML Key Values
func (m *Map) ToXML() (XMLKeyValues, error) {
	var kvs XMLKeyValues
	for key, val := range *m {
		kvs.KeyValues = append(
			kvs.KeyValues,
			XMLKeyValue{
				Name:  key,
				Value: val,
			},
		)
	}
	return kvs, nil
}