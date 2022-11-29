package ngram

type MemStorage struct {
	index map[uint32]nGramValue
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		index: map[uint32]nGramValue{},
	}
}

func (m *MemStorage) IncrementTokenInHashes(hashes []uint32, id TokenID) error {
	for _, hash := range hashes {
		if m.index[hash] == nil {
			m.index[hash] = make(map[TokenID]int)
		}
		m.index[hash][id]++
	}
	return nil
}

func (m *MemStorage) CountNGrams(inputNgrams []uint32) (map[TokenID]int, error) {
	counters := make(map[TokenID]int)
	for _, ngramHash := range inputNgrams {
		for tok := range m.index[ngramHash] {
			counters[tok]++
		}
	}
	return counters, nil
}
