package ngram

type Storage interface {
	IncrementTokenInHashes(hashes []uint32, id TokenID) error

	CountNGrams(inputNgrams []uint32) (map[TokenID]int, error)
}
