package main

import (
	"encoding/csv"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	log.Println("Starting")

	req, err := http.NewRequest("GET", "https://www.ta.org.jm/routes-and-fares", nil)
	HandleError(err)
	resp, err := http.DefaultClient.Do(req)
	tokenizer := html.NewTokenizer(resp.Body)
	info_index := 0
	region := ""
	origin := ""
	destination := ""
	fare := ""
	file, err := os.Create("result.csv")
	HandleError(err)
	for {
		tt := tokenizer.Next()
		t := tokenizer.Token()
		err := tokenizer.Err()
		if err == io.EOF {
			break
		}
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			isFareTable := string(t.Data) == "div" && strings.Contains(t.String(), `class="fare-table"`)
			isRegionStart := string(t.Data) == "div" && strings.Contains(t.String(), `class="region"`)
			isTableData := string(t.Data) == "div" && strings.Contains(t.String(), `class="table-row"`)
			isTableInfo := string(t.Data) == "div" && strings.Contains(t.String(), `class="column`)
			if isRegionStart {

				region = strings.Replace(strings.Split(strings.Split(t.String(), `"`)[1], `"`)[0], "_region", "", 1)
				log.Println("Found Region:" + region)
			}
			if isFareTable {
				log.Println("Found Fare Table:" + string(t.String()))
			}
			if isTableData {
				log.Println("Found Table Data:" + string(t.String()))
			}
			if isTableInfo {
				tokenizer.Next()
				switch info_index {
				case 0:
					origin = string(tokenizer.Token().String())
				case 1:
					destination = string(tokenizer.Token().String())
				case 2:
					fare = string(tokenizer.Token().String())

				}
				info_index = info_index + 1
			}
			if info_index == 3 {
				info_index = 0
				WriteInfo(file, region, origin, destination, fare)
			}

		}
	}
	HandleError(file.Close())

}
func WriteInfo(file *os.File, r string, o string, d string, f string) {
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err := writer.Write([]string{r, o, d, f})
	log.Println("Write Info:" + r + " " + o + " " + d + " " + f)
	HandleError(err)
	writer.Flush()

}

func HandleError(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}
}
