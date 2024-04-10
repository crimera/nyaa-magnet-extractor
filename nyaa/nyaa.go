package nyaa

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const (
	TRANS    = "http://localhost:9091/transmission/rpc"
	BASE_URL = "https://nyaa.si/"
)

func Err(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func GetTransmission() (isTrue bool) {
	page, _ := http.Get(TRANS)
	if page != nil {
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

	req, e := http.NewRequest("POST", TRANS, bytes.NewBuffer([]byte(data)))
	Err(e)

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

func Get(url string) (*http.Response, error) {
	client := &http.Client{}
	req, e := http.NewRequest("GET", url, nil)
	Err(e)

	req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; WOW64; x64; rv:105.0esr) Gecko/20010101 Firefox/105.0esr`)

	return client.Do(req)
}

func GetMagnets(url string) (magnets []string) {
	client := &http.Client{}
	req, e := http.NewRequest("GET", url, nil)
	Err(e)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; WOW64; x64; rv:105.0esr) Gecko/20010101 Firefox/105.0esr`)
	page, e := client.Do(req)
	Err(e)

	doc, e := goquery.NewDocumentFromResponse(page)
	Err(e)

	doc.Find("a[href^=magnet]").Each(func(_ int, s *goquery.Selection) {
		magnet, m := s.Attr("href")
		if m {
			magnets = append(magnets, magnet)
		}
	})

	defer page.Body.Close()

	return magnets
}

func Query(q string, sort string, order string) (items []Item) {
	url, _ := url.Parse(BASE_URL)
	query := url.Query()
	query.Set("q", q)
	query.Set("s", sort)
	query.Set("o", order)
	url.RawQuery = query.Encode()

	page, e := Get(url.String())
	Err(e)

	doc, e := goquery.NewDocumentFromResponse(page)
	Err(e)

	doc.Find("tbody").First().Find("tr").Each(func(i int, s *goquery.Selection) {
		title := ""
		category := ""
		magnet := ""
		size := ""

		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				category, _ = s.Find("a[href][title]").Last().Attr("title")
			case 1:
				title, _ = s.Find("a[href][title]").Last().Attr("title")
			case 2:
				magnet, _ = s.Find("a[href]").Last().Attr("href")
			case 3:
				size = s.Text()
			}
		})

		items = append(items, Item{category, title, size, magnet})
	})

	return items
}

func AddToClient() {
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

type Item struct {
	Category string
	Name     string
	Size     string
	Magnet   string
}
