package storeagentnotion

import (
	"log"

	"github.com/chnyangjie/ohmynotionagent/notionagent"
	"github.com/chnyangjie/ohmyspider/webpipeline"
	"github.com/jomei/notionapi"
)

type StoreAgentNotion struct {
	notionClient *notionapi.Client
	agentId      string
	isUniqFunc   webpipeline.IsUniqFunction
}

func NewStoreAgentNotion(agentId, notionToken string, uniqFunc webpipeline.IsUniqFunction) *StoreAgentNotion {
	return &StoreAgentNotion{
		notionClient: notionapi.NewClient(notionapi.Token(notionToken)),
		agentId:      agentId,
		isUniqFunc:   uniqFunc,
	}
}

func (a *StoreAgentNotion) CanStore(req webpipeline.StoreRequest) bool {
	if len(req.StoreParams) != 1 {
		return false
	}
	for _, content := range req.StoreContent {
		_, ok := content.(map[string]notionapi.Property)
		if !ok {
			log.Printf("Content Type is not notionapi.Property: %v", content)
			return false
		}
	}
	return true
}

func (a *StoreAgentNotion) DoStore(req webpipeline.StoreRequest) error {
	databaseId := req.StoreParams[0]
	for _, content := range req.StoreContent {
		data := content.(map[string]notionapi.Property)
		notionagent.CreateNewPageInDatabase(a.notionClient, notionapi.DatabaseID(databaseId), data)
	}
	return nil
}

func (a *StoreAgentNotion) AgentId() string {
	return a.agentId
}

func (a *StoreAgentNotion) UniqFunc() webpipeline.IsUniqFunction {
	return a.isUniqFunc
}

func (a *StoreAgentNotion) Close() error {
	return nil
}
