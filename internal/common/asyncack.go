package common

import "sync"

type MsgAckChan struct {
	Mu       sync.Mutex
	Id       string
	IsClosed bool
	DataChan chan interface{}
}

func (mac *MsgAckChan) TryCloseChan() {
	mac.Mu.Lock()
	defer mac.Mu.Unlock()
	if !mac.IsClosed {
		close(mac.DataChan)
		mac.IsClosed = true
	}
}

func (mac *MsgAckChan) TrySendDataAndCloseChan(data interface{}) bool {
	mac.Mu.Lock()
	defer mac.Mu.Unlock()
	if !mac.IsClosed {
		mac.DataChan <- data
		close(mac.DataChan)
		mac.IsClosed = true
		return true
	}
	return false
}
