package worker

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/clients"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/config"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
)

type (
	DPWorker struct {
		logger   commons.TedgeLogger
		cli      *clients.DPReportClient
		dataChan chan interface{}
		doneChan chan struct{}
	}

	DPWorkerPool struct {
		mutex        sync.Mutex
		workers      []*DPWorker
		workersCount uint32
		lastWorker   uint32
	}
)

//////////////////////////////////////////////////////////////////////////////////////////////////////
func NewDPWorkerPool(cliCfg config.ClientInfo, logger commons.TedgeLogger) (*DPWorkerPool, error) {
	var maxWorker uint32 = 16
	var chanCap int32 = 64

	workers := make([]*DPWorker, 0, maxWorker)
	for i := 0; i < int(maxWorker); i++ {
		worker, err := NewDPWorker(cliCfg, chanCap, logger)
		if err != nil {
			return nil, err
		}
		workers = append(workers, worker)
		go worker.run()
	}

	return &DPWorkerPool{
		workersCount: maxWorker,
		lastWorker:   0,
		workers:      workers,
	}, nil
}

func (wp *DPWorkerPool) PutData(data interface{}) error {
	wp.mutex.Lock()
	index := wp.lastWorker % wp.workersCount
	wp.lastWorker++
	wp.mutex.Unlock()
	return wp.workers[index].putData(data)
}

func (wp *DPWorkerPool) Stop() {
	wp.mutex.Lock()
	for _, worker := range wp.workers {
		worker.stop()
	}
	wp.mutex.Unlock()
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
func NewDPWorker(cliCfg config.ClientInfo, chanCap int32, logger commons.TedgeLogger) (*DPWorker, error) {
	cli, err := clients.NewDPReportClient(cliCfg)
	if err != nil {
		return nil, err
	}

	worker := DPWorker{
		logger:   logger,
		cli:      cli,
		dataChan: make(chan interface{}, chanCap),
		doneChan: make(chan struct{}),
	}
	return &worker, nil
}

func (w *DPWorker) putData(data interface{}) error {
	select {
	case <-time.After(time.Second * 1):
		return errors.New("put data to chan timeout")
	case w.dataChan <- data:
		return nil
	}
}

func (w *DPWorker) run() {
	for {
		select {
		case <-w.doneChan:
			w.cli.Conn.Close()
			return
		case data := <-w.dataChan:
			switch data.(type) {
			case *proto.Events:
				subCtx, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout)
				if _, err := w.cli.Add(subCtx, data.(*proto.Events)); err != nil {
					w.logger.Errorf("DPWorker: send with dp data error: %s", err)
					cancel()
				} else {
					cancel()
				}
			case *proto.WithoutDpReport:
				subCtx, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout)
				if _, err := w.cli.WithoutDpReportData(subCtx, data.(*proto.WithoutDpReport)); err != nil {
					w.logger.Errorf("DPWorker: send without dp data error: %s", err)
				}
				cancel()
			default:
				w.logger.Errorf("DPWorker unsupported type:%v", reflect.TypeOf(data))
			}
		}
	}
}

func (w *DPWorker) stop() {
	close(w.doneChan)
}
