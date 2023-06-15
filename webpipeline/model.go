package webpipeline

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
	UniqId         *string
	Source         *string
	IsUniqFunction IsUniqFunction
	AgentId        *string

	StoreParams  []string
	StoreContent []interface{}
}
type StoreFunction func(request StoreRequest) error
type IsUniqFunction func(source, uniqId string) bool
