package ngram

import (
	"encoding/binary"
	"errors"
	"github.com/dgraph-io/badger/v3"
	uuid "github.com/satori/go.uuid"
)

type BadgerStorage struct {
	db *badger.DB
}

func key(hash uint32, id TokenID) []byte {
	k := make([]byte, 4+uuid.Size)
	binary.LittleEndian.PutUint32(k, hash)
	copy(k[4:], id[:])
	return k
}

func tokenFromKey(k []byte) TokenID {
	var id TokenID
	copy(id[:], k[4:])
	return id
}

func value(val uint32) []byte {
	v := make([]byte, 4)
	binary.LittleEndian.PutUint32(v, val)
	return v
}

func hashPrefix(hash uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, hash)
	return b
}

func NewBadgerStorage(badgerFolder string) (*BadgerStorage, error) {
	opts := badger.DefaultOptions(badgerFolder)
	opts.NumVersionsToKeep = 0
	opts.CompactL0OnClose = true
	opts.NumLevelZeroTables = 1
	opts.NumLevelZeroTablesStall = 2
	opts.ValueLogFileSize = 1024 * 1024 * 10

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerStorage{
		db: db,
	}, nil
}

func iteratorOptsOnPrefix(prefix []byte) badger.IteratorOptions {
	opts := badger.DefaultIteratorOptions
	opts.Prefix = prefix
	return opts
}

func (m *BadgerStorage) IncrementTokenInHashes(hashes []uint32, id TokenID) error {
	return m.db.Update(func(txn *badger.Txn) error {
		for _, hash := range hashes {
			k := key(hash, id)
			item, err := txn.Get(k)
			if err != nil {
				if errors.Is(err, badger.ErrKeyNotFound) {
					if err = txn.Set(k, value(1)); err != nil {
						return err
					}
					continue
				}
				return err
			}
			if err = item.Value(func(val []byte) error {
				v := binary.LittleEndian.Uint32(val)
				v++
				binary.LittleEndian.PutUint32(val, v)
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	})

}

func (m *BadgerStorage) CountNGrams(inputNgrams []uint32) (map[TokenID]int, error) {
	result := make(map[TokenID]int)
	err := m.db.View(func(txn *badger.Txn) error {
		for _, ngramHash := range inputNgrams {
			it := txn.NewIterator(iteratorOptsOnPrefix(hashPrefix(ngramHash)))
			for it.Rewind(); it.Valid(); it.Next() {
				id := tokenFromKey(it.Item().Key())
				result[id]++
			}
			it.Close()
		}
		return nil
	})
	return result, err
}
