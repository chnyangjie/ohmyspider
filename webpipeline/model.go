package webpipeline

import (
	"github.com/jomei/notionapi"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

type StartConfig struct {
	NotionToken    string
	LarkAppId      string
	LarkToken      string
	OneTimeChannel chan StoreRequest
}

type StoreType string

const (
	StoreTypeNotionDatabase StoreType = "NOTION_DATABASE"
	StoreTypeLarkBitable    StoreType = "LARK_BITABLE"
	StoreTypeFile           StoreType = "FILE"
)

type StoreRequest struct {
	UniqId         string
	Source         string
	IsUniqFunction IsUniqFunction

	NotionDatabase notionapi.DatabaseID
	NotionContent  map[string]notionapi.Property

	LarkContent  []*larkbitable.AppTableRecord
	LarkAppToken string
	LarkTableId  string

	FilePath    string
	FileContent []byte

	SendToChannel bool
	ExtraContent  interface{}
}
type StoreFunction func(request StoreRequest) error
type IsUniqFunction func(source, uniqId string) bool
