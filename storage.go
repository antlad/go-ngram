package ngram

type Storage interface {
	IncrementInHashAndToken(hash uint32, id TokenID)

	CountNGrams(inputNgrams []uint32) map[TokenID]int
}
