package gdelt

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"example.com/m/v2/global"
)

func GetRawData(date time.Time) ([]Entry, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	quit := make(chan bool)
	collection := []Entry{}
	wg := sync.WaitGroup{}
	for i := 0; i <= 24*4-1; i++ {
		wg.Add(1)
		go func() {
			download(fmt.Sprintf("http://data.gdeltproject.org/gdeltv2/%v.gkg.csv.zip", start.Format("20060102150405")), quit, &collection)
			wg.Done()
		}()
		start = start.Add(time.Minute * 15)
	}
	wg.Wait()
	close(quit)
	if len(quit) != 0 {
		return nil, errors.New("download failed")
	}

	return collection, nil
}

func download(url string, quit chan bool, collection *[]Entry) {
	for r := global.DownloadRetry; r > 0; r-- {
		global.Log.Println("downloading from ", url)
		resp, err := http.Get(url)
		if err != nil {
			global.Log.Println("download from ", url, " failed, ", err.Error(), " remaining reties: ", r)
			continue
		}
		global.Log.Println("download from ", url, " succeed")
		defer resp.Body.Close()

		buff := bytes.NewBuffer([]byte{})
		size, err := io.Copy(buff, resp.Body)
		if err != nil {
			global.Log.Println("copy response ", url, " failed, ", err.Error(), " remaining reties: ", r)
			continue
		}
		reader := bytes.NewReader(buff.Bytes())
		csvFile, err := zip.NewReader(reader, size)
		if err != nil {
			global.Log.Println("unzip ", url, " failed, ", err.Error(), " remaining reties: ", r)
			continue
		}

		for _, f := range csvFile.File {
			content, err := f.Open()
			if err != nil {
				global.Log.Println("open ", f.Name, " failed, ", err.Error(), " remaining reties: ", r)
				continue
			}
			scanner := bufio.NewScanner(content)
			var co = []Entry{}
			for scanner.Scan() {
				line := strings.Split(scanner.Text(), "\t")
				DATE := line[1]
				DocumentIdentifier := line[4]
				V2Tone := line[15]
				Themes := line[7]
				Organizations := line[13]
				if len(Organizations) != 0 && len(Themes) != 0 {
					co = append(co, Entry{DATE, DocumentIdentifier, V2Tone, Themes, Organizations})
				}
			}
			global.Log.Println("process ", f.Name, " succeed")
			*collection = append(*collection, co...)
		}
		return
	}
	quit <- true
	return

}
