package webpipeline

type StoreAgent interface {
	CanStore(req StoreRequest) bool
	DoStore(req StoreRequest) error
	AgentId() string
	UniqFunc() IsUniqFunction
	Close() error
}
