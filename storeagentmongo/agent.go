package storeagentmongo

import (
	"context"
	"log"

	"github.com/chnyangjie/ohmyspider/webpipeline"
	"go.mongodb.org/mongo-driver/mongo"
)

type StoreAgentMongo struct {
	agentId     string
	mongoClient *mongo.Client
	isUniqFunc  webpipeline.IsUniqFunction
}

func NewStoreAgentMongo(agentId string, client *mongo.Client, uniqFunc webpipeline.IsUniqFunction) *StoreAgentMongo {
	return &StoreAgentMongo{
		agentId:     agentId,
		mongoClient: client,
		isUniqFunc:  uniqFunc,
	}
}

func (a *StoreAgentMongo) CanStore(req webpipeline.StoreRequest) bool {
	if len(req.StoreParams) != 2 {
		return false
	}
	for _, content := range req.StoreContent {
		_, ok := content.(map[string]interface{})
		if !ok {
			log.Printf("Content Type is not notionapi.Property: %v", content)
			return false
		}
	}
	return true
}

func (a *StoreAgentMongo) DoStore(req webpipeline.StoreRequest) error {
	database := req.StoreParams[0]
	collection := req.StoreParams[1]
	result, err := a.mongoClient.Database(database).Collection(collection).InsertMany(context.Background(), req.StoreContent)
	if err != nil {
		log.Printf("InsertMany error: %v", err)
		return err
	}
	log.Printf("InsertMany result: %v", result)
	return nil
}

func (a *StoreAgentMongo) AgentId() string {
	return a.agentId
}

func (a *StoreAgentMongo) UniqFunc() webpipeline.IsUniqFunction {
	return a.isUniqFunc
}

func (a *StoreAgentMongo) Close() error {
	return a.mongoClient.Disconnect(context.Background())
}
