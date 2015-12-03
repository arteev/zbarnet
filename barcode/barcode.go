package barcode

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
)

//Errors for barcode
var (
	ErrDataBufferEmpty = errors.New("The data buffer is empty")
	ErrFormatError     = errors.New("Unknown data format")
	ErrSkipRootTag     = errors.New("Skip root tag")
)

type result struct {
	XMLName xml.Name `xml:"index"`
	Barcode BarCode  `xml:"symbol"`
}

//A BarCode incoming barcode from buffer
type BarCode struct {
	Type    string `xml:"type,attr" json:"type"`
	Quality int    `xml:"quality,attr" json:"quality"`
	Data    []byte `xml:"data" json:"data"`
}

//A Mode output zbarcam
type Mode int

//Pre defined the modes output zbarcam
const (
	ModeRaw    = Mode(0)
	ModeXML    = Mode(1)
	ModeNative = Mode(2)
)

//MustBarCode decode and create *Barcode from buffer with xml format or native format of zbarcam
func MustBarCode(data []byte, mode Mode) (*BarCode, error) {
	v := result{}
	if len(data) == 0 {
		return nil, ErrDataBufferEmpty
	}
	sData := strings.TrimSpace(string(data))
	if mode == ModeXML {
		//xml format
		if strings.HasPrefix(sData, "<barcodes") || strings.HasPrefix(sData, "</source></barcodes>") {
			return nil, ErrSkipRootTag
		}
		err := xml.Unmarshal(data, &v)
		if err != nil {
			return nil, err
		}
	} else if mode == ModeNative {
		//native format
		idx := bytes.Index(data, []byte{':'})
		if idx == -1 {
			return nil, ErrFormatError
		}
		v.Barcode.Type = string(data[0:idx])
		v.Barcode.Data = data[idx:]
		v.Barcode.Quality = 1
	} else {
		//Raw
		v.Barcode.Type = "raw"
		v.Barcode.Quality = 1
		v.Barcode.Data = data

	}
	return &v.Barcode, nil
}

//ToJSON returns BarCode as formatted JSON
//or error if serialization failed
func (b *BarCode) ToJSON() (string, error) {
	out, err := json.MarshalIndent(&b, "", "")
	if err != nil {
		return "", err
	}
	return string(out), nil
}
