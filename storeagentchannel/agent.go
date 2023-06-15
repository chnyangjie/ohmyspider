package storeagentchannel

import "github.com/chnyangjie/ohmyspider/webpipeline"

type ChannelStoreAgent struct {
	outputChannel chan webpipeline.StoreRequest
	agentId       string
}

func NewChannelStoreAgent(agentId string, outputChannel chan webpipeline.StoreRequest) *ChannelStoreAgent {
	return &ChannelStoreAgent{
		outputChannel: outputChannel,
		agentId:       agentId,
	}
}

func (a *ChannelStoreAgent) CanStore(req webpipeline.StoreRequest) bool {
	return true
}

func (a *ChannelStoreAgent) DoStore(req webpipeline.StoreRequest) error {
	a.outputChannel <- req
	return nil
}

func (a *ChannelStoreAgent) AgentId() string {
	return a.agentId
}

func (a *ChannelStoreAgent) UniqFunc() webpipeline.IsUniqFunction {
	return nil
}
func (a *ChannelStoreAgent) Close() error {
	close(a.outputChannel)
	return nil
}
