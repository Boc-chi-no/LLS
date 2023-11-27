package db

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson"
	"linkshortener/lib/tool"
	"linkshortener/log"
	"linkshortener/setting"
	"strconv"
	"time"
)

type LlsBadgerDB struct {
	DatabaseName string
	ConnectName  string
	BadgerDB     *badger.DB
}

type BadgerDBTable struct {
	tableName string
	db        *LlsBadgerDB
}

func (db *LlsBadgerDB) SetDB(_, _ string) *LlsBadgerDB {
	return db
}

func (db *LlsBadgerDB) SetBadgerDB(badgerDB *badger.DB) *LlsBadgerDB {
	db.BadgerDB = badgerDB
	return db
}

func NewBadgerDB() *LlsBadgerDB {
	log.InfoPrint("Using the BadgerDB as a data source")
	opts := badger.DefaultOptions(tool.If(setting.Cfg.BadgerDB.WithInMemory, "", setting.Cfg.BadgerDB.Path).(string)).WithInMemory(setting.Cfg.BadgerDB.WithInMemory)
	badgerDB, err := badger.Open(opts)
	if err != nil {
		log.PanicPrint("Open BadgerDB File failed: %s", err)
	}
	var db LlsBadgerDB
	return db.SetBadgerDB(badgerDB)
}

// NewBadgerDBTable Initialization table
func NewBadgerDBTable(db *LlsBadgerDB, tableName string) Tabler {
	var table = &BadgerDBTable{}
	table.tableName = tableName
	table.SetDB(db)
	return table
}

// SetBadgerDBTable Setting table
func SetBadgerDBTable(db *LlsBadgerDB, tableName string) Tabler {
	var table = &BadgerDBTable{}
	table.tableName = tableName
	table.SetDB(db)
	return table
}

func (b *BadgerDBTable) SetDB(db interface{}) {
	badgerDB, ok := db.(*LlsBadgerDB)
	if ok {
		b.db = badgerDB
	}
}

func (b *BadgerDBTable) getDB() *badger.DB {
	return b.db.BadgerDB
}

func (b *BadgerDBTable) InsertOne(document interface{}, key string, autoKey bool) error {
	key = tool.ConcatStrings(b.tableName, ":", key)

	if autoKey {
		key = tool.ConcatStrings(key, ":", time.Now().Format("20060102150405"), ":", strconv.FormatUint(tool.GlobalCounterSafeAdd(1), 16))
	}

	db := b.getDB()
	val, err := tool.MarshalJsonByBson(document)
	if err != nil {
		log.ErrorPrint("InsertOne Marshal document error: %v", err)
		return err
	}
	dbErr := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), val)
	})
	if dbErr != nil {
		log.ErrorPrint("InsertOne Update document error: %v", dbErr)
	}
	return dbErr
}

func (b *BadgerDBTable) UpdateOne(filter interface{}, result interface{}) error {
	//TODO: implement me
	log.ErrorPrint("implement UpdateOne")
	return nil
}

func (b *BadgerDBTable) UpdateByID(id string, update interface{}) error {
	if update == nil {
		return log.Errorf("update cannot be nil")
	}
	var updateData map[string]interface{}
	updateDataBson, ok := update.(bson.M)

	if ok {
		updateDataBsonMap := map[string]interface{}(updateDataBson)
		updateDataBsonMapSet, keyExists := updateDataBsonMap["$set"]
		if !keyExists {
			return log.Errorf("update requires $set")
		}
		updateDataBsonMapSetMap, updateDataBsonMapSetOk := updateDataBsonMapSet.(bson.M)
		if updateDataBsonMapSetOk {
			updateData = updateDataBsonMapSetMap
		} else {
			return log.Errorf("update.$set needs to be of type bson.M")
		}
	} else {
		return log.Errorf("update needs to be of type bson.M")
	}

	db := b.getDB()
	key := tool.ConcatStrings(b.tableName, ":", id)

	err := db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		mMap := make(map[string]interface{})
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &mMap)
		})
		if err != nil {
			return log.Errorf("BadgerDB Value Read Error: %s", err)
		}

		for updateKey, updateValue := range updateData {
			_, keyExists := mMap[updateKey]
			if !keyExists {
				return log.Errorf("the Key of update does not exist in the data")
			}
			mMap[updateKey] = updateValue
		}

		newValueBytes, err := json.Marshal(mMap)
		if err != nil {
			return log.Errorf("MarshalJsonByBson Error: %s", err)
		}
		err = txn.Set([]byte(key), newValueBytes)
		if err != nil {
			return log.Errorf("BadgerDB Set Error: %s", err)
		}
		return nil
	})

	return err
}

func (b *BadgerDBTable) FindByID(id interface{}, result interface{}) error {
	//TODO: implement me
	log.ErrorPrint("implement FindByID")
	return nil
}

func (b *BadgerDBTable) FindOne(filter interface{}, result interface{}) error {
	//TODO: implement me
	log.ErrorPrint("implement FindOne")
	return nil
}

func (b *BadgerDBTable) Find(filter interface{}, result interface{}, opt *FindOptions) error {
	var findFilter map[string]interface{}
	if filter != nil {
		findFilterBson, ok := filter.(bson.D)
		if ok {
			findFilter = findFilterBson.Map()
		}
	}

	if opt.Key == "" {
		log.ErrorPrint("BadgerDB requires Key")
		return fmt.Errorf("BadgerDB requires Key")
	}
	db := b.getDB()
	key := tool.ConcatStrings(b.tableName, ":", opt.Key)
	mSlice := make([]map[string]interface{}, 0)

	if opt.PrefixScans {
		err := db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.Prefix = []byte(key)
			it := txn.NewIterator(opts)
			defer it.Close()

			index := 0
			for it.Rewind(); it.Valid() && index < int(opt.Skip); it.Next() {
				index++
			}
			for index < int(opt.Skip+opt.Limit) && it.Valid() {
				item := it.Item()
				mMap := make(map[string]interface{})
				err := item.Value(func(val []byte) error {
					return json.Unmarshal(val, &mMap)
				})
				if err != nil {
					return err
				}

				if findFilter == nil || tool.IsDataMatchingFilter(mMap, findFilter) {
					mSlice = append(mSlice, mMap)
				}

				it.Next()
				index++
			}

			mSliceJson, _ := json.Marshal(mSlice)
			return tool.UnmarshalJsonByBson(mSliceJson, result)
		})
		return err
	} else {
		err := db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(key))
			if err != nil {
				return err
			}
			err = item.Value(func(val []byte) error {
				mMap := make(map[string]interface{})
				err := json.Unmarshal(val, &mMap)
				if err != nil {
					return err
				}
				if findFilter == nil || tool.IsDataMatchingFilter(mMap, findFilter) {
					mSlice = append(mSlice, mMap)
				}
				mSliceJson, _ := json.Marshal(mSlice)
				return tool.UnmarshalJsonByBson(mSliceJson, result)
			})
			return err
		})
		if err != nil {
			log.ErrorPrint("BadgerDB Find Error: %s", err)
		}

		return err
	}
}

func (b *BadgerDBTable) CreateOneIndex(interface{}, ...interface{}) error {
	return nil
}

func (b *BadgerDBTable) CountDocuments(_ interface{}, opt *FindOptions) (int64, error) {
	var count int64 = 0
	if opt.Key == "" {
		log.ErrorPrint("BadgerDB requires Key")
		return count, fmt.Errorf("BadgerDB requires Key")
	}
	db := b.getDB()
	key := tool.ConcatStrings(b.tableName, ":", opt.Key)

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(key)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte(key)); it.ValidForPrefix([]byte(key)); it.Next() {
			count++
		}
		return nil
	})

	return count, err
}
