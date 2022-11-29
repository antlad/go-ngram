package ngram

type MemStorage struct {
	index map[uint32]nGramValue
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		index: map[uint32]nGramValue{},
	}
}

func (m *MemStorage) IncrementInHashAndToken(hash uint32, id TokenID) {
	if m.index[hash] == nil {
		m.index[hash] = make(map[TokenID]int)
	}
	m.index[hash][id]++
}

func (m *MemStorage) CountNGrams(inputNgrams []uint32) map[TokenID]int {
	counters := make(map[TokenID]int)
	for _, ngramHash := range inputNgrams {
		for tok := range m.index[ngramHash] {
			counters[tok]++
		}
	}
	return counters
}
