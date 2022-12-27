package badgdb

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/kvdb/option"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/kvdb/ssort"
	"go.uber.org/atomic"
)

type BadgerDB struct {
	bagdb      *badger.DB
	keyCount   atomic.Int64
	dbPath     string
	defaultTTL int
}

func NewBadgerDB(dbpath string, dttl int) (*BadgerDB, error) {
	option := badger.DefaultOptions(dbpath).WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(option)
	if err != nil {
		return nil, err
	}

	badged := BadgerDB{
		defaultTTL: dttl,
		bagdb:      db,
		dbPath:     dbpath,
	}

	badged.Init()
	return &badged, nil
}

func (db *BadgerDB) Init() {
	count, _ := db.initKeyCount()
	db.keyCount.Store(count)
}

func (db *BadgerDB) Close() error {
	return db.bagdb.Close()
}

func (db *BadgerDB) KeyCount() int64 {
	return db.keyCount.Load()
}

func (db *BadgerDB) Get(key []byte) ([]byte, error) {
	var valCopy []byte
	err := db.bagdb.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		valCopy, err = item.ValueCopy(nil)
		return err
	})

	return valCopy, err
}

func (db *BadgerDB) Set(key []byte, value []byte) error {
	err := db.bagdb.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, value)
		err := txn.SetEntry(e)
		return err
	})

	if err == nil {
		db.keyCount.Inc()
	}

	return err
}

func (db *BadgerDB) SetBatch(list map[string][]byte, ttl int) error {
	return db.bagdb.Update(func(txn *badger.Txn) error {
		for key, val := range list {
			entry := badger.NewEntry([]byte(key), val)
			if ttl > 0 {
				entry.WithTTL(time.Duration(ttl) * time.Second)
			}

			if err := txn.SetEntry(entry); err != nil {
				return fmt.Errorf("set entry error %w", err)
			}
		}
		return nil
	})
}

func (db *BadgerDB) Del(key []byte) error {
	err := db.bagdb.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})

	if err == nil {
		db.keyCount.Dec()
	}

	return err
}

func (db *BadgerDB) DelBatch(keys ...[]byte) error {
	err := db.bagdb.Update(func(txn *badger.Txn) error {
		for _, key := range keys {
			if err := txn.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})

	if err == nil {
		db.keyCount.Sub(int64(len(keys)))
	}
	return err
}

// SetWithTTL 带ttl 存储，ttl<=0是永久存储
func (db *BadgerDB) SetWithTTL(key []byte, value []byte, ttl int) error {
	err := db.bagdb.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, value)
		if ttl > 0 {
			entry.WithTTL(time.Duration(ttl) * time.Second)
		}
		return txn.SetEntry(entry)
	})

	if err == nil {
		db.keyCount.Inc()
	}

	return err
}

func (db *BadgerDB) SetWithDTTL(key []byte, value []byte) error {
	return db.SetWithTTL(key, value, db.defaultTTL)
}

func (db *BadgerDB) Range(f func(key, value []byte) bool, opt *option.RangeOption) error {
	return db.RangeWithOption(f, nil, opt)
}

//先调用f2排序、再调用f处理数据
func (db *BadgerDB) RangeSort(f func(key, value []byte) bool, f2 ssort.SSort, opt *option.RangeOption) error {
	return db.RangeWithOption(f, f2, opt)
}

func (db *BadgerDB) RangeWithOption(f func(key, value []byte) bool, f2 ssort.SSort, opt *option.RangeOption) error {
	if opt == nil {
		opt = option.DefaultOption()
	}

	//1.Get Keys
	var keys []string
	if opt.Prefix == "" {
		keys = db.getAllKeys()
	} else {
		keys = db.getPrefixKeys(opt.Prefix)
	}

	//2.sort and process keys/values
	if f2 != nil {
		f2(keys)
	}

	//3.process keys/values
	return db.processWithOption(f, keys, *opt)
}

////////////////////////////////////////////////////////////////////////////////////////////
func (db *BadgerDB) getAllKeys() []string {
	var keys []string
	db.bagdb.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().KeyCopy(nil)
			keys = append(keys, string(key))
		}

		return nil
	})

	return keys
}

func (db *BadgerDB) getPrefixKeys(prefix string) []string {
	var keys []string
	db.bagdb.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		prefixB := []byte(prefix)
		for it.Seek(prefixB); it.ValidForPrefix(prefixB); it.Next() {
			key := it.Item().KeyCopy(nil)
			keys = append(keys, string(key))
		}
		return nil
	})

	return keys
}

func (db *BadgerDB) initKeyCount() (int64, error) {
	var count int64 = 0
	err := db.bagdb.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			count++
		}

		return nil
	})

	return count, err
}

func (db *BadgerDB) processWithOption(f func(key, value []byte) bool, keys []string, op option.RangeOption) error {
	num := 0
	var err error
	for _, key := range keys {
		err = db.bagdb.Update(func(txn *badger.Txn) error {
			if op.SleepMs > 0 {
				time.Sleep(time.Duration(op.SleepMs) * time.Millisecond)
			}

			//1.Get Value
			item, err := txn.Get([]byte(key))
			if err != nil {
				return fmt.Errorf("get key:%s err:%s", key, err)
			}

			//2.ValueCopy
			var value []byte
			if op.LoadValue {
				value, err = item.ValueCopy(nil)
				if err != nil {
					return fmt.Errorf("key:%s copy value err:%s", key, err)
				}

				//3.process Value
				ret := f([]byte(key), value)
				if !ret {
					return fmt.Errorf("process ret:%v", ret)
				}
			}

			//4.delete key
			if op.DelFlag {
				err := txn.Delete([]byte(key))
				if err == nil {
					db.keyCount.Dec()
				} else {
					return fmt.Errorf("key:%s delete err:%s", key, err)
				}
			}

			return nil
		})

		if err != nil {
			return err
		}

		num++
		if op.MaxCount > 0 && num >= op.MaxCount {
			break
		}
	}

	return err
}
