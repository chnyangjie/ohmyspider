package webpipeline

import (
	"log"
	"sync"

	"github.com/jomei/notionapi"
	lark "github.com/larksuite/oapi-sdk-go/v3"
)

type PipelineExecutor struct {
	storeChann   chan StoreRequest
	wg           sync.WaitGroup
	config       StartConfig
	notionClient *notionapi.Client
	larkClient   *lark.Client
}

func NewPipelineExecutor(config StartConfig) PipelineExecutor {
	return PipelineExecutor{
		storeChann: make(chan StoreRequest, 1000),
		wg:         sync.WaitGroup{},
		config:     config,
	}
}

func (e *PipelineExecutor) Start() {
	if e.config.NotionToken != "" {
		e.notionClient = notionapi.NewClient(notionapi.Token(e.config.NotionToken))
	}
	if e.config.LarkToken != "" && e.config.LarkAppId != "" {
		e.larkClient = lark.NewClient(e.config.LarkAppId, e.config.LarkToken)
	}
	go e.startConsume()
}
func (e *PipelineExecutor) Stop() {
	e.wg.Wait()
	close(e.storeChann)
}
func (e *PipelineExecutor) Store(req StoreRequest) {
	e.wg.Add(1)
	e.storeChann <- req
}

func (e *PipelineExecutor) startConsume() {
	for {
		select {
		case request := <-e.storeChann:
			{
				func(wg *sync.WaitGroup) {
					defer wg.Done()
					if request.IsUniqFunction != nil {
						if request.UniqId == "" || request.Source == "" {
							return
						}
						if !request.IsUniqFunction(request.Source, request.UniqId) {
							return
						}
					}
					if request.FilePath != "" {
						StoreToFile(request)
					} else if len(request.LarkContent) != 0 {
						if e.larkClient != nil {
							StoreToLarkBitable(e.larkClient, request)
						}
					} else {
						if e.notionClient != nil {
							StoreToNotionDatabase(e.notionClient, request)
						}
					}
					if request.SendToChannel && e.config.OneTimeChannel != nil {
						e.config.OneTimeChannel <- request
					}
				}(&e.wg)
			}
		default:
			{
				log.Printf("No request to process")
				return
			}
		}
	}
}
