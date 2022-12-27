package kvdb

import (
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/kvdb/badgdb"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/kvdb/option"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/kvdb/ssort"
)

var _ SDB = new(badgdb.BadgerDB)

type SDB interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte) error
	Del(key []byte) error

	SetWithTTL(key []byte, value []byte, ttl int) error
	SetWithDTTL(key []byte, value []byte) error
	SetBatch(list map[string][]byte, ttl int) error

	Range(f func(key, value []byte) bool, opt *option.RangeOption) error
	RangeSort(f func(key, value []byte) bool, f2 ssort.SSort, opt *option.RangeOption) error

	DelBatch(keys ...[]byte) error

	KeyCount() int64
	Close() error
}
