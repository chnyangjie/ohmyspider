package crawler

import (
	"github.com/chnyangjie/ohmyspider/webcrawler"
	"github.com/chnyangjie/ohmyspider/webpipeline"
)

type Crawler interface {
	Start(crawler *webcrawler.CrawlerExecutor, pipeline *webpipeline.PipelineExecutor)
	Stop()
	Crawl(req webcrawler.HTTPRequest)
}
