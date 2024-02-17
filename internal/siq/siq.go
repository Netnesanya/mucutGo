package siq

import (
	"encoding/xml"
)

type Package struct {
	XMLName xml.Name `xml:"package"`
	Name    string   `xml:"name,attr"`
	Rounds  []Round  `xml:"rounds>round"`
}

type Round struct {
	Name   string  `xml:"name,attr"`
	Themes []Theme `xml:"themes>theme"`
}

type Theme struct {
	Name      string     `xml:"name,attr"`
	Questions []Question `xml:"questions>question"`
}

type Question struct {
	Price    int      `xml:"price,attr"`
	Scenario []Atom   `xml:"scenario>atom"`
	Right    []string `xml:"right>answer"`
}

type Atom struct {
	Type string `xml:"type,attr"`
	Text string `xml:",chardata"`
}
