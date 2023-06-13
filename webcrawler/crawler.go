package webcrawler

import "github.com/chnyangjie/ohmyspider/webpipeline"

type Crawler interface {
	Start(crawler *CrawlerExecutor, pipeline *webpipeline.PipelineExecutor)
	Stop()
	Crawl(req HTTPRequest)
}

type CrawlCallback func(request HTTPRequest, response []byte)
