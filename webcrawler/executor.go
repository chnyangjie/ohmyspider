package webcrawler

import (
	"log"
	"sync"
)

type CrawlerExecutor struct {
	crawlerChannel chan HTTPRequest
	wg             sync.WaitGroup
}

func NewCrawlerExecutor() CrawlerExecutor {
	return CrawlerExecutor{
		crawlerChannel: make(chan HTTPRequest, 10000),
		wg:             sync.WaitGroup{},
	}
}
func (e *CrawlerExecutor) Start() {
	go e.startCrawl()
}

func (e *CrawlerExecutor) startCrawl() {
	for {
		select {
		case req := <-e.crawlerChannel:
			{
				func() {
					defer e.wg.Done()
					log.Printf("Crawling %+v", req)
					resp, err := execute(req)
					if err != nil || resp == nil {
						log.Printf("Crawling %+v failed", req)
						return
					}
					if req.Callback != nil {
						req.Callback(req, resp)
					} else {
						log.Printf("Crawling %+v without callback", req)
					}
				}()
			}
		default:
			{
				log.Printf("No request to process")
				return
			}
		}
	}
}
func (e *CrawlerExecutor) Crawl(req HTTPRequest) {
	e.wg.Add(1)
	e.crawlerChannel <- req
}

func (e *CrawlerExecutor) Stop() {
	e.wg.Wait()
	close(e.crawlerChannel)
}
