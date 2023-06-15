package storeagentfile

import (
	"log"
	"os"

	"github.com/chnyangjie/ohmyspider/webpipeline"
)

type StoreAgentFile struct {
	agentId  string
	filePath string
	file     *os.File
}

func NewStoreAgentFile(agentId string, filePath string) *StoreAgentFile {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Open file error: %v", err)
		return nil
	}
	return &StoreAgentFile{
		agentId:  agentId,
		filePath: filePath,
		file:     file,
	}
}

func (a *StoreAgentFile) CanStore(req webpipeline.StoreRequest) bool {
	if len(req.StoreParams) != 1 {
		return false
	}
	for _, content := range req.StoreContent {
		_, ok := content.(string)
		if !ok {
			return false
		}
	}
	return true
}
func (a *StoreAgentFile) DoStore(req webpipeline.StoreRequest) error {
	for _, content := range req.StoreContent {
		data := content.(string)
		a.file.Write([]byte(data))
	}
	return nil
}
func (a *StoreAgentFile) AgentId() string {
	return a.agentId
}
func (a *StoreAgentFile) Close() error {
	a.file.Close()
	return nil
}
func (a *StoreAgentFile) UniqFunc() webpipeline.IsUniqFunction {
	return nil
}
