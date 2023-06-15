package storeagentmysql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chnyangjie/ohmyspider/webpipeline"
)

type StoreAgentMysql struct {
	agentId    string
	isUniqFunc webpipeline.IsUniqFunction
	db         *sql.DB
}

func NewStoreAgentMysql(agentId, username, password, dbName, lunk string, uniqFunc webpipeline.IsUniqFunction) *StoreAgentMysql {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", username, password, lunk, dbName))
	if err != nil {
		log.Printf("Open mysql error: %v", err)
		return nil
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(1)
	return &StoreAgentMysql{
		agentId:    agentId,
		isUniqFunc: uniqFunc,
		db:         db,
	}
}

func (a *StoreAgentMysql) AgentId() string {
	return a.agentId
}
func (a *StoreAgentMysql) Close() error {
	return a.db.Close()
}

func (a *StoreAgentMysql) CanStore(req webpipeline.StoreRequest) bool {
	for _, content := range req.StoreContent {
		_, ok := content.(string)
		if !ok {
			return false
		}
	}
	return true
}
func (a *StoreAgentMysql) DoStore(req webpipeline.StoreRequest) error {
	for _, content := range req.StoreContent {
		data := content.(string)
		_, err := a.db.Exec(data)
		if err != nil {
			log.Printf("Exec error: %v", err)
		}
	}
	return nil
}

func (a *StoreAgentMysql) UniqFunc() webpipeline.IsUniqFunction {
	return a.isUniqFunc
}
