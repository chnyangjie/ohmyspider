package webpipeline

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/chnyangjie/ohmynotionagent/notionagent"
	"github.com/jomei/notionapi"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

var (
	storeChann = make(chan StoreRequest, 1000)
	wg         = sync.WaitGroup{}
)

func Save(request StoreRequest) {
	wg.Add(1)
	storeChann <- request
}

func Start(startConfig StartConfig) {
	var notionClient *notionapi.Client
	var larkClient *lark.Client
	if startConfig.NotionToken != "" {
		notionClient = notionapi.NewClient(notionapi.Token(startConfig.NotionToken))
	}
	if startConfig.LarkToken != "" && startConfig.LarkAppId != "" {
		larkClient = lark.NewClient(startConfig.LarkAppId, startConfig.LarkToken)
	}
	go func() {
		for {
			select {
			case request := <-storeChann:
				{
					func(wg *sync.WaitGroup) {
						defer wg.Done()
						if request.FilePath != "" {
							StoreToFile(request)
						} else if len(request.LarkContent) != 0 {
							if larkClient != nil {
								StoreToLarkBitable(larkClient, request)
							}
						} else {
							if notionClient != nil {
								StoreToNotionDatabase(notionClient, request)
							}
						}
						if request.SendToChannel && startConfig.OneTimeChannel != nil && request.ContentJson != "" {
							startConfig.OneTimeChannel <- request.ContentJson
						}
					}(&wg)
				}
			}
		}
	}()
}
func StoreToLarkBitable(larkClient *lark.Client, request StoreRequest) {
	log.Printf("Store request: %v", len(request.LarkContent))
	req := larkbitable.NewBatchCreateAppTableRecordReqBuilder().
		AppToken(request.LarkAppToken).
		TableId(request.LarkTableId).
		Body(larkbitable.NewBatchCreateAppTableRecordReqBodyBuilder().
			Records(request.LarkContent).
			Build()).
		Build()
	larkClient.Bitable.AppTableRecord.BatchCreate(context.Background(), req)
}
func StoreToNotionDatabase(client *notionapi.Client, request StoreRequest) {
	log.Printf("Store request: %v", len(request.Content))
	notionagent.CreateNewPageInDatabase(client, request.Database, request.Content)
}
func StoreToFile(request StoreRequest) {
	log.Printf("Store request: %v", len(request.FileContent))
	file, err := os.OpenFile(request.FilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Open file error: %v", err)
		return
	}
	defer file.Close()
	file.Write(request.FileContent)
}

func Stop() {
	wg.Wait()
}
