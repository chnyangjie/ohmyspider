package webpipeline

import (
	"context"
	"log"
	"os"

	"github.com/chnyangjie/ohmynotionagent/notionagent"
	"github.com/jomei/notionapi"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

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
	log.Printf("Store request: %v", len(request.NotionContent))
	notionagent.CreateNewPageInDatabase(client, request.NotionDatabase, request.NotionContent)
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
