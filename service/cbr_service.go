package service

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"os"
)

type ValCurs struct {
	Date   string `xml:",attr"`
	Valute []struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Nominal  string `xml:"Nominal"`
		Name     string `xml:"Name"`
		Value    string `xml:"Value"`
	} `xml:"Valute"`
}

func GetCbrCurrency() (ValCurs, error) {
	resp, err := http.Get(os.Getenv("FiatCurrencyUrl"))
	valCurs := ValCurs{}
	if err != nil {
		log.Print(err)
		return ValCurs{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Print(err)
		return valCurs, err
	}
	decode := xml.NewDecoder(resp.Body)
	decode.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}
	err = decode.Decode(&valCurs)
	if err != nil {
		log.Println(err)
	}
	return valCurs, nil
}
