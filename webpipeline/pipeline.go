package webpipeline

import (
	"log"
	"sync"
)

type PipelineExecutor struct {
	storeChann chan StoreRequest
	wg         sync.WaitGroup
	quit       chan bool
	agentList  []StoreAgent
}

func NewPipelineExecutor(agentList []StoreAgent) PipelineExecutor {
	return PipelineExecutor{
		storeChann: make(chan StoreRequest, 1000),
		wg:         sync.WaitGroup{},
		quit:       make(chan bool),
		agentList:  agentList,
	}
}

func (e *PipelineExecutor) Start() {
	go e.startConsume()
}
func (e *PipelineExecutor) Stop() {
	e.wg.Wait()
	e.quit <- true
	close(e.storeChann)
}
func (e *PipelineExecutor) Store(req StoreRequest) {
	e.wg.Add(1)
	e.storeChann <- req
}

func (e *PipelineExecutor) startConsume() {
	for {
		select {
		case request := <-e.storeChann:
			{
				func(wg *sync.WaitGroup) {
					defer wg.Done()
					for _, agent := range e.agentList {
						if agent == nil {
							continue
						}
						if request.AgentId != nil && *request.AgentId != agent.AgentId() {
							continue
						}
						if agent.CanStore(request) {
							if isUniq(agent, request) {
								agent.DoStore(request)
							}
						}
					}
				}(&e.wg)
			}
		case <-e.quit:
			{
				log.Printf("Closing pipeline")
				for _, agent := range e.agentList {
					if agent == nil {
						continue
					}
					err := agent.Close()
					if err != nil {
						log.Printf("Closing agent %s failed: %+v", agent.AgentId(), err)
					}
				}
				return
			}
		}
	}
}

func isUniq(agent StoreAgent, req StoreRequest) bool {
	if req.Source == nil || req.UniqId == nil {
		return true
	}
	if req.IsUniqFunction != nil {
		if !req.IsUniqFunction(*req.Source, *req.UniqId) {
			return false
		}
	}
	if agent.UniqFunc() != nil {
		if !agent.UniqFunc()(*req.Source, *req.UniqId) {
			return false
		}
	}
	return true
}
