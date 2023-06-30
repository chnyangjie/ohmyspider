package storeagentfile

import (
	"log"
	"os"
	"path/filepath"

	"github.com/chnyangjie/ohmyspider/webpipeline"
)

type StoreAgentFile struct {
	agentId  string
	filePath string
	file     *os.File
}

func NewStoreAgentFile(agentId string, filePath string) *StoreAgentFile {
	var f *os.File
	if filePath != "" {
		f, _ = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	}
	return &StoreAgentFile{
		agentId:  agentId,
		filePath: filePath,
		file:     f,
	}
}

func (a *StoreAgentFile) CanStore(req webpipeline.StoreRequest) bool {
	if len(req.StoreParams) != len(req.StoreContent) && a.file == nil {
		log.Printf("StoreAgentFile: StoreParams and StoreContent length not equal, and file is nil")
		return false
	}
	for _, content := range req.StoreContent {
		_, ok := content.([]byte)
		if !ok {
			log.Printf("StoreAgentFile: StoreContent is not []byte")
			return false
		}
	}
	return true
}
func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
func (a *StoreAgentFile) DoStore(req webpipeline.StoreRequest) error {
	if len(req.StoreParams) == len(req.StoreContent) {
		for i := 0; i < len(req.StoreParams); i++ {
			data := req.StoreContent[i].([]byte)
			file, err := create(req.StoreParams[i])
			if err != nil {
				log.Printf("StoreAgentFile: create file error: %v", err)
				return err
			}
			defer file.Close()
			file.Write(data)
		}
		return nil
	} else if a.file != nil {
		for _, content := range req.StoreContent {
			data := content.(string)
			a.file.Write([]byte(data))
		}
	}
	return nil
}
func (a *StoreAgentFile) AgentId() string {
	return a.agentId
}
func (a *StoreAgentFile) Close() error {
	if a.file != nil {
		a.file.Close()
	}
	return nil
}
func (a *StoreAgentFile) UniqFunc() webpipeline.IsUniqFunction {
	return nil
}
