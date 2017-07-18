package crawler

import (
	"github.com/astaxie/beego"
	"sync"
)

// 爬虫

type Crawler struct {
	schemes map[int]*CrawlerScheme
	wg      *sync.WaitGroup
}

func NewCrawler() *Crawler {
	c := new(Crawler)
	c.schemes = make(map[int]*CrawlerScheme)
	c.wg = new(sync.WaitGroup)
	return c
}

func (c *Crawler) AddScheme(s *CrawlerScheme) {
	if _, ok := c.schemes[s.Id]; ok {
		beego.Warn(s.Id, " is exist.")
		return
	}
	c.schemes[s.Id] = s
}

func (c *Crawler) Run() {
	for _, sche := range c.schemes {
		c.wg.Add(1)
		go sche.run(c.wg)
	}
	c.wg.Wait()
}
