package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var (
	TRANS = "http://localhost:9091/transmission/rpc"
)

func Err(e error) {
  if e!=nil {
    log.Fatal(e)
  }
}

func GetTransmission() (isTrue bool) {
	page, _ := http.Get(TRANS)
	if page!=nil {
		isTrue = true
	}
	
	return isTrue
}

func AddTorrent(path string, magnet string) {
	client := &http.Client{}
	t := GetPage(TRANS).Header

	data := fmt.Sprintf(`{
		"method": "torrent-add",
		"arguments": {
			"paused": "false",
			"download-dir": "%s",
			"filename": "%s"
		}
	}`, path, magnet)
	
	req, e := http.NewRequest("POST", TRANS, bytes.NewBuffer([]byte(data))); Err(e)
	req.Header = t
	_, er := client.Do(req)
	Err(er)

	defer req.Body.Close()

}

func GetPage(url string) *http.Response {
  page, e := http.Get(url)
  Err(e)
  return page
}

func GetMagnets(url string) (magnets []string)  {
	client := &http.Client{}
	req, e := http.NewRequest("GET", url, nil); Err(e)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; WOW64; x64; rv:105.0esr) Gecko/20010101 Firefox/105.0esr`)
	page, e := client.Do(req); Err(e)

  doc, e := goquery.NewDocumentFromResponse(page); Err(e) 

  doc.Find("a[href^=magnet]").Each(func(_ int, s *goquery.Selection) {
    magnet, m := s.Attr("href")
    if m {
      magnets = append(magnets, magnet)
    }
  })


  defer page.Body.Close()

  return magnets
}

func main()  {
	var wg sync.WaitGroup
	url := flag.String("u", "", "-u [url]")
	wd, _ := os.Getwd()
	path := flag.String("p", wd, "-p [path]")
	flag.Parse()

  magnets := GetMagnets(*url)  
  wg.Add(len(magnets))
  for _, magnet := range magnets {
  	go func(path string, magnet string) {
  		AddTorrent(path, magnet)
  		wg.Done()
  	}(*path, magnet)
  }

  wg.Wait()

}
