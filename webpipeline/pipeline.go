package webpipeline

import "sync"

type PipelineExecutor struct {
	storeChann chan StoreRequest
	wg         sync.WaitGroup
}

func NewPipelineExecutor() PipelineExecutor {
	return PipelineExecutor{
		storeChann: make(chan StoreRequest, 1000),
		wg:         sync.WaitGroup{},
	}
}

func (e *PipelineExecutor) Start() {

}
func (e *PipelineExecutor) Stop() {

}
func (e *PipelineExecutor) Store() {

}
