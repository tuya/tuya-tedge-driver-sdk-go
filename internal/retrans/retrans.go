package retrans

import (
	"fmt"
	"strings"
	"time"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/kvdb"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/kvdb/badgdb"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/kvdb/option"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/kvdb/ssort"
)

const (
	ProcessDPInterval     = 2
	ProcessAtopInterval   = 30
	DefaultTTL            = 60 * 60 * 24 * 7
	ReportDpInterval      = 200
	DefaultMaxKeyNum      = 20000
	DefaultMaxDBSize      = 10 * 1024 * 1024
	DefaultEliminationNum = 100
	DefaultMaxAtopValue   = 1024 * 1024

	DefaultDPDBPath   = "/mnt/localdb/dp"
	DefaultAtopDBPath = "/mnt/localdb/atop"
	KeyFormat         = "%s_%d" //DP:   ${cid}_${time}
	AtopKeySep        = "|"     //Atop: ${api}|${version}_${time}
)

var _ kvdb.SDB = new(badgdb.BadgerDB)

type ConnStatItf interface {
	GetCloudStatus() bool
	//SetConnStat(stat bool)
}

type MessageItf interface {
	DPRePublish(key string, data []byte) error
	AtopReReport(key string, data []byte) error
}

type RtOption struct {
	DPDBPath   string //DP重传数据库
	AtopDBPath string //Atop重传数据库
	Interval   int    //DP重传间隔, 毫秒
	TTL        int    //ttl, 秒
	MaxDBSize  int
	MaxKeyNum  int64
}

type ReTransfer struct {
	dpDB       kvdb.SDB
	atopDB     kvdb.SDB
	option     RtOption
	ConnStat   ConnStatItf
	MessagePub MessageItf
	logger     commons.TedgeLogger
	stopChan   chan struct{}
}

func DefaultRtOption() *RtOption {
	rtOption := RtOption{
		DPDBPath:   DefaultDPDBPath,
		AtopDBPath: DefaultAtopDBPath,
		TTL:        DefaultTTL,
		Interval:   ReportDpInterval,
		MaxKeyNum:  DefaultMaxKeyNum,
		MaxDBSize:  DefaultMaxDBSize,
	}

	return &rtOption
}

func NewReTransfer(logger commons.TedgeLogger, roption RtOption) (*ReTransfer, error) {
	reTransfer := ReTransfer{
		logger:   logger,
		option:   roption,
		stopChan: make(chan struct{}),
	}

	var err error
	if roption.DPDBPath != "" {
		reTransfer.dpDB, err = badgdb.NewBadgerDB(roption.DPDBPath, roption.TTL)
		if err != nil {
			return nil, err
		}
	}

	if roption.AtopDBPath != "" {
		reTransfer.atopDB, err = badgdb.NewBadgerDB(roption.AtopDBPath, roption.TTL)
		if err != nil {
			return nil, err
		}
	}

	return &reTransfer, nil
}

func (rt *ReTransfer) SetConnItf(connItf ConnStatItf) {
	rt.ConnStat = connItf
}

func (rt *ReTransfer) SetMessageItf(messagePub MessageItf) {
	rt.MessagePub = messagePub
}

func (rt *ReTransfer) GetConnStat() bool {
	return rt.ConnStat.GetCloudStatus()
}

//单独测试时使用
//func (rt *ReTransfer) SetConnStat(stat bool) {
//	//rt.ConnStat.SetCloudStatus(stat)
//}

///////////////////////////////////////////////////////////////////////////////////////////////
func (rt *ReTransfer) Exit() error {
	rt.dpDB.Close()
	close(rt.stopChan)
	return nil
}

func (rt *ReTransfer) Running() error {
	tickerDP := time.NewTicker(ProcessDPInterval * time.Second)     //mqtt
	tickerAtop := time.NewTicker(ProcessAtopInterval * time.Second) //atop
	for {
		select {
		case <-tickerDP.C:
			if rt.GetConnStat() {
				rt.ProcessDPReTransfer()
			} else {
				rt.ProcessElimination(rt.dpDB)
			}
		case <-tickerAtop.C:
			if rt.GetConnStat() {
				rt.ProcessAtopReTransfer()
			} else {
				rt.ProcessElimination(rt.atopDB)
			}
		case <-rt.stopChan:
			return nil
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//存储数据量太大，淘汰最先保存的部分
func (rt *ReTransfer) ProcessElimination(datadb kvdb.SDB) {
	lc := rt.logger
	if datadb.KeyCount() < rt.option.MaxKeyNum {
		return
	}

	f := func(key, value []byte) bool {
		lc.Errorf("1.ProcessElimination to delete key:%s", key)
		return true
	}

	eOption := option.NewRangeOption(true, 0, DefaultEliminationNum, true, "")
	datadb.RangeSort(f, ssort.SuffixSort, eOption)
}

func genNewKey(key []byte) []byte {
	newKey := fmt.Sprintf(KeyFormat, key, time.Now().UnixNano())
	return []byte(newKey)
}

func getKeyPrefix(key string) string {
	lastIndex := strings.LastIndexByte(key, ssort.Sep[0])
	return key[0:lastIndex]
}

///////////////////////////////////////////////////////////////////////////////////////////////
//DP 模型，DP断网重传
func (rt *ReTransfer) SaveDPKV(key []byte, value []byte) error {
	lc := rt.logger
	count := rt.dpDB.KeyCount()
	if count >= rt.option.MaxKeyNum {
		lc.Errorf("SaveDPKV: discard key:%s, overflow key num:%d, max:%d", key, count, rt.option.MaxKeyNum)
		return fmt.Errorf("overflow key num:%d, max:%d", count, rt.option.MaxKeyNum)
	}

	newKey := genNewKey(key)
	err := rt.dpDB.SetWithDTTL(newKey, value)
	if err != nil {
		lc.Errorf("SaveDPKV SetWithDTTL key:%s, err:%s", newKey, err)
		return err
	}

	return nil
}

//1.如果边缘网关和云端断连，则返回false，退出重传
//2.如果重传接口返回失败，则返回false，退出重传
func (rt *ReTransfer) ProcessDPReTransfer() error {
	lc := rt.logger
	f := func(key, value []byte) bool {
		connStat := rt.GetConnStat()
		if !connStat {
			lc.Errorf("1.ProcessDPReTransfer connStat:%v, return", connStat)
			return false
		}

		prefixKey := getKeyPrefix(string(key))
		err := rt.MessagePub.DPRePublish(prefixKey, value)
		if err != nil {
			lc.Errorf("2.ProcessDPReTransfer DPRePublish key:%s err:%s, return", key, err)
			//1. 只能通过网络连接是否正常决定是否重传成功
			//2. DPRePublish 失败，可能是格式错误，不能返回失败，否则这条数据永远无法发送成功!!!
			//return false
		}
		return true
	}

	rOption := option.NewRangeOption(true, rt.option.Interval, 0, true, "")
	err := rt.dpDB.RangeSort(f, ssort.SuffixSort, rOption)
	if err != nil {
		lc.Errorf("3.ProcessDPReTransfer RangeSort err:%v", err)
	}

	return err
}

func (rt *ReTransfer) PrintDPKeys() error {
	f := func(key, value []byte) bool {
		rt.logger.Infof(">>>PrintDPKeys key:%s, value:%s", key, value)
		return true
	}

	return rt.dpDB.RangeSort(f, ssort.SuffixSort, nil)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//Atop 失败重传
func (rt *ReTransfer) ProcessAtopReTransfer() error {
	lc := rt.logger
	f := func(key, value []byte) bool {
		if !rt.GetConnStat() {
			return false
		}

		prefixKey := getKeyPrefix(string(key))
		err := rt.MessagePub.AtopReReport(prefixKey, value)
		if err != nil {
			lc.Errorf("1.ProcessAtopReTransfer AtopReReport key:%s err:%s, return", key, err)
			//AtopReReport 返回超时、或者Atop网关本身错误时，才能返回err
			return false
		}

		lc.Infof("2.ProcessAtopReTransfer AtopReReport succ, key:%s", key)
		return true
	}

	opt := option.NewRangeOption(true, rt.option.Interval, 0, true, "")
	err := rt.atopDB.RangeSort(f, ssort.SuffixSort, opt)
	if err != nil {
		lc.Errorf("3.ProcessAtopReTransfer RangeSort err:%v", err)
	}
	return err
}

func (rt *ReTransfer) SaveAtopKV(api, version string, value []byte) error {
	lc := rt.logger
	count := rt.atopDB.KeyCount()
	if count >= rt.option.MaxKeyNum {
		lc.Errorf("1.SaveAtopKV: discard api:%s version:%s, overflow key num:%d, max:%d", api, version, count, rt.option.MaxKeyNum)
		return fmt.Errorf("overflow key num:%d, max:%d", count, rt.option.MaxKeyNum)
	}

	valueSize := len(value)
	if valueSize > DefaultMaxAtopValue {
		lc.Errorf("2.SaveAtopKV: discard api:%s version:%s, valueSize:%d too big", api, version, valueSize)
		return fmt.Errorf("api:%s version:%s, valueSize:%d too big", api, version, valueSize)
	}

	key := fmt.Sprintf("%s%s%s", api, AtopKeySep, version)
	newKey := genNewKey([]byte(key))
	err := rt.atopDB.SetWithDTTL(newKey, value)
	if err != nil {
		lc.Errorf("3.SaveAtopKV SetWithDTTL key:%s, err:%s", newKey, err)
		return err
	}

	return nil
}

func (rt *ReTransfer) PrintAtopKeys() error {
	f := func(key, value []byte) bool {
		rt.logger.Infof(">>>PrintAtopKeys key:%s, value:%s", key, value)
		return true
	}

	return rt.atopDB.RangeSort(f, ssort.SuffixSort, nil)
}
