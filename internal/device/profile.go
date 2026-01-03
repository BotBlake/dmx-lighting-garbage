package device

import "encoding/xml"

type DDFDevice struct {
	XMLName   xml.Name  `xml:"device"`
	Functions Functions `xml:"functions"`
}

type Functions struct {
	RGB     *RGB      `xml:"rgb"`
	Dimmer  *Dimmer   `xml:"dimmer"`
	RawList []RawChan `xml:"raw"`
}

type RGB struct {
	Red   Channel `xml:"red"`
	Green Channel `xml:"green"`
	Blue  Channel `xml:"blue"`
}

type Dimmer struct {
	DMXChannel int `xml:"dmxchannel,attr"`
}

type RawChan struct {
	Name       string `xml:"name,attr"`
	DMXChannel int    `xml:"dmxchannel,attr"`
}

type Channel struct {
	DMXChannel int `xml:"dmxchannel,attr"`
}
