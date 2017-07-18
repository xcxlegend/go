package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

type BodyParseType int

const (
	BODYPARSETYPE_HTML BodyParseType = iota
	BODYPARSETYPE_JSON
)

// 爬虫任务

type CrawlerScheme struct {
	Id     int // 标识
	option *CrawlerSchemeOption
	wg     *sync.WaitGroup
}

type CrawlerSchemeOption struct {
	ListUri          string // http://xx.xx.com/abc/def-%d
	MaxListPager     int
	ListParseType    BodyParseType
	ListParsePattern string // preg
	PageParseType    BodyParseType
	PageParsePattern string
	MaxGONumber      int // 同时进行的数量
}

func NewCrawlerScheme(id int, option *CrawlerSchemeOption) *CrawlerScheme {
	c := &CrawlerScheme{
		Id:     id,
		option: option,
		wg:     new(sync.WaitGroup),
	}
	return c
}

func (c *CrawlerScheme) run(wg *sync.WaitGroup) {
	defer wg.Done()
	if c.MaxListPager > 0 {
		for i := 1; i <= c.MaxListPager; i++ {
			var uri = fmt.Sprintf(c.option.ListUri, i)
			var body = httpGet(uri)
			if body == nil {
				break
			}
			var urls = searchByPatterns(*body, c.option.ListParsePattern, 1)
		}
	}
}

func (c *CrawlerScheme) handler() {
	for {

	}
}

func searchByPatterns(text string, pattern string, index int) []string {
	var req = regexp.MustCompile(pattern)
	var res = req.FindAllStringSubmatch(text, -1)
	var search = []string{}
	if len(res) > 0 {
		for _, l := range res {
			if len(l) > index-1 {
				search = append(search, l[index])
			}
		}
	}
	return search
}

func httpGet(uri string) *string {
	resp, err := http.Get(uri)
	if err != nil {
		return nil
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	return *string(body)
}

// listPattern \d
//
