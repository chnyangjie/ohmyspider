package storeagentlark

import (
	"context"
	"log"

	"github.com/chnyangjie/ohmyspider/webpipeline"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

type StoreAgentLark struct {
	agentId    string
	larkClient *lark.Client
	isUniqFunc webpipeline.IsUniqFunction
}

func NewStoreAgentLark(agentId string, larkAppid, larkSecretToken string, uniqFunc webpipeline.IsUniqFunction) *StoreAgentLark {
	return &StoreAgentLark{
		agentId:    agentId,
		larkClient: lark.NewClient(larkAppid, larkSecretToken),
		isUniqFunc: uniqFunc,
	}
}

func (a *StoreAgentLark) AgentId() string {
	return a.agentId
}

func (a *StoreAgentLark) Close() error {
	return nil
}

func (a *StoreAgentLark) CanStore(req webpipeline.StoreRequest) bool {
	if len(req.StoreParams) != 2 {
		return false
	}
	for _, content := range req.StoreContent {
		_, ok := content.(*larkbitable.AppTableRecord)
		if !ok {
			log.Printf("content is not *larkbitable.AppTableRecord, content: %v", content)
			return false
		}
	}
	return true
}

func (a *StoreAgentLark) DoStore(req webpipeline.StoreRequest) error {
	contentList := []*larkbitable.AppTableRecord{}
	for _, content := range req.StoreContent {
		data := content.(*larkbitable.AppTableRecord)
		contentList = append(contentList, data)
	}
	appToken := req.StoreParams[0]
	tableId := req.StoreParams[1]
	bitableReq := larkbitable.NewBatchCreateAppTableRecordReqBuilder().
		AppToken(appToken).
		TableId(tableId).
		Body(larkbitable.NewBatchCreateAppTableRecordReqBodyBuilder().
			Records(contentList).
			Build()).
		Build()
	resp, err := a.larkClient.Bitable.AppTableRecord.BatchCreate(context.Background(), bitableReq)
	if err != nil {
		log.Printf("larkClient.Bitable.AppTableRecord.BatchCreate error: %v", err)
	}
	log.Printf("larkClient.Bitable.AppTableRecord.BatchCreate resp: %v", resp)
	return err
}

func (a *StoreAgentLark) UniqFunc() webpipeline.IsUniqFunction {
	return a.isUniqFunc
}
