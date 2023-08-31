package storeagentpostgresql

import (
	"database/sql"
	"log"

	"github.com/chnyangjie/ohmyspider/webpipeline"
	_ "github.com/lib/pq"
)

type StoreAgentPostgresql struct {
	agentId    string
	isUniqFunc webpipeline.IsUniqFunction
	db         *sql.DB
}

func NewStoreAgentPostgresql(agentId, link string, uniqFunc webpipeline.IsUniqFunction) *StoreAgentPostgresql {
	db, err := sql.Open("postgres", link)
	if err != nil {
		log.Printf("Open postgresql error: %v", err)
		return nil
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(1)
	return &StoreAgentPostgresql{
		agentId:    agentId,
		isUniqFunc: uniqFunc,
		db:         db,
	}
}

func (a *StoreAgentPostgresql) AgentId() string {
	return a.agentId
}
func (a *StoreAgentPostgresql) Close() error {
	return a.db.Close()
}

func (a *StoreAgentPostgresql) CanStore(req webpipeline.StoreRequest) bool {
	for _, content := range req.StoreContent {
		_, ok := content.(string)
		if !ok {
			return false
		}
	}
	return true
}
func (a *StoreAgentPostgresql) DoStore(req webpipeline.StoreRequest) error {
	for _, content := range req.StoreContent {
		data := content.(string)
		_, err := a.db.Exec(data)
		if err != nil {
			log.Printf("Exec error: %v", err)
		}
	}
	return nil
}

func (a *StoreAgentPostgresql) UniqFunc() webpipeline.IsUniqFunction {
	return a.isUniqFunc
}
