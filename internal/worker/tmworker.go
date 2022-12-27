package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/clients"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/config"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
)

type (
	TMWorker struct {
		logger   commons.TedgeLogger
		cli      *clients.TyModelReportClient
		dataChan chan *proto.ThingModelMsg
		doneChan chan struct{}
	}

	TMWorkerPool struct {
		mutex        sync.Mutex
		workers      []*TMWorker
		workersCount uint32
		lastWorker   uint32
		ackMap       sync.Map
	}
)

//////////////////////////////////////////////////////////////////////////////////////////////////////
func NewTMWorkerPool(cliCfg config.ClientInfo, logger commons.TedgeLogger) (*TMWorkerPool, error) {
	var maxWorker uint32 = 16
	var chanCap int32 = 64

	workers := make([]*TMWorker, 0, maxWorker)
	for i := 0; i < int(maxWorker); i++ {
		worker, err := NewTMWorker(cliCfg, chanCap, logger)
		if err != nil {
			return nil, err
		}

		workers = append(workers, worker)
		go worker.run()
	}

	return &TMWorkerPool{
		workersCount: maxWorker,
		lastWorker:   0,
		workers:      workers,
	}, nil
}

func (wp *TMWorkerPool) Stop() {
	wp.mutex.Lock()
	for _, worker := range wp.workers {
		worker.stop()
	}
	wp.mutex.Unlock()
}

func (wp *TMWorkerPool) PutData(data *proto.ThingModelMsg) error {
	wp.mutex.Lock()
	index := wp.lastWorker % wp.workersCount
	wp.lastWorker++
	wp.mutex.Unlock()
	return wp.workers[index].putData(data)
}

func (wp *TMWorkerPool) PutDataWithMsgId(id string, data *proto.ThingModelMsg) (interface{}, error) {
	wp.mutex.Lock()
	index := wp.lastWorker % wp.workersCount
	wp.lastWorker++
	wp.mutex.Unlock()

	ch := wp.genAckChan(id)
	if err := wp.workers[index].putData(data); err != nil {
		ch.TryCloseChan()
		return nil, err
	}

	select {
	case <-time.After(5 * time.Second):
		ch.TryCloseChan()
		return nil, errors.New("wait response timeout")
	case resp := <-ch.DataChan:
		return resp, nil
	}
}

func (wp *TMWorkerPool) StoreMsgId(id string, ch *common.MsgAckChan) {
	wp.ackMap.Store(id, ch)
}

func (wp *TMWorkerPool) DeleteMsgId(id string) {
	wp.ackMap.Delete(id)
}

func (wp *TMWorkerPool) LoadMsgChan(id string) (interface{}, bool) {
	return wp.ackMap.Load(id)
}

func (wp *TMWorkerPool) genAckChan(id string) *common.MsgAckChan {
	ack := &common.MsgAckChan{
		Id:       id,
		DataChan: make(chan interface{}, 1),
	}
	wp.ackMap.Store(id, ack)
	return ack
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
func NewTMWorker(cliCfg config.ClientInfo, chanCap int32, logger commons.TedgeLogger) (*TMWorker, error) {
	cli, err := clients.NewTyModelReportClient(cliCfg)
	if err != nil {
		return nil, err
	}
	worker := TMWorker{
		logger:   logger,
		cli:      cli,
		dataChan: make(chan *proto.ThingModelMsg, chanCap),
		doneChan: make(chan struct{}),
	}

	return &worker, nil
}

func (w *TMWorker) putData(data *proto.ThingModelMsg) error {
	select {
	case <-time.After(time.Second * 1):
		return errors.New("put data to chan timeout")
	case w.dataChan <- data:
		return nil
	}
}

func (w *TMWorker) run() {
	for {
		select {
		case <-w.doneChan:
			w.stop()
			return
		case data := <-w.dataChan:
			w.logger.Debugf("TMWorker: send thing model report data: %+v", data)
			subCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			if _, err := w.cli.ThingModelMsgReport(subCtx, data); err != nil {
				w.logger.Errorf("thing model msg report error: %s", err)
			}
			cancel()
		}
	}
}

func (w *TMWorker) stop() {
	close(w.doneChan)
}
