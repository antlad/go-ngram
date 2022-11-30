package ngram

import (
	"errors"
	"github.com/spaolacci/murmur3"
	"math"
)

const (
	maxN       = 8
	defaultPad = "$"
	defaultN   = 3
)

// TokenID is just id of the token
type TokenID int

type nGramValue map[TokenID]int

// NGramIndex can be initialized by default (zeroed) or created with "NewNgramIndex"
type NGramIndex struct {
	pad     string
	n       int
	storage Storage
	warp    float64
}

// SearchResult contains token id and similarity - value in range from 0.0 to 1.0
type SearchResult struct {
	TokenID    TokenID
	Similarity float64
}

func (ngram *NGramIndex) splitInput(str string) ([]uint32, error) {
	if len(str) == 0 {
		return nil, errors.New("empty string")
	}
	pad := ngram.pad
	n := ngram.n
	input := pad + str + pad
	prevIndexes := make([]int, maxN)
	var counter int
	results := make([]uint32, 0)

	for index := range input {
		counter++
		if counter > n {
			top := prevIndexes[(counter-n)%maxN]
			substr := input[top:index]
			hash := murmur3.Sum32([]byte(substr))
			results = append(results, hash)
		}
		prevIndexes[counter%maxN] = index
	}

	for i := n - 1; i > 1; i-- {
		if len(input) >= i {
			top := prevIndexes[(len(input)-i)%maxN]
			substr := input[top:]
			hash := murmur3.Sum32([]byte(substr))
			results = append(results, hash)
		}
	}

	return results, nil
}

func (ngram *NGramIndex) init() {
	if ngram.pad == "" {
		ngram.pad = defaultPad
	}
	if ngram.n == 0 {
		ngram.n = defaultN
	}
	if ngram.warp == 0.0 {
		ngram.warp = 1.0
	}
}

type Option func(*NGramIndex) error

// SetPad must be used to pass padding character to NGramIndex c-tor
func SetPad(c rune) Option {
	return func(ngram *NGramIndex) error {
		ngram.pad = string(c)
		return nil
	}
}

// SetN must be used to pass N (gram size) to NGramIndex c-tor
func SetN(n int) Option {
	return func(ngram *NGramIndex) error {
		if n < 2 || n > maxN {
			return errors.New("bad 'n' value for n-gram index")
		}
		ngram.n = n
		return nil
	}
}

// SetWarp must be used to pass warp to NGramIndex c-tor
func SetWarp(warp float64) Option {
	return func(ngram *NGramIndex) error {
		if warp < 0.0 || warp > 1.0 {
			return errors.New("bad 'warp' value for n-gram index")
		}
		ngram.warp = warp
		return nil
	}
}

// NewNGramIndex is N-gram index c-tor. In most cases must be used withot parameters.
// You can pass parameters to c-tor using functions SetPad, SetWarp and SetN.
func NewNGramIndex(storage Storage, opts ...Option) (*NGramIndex, error) {
	ngram := &NGramIndex{
		storage: storage,
	}
	for _, opt := range opts {
		if err := opt(ngram); err != nil {
			return nil, err
		}
	}
	ngram.init()
	return ngram, nil
}

// Add token to index. Function returns token id, this id can be converted
// to string with function "GetString".
func (ngram *NGramIndex) Add(input string, id TokenID) error {
	results, err := ngram.splitInput(input)
	if err != nil {
		return err
	}

	return ngram.storage.IncrementTokenInHashes(results, id)
}

// countNgrams maps matched tokens to the number of ngrams, shared with input string
func (ngram *NGramIndex) countNgrams(inputNgrams []uint32) (map[TokenID]int, error) {
	return ngram.storage.CountNGrams(inputNgrams)
}

func validateThresholdValues(thresholds []float64) (float64, error) {
	var tval float64
	if len(thresholds) == 1 {
		tval = thresholds[0]
		if tval < 0.0 || tval > 1.0 {
			return 0.0, errors.New("threshold must be in range (0, 1)")
		}
	} else if len(thresholds) > 1 {
		return 0.0, errors.New("too many arguments")
	}
	return tval, nil
}

func (ngram *NGramIndex) match(input string, tval float64) ([]SearchResult, error) {
	inputNgrams, err := ngram.splitInput(input)
	if err != nil {
		return nil, err
	}
	output := make([]SearchResult, 0)
	tokenCount, err := ngram.countNgrams(inputNgrams)
	if err != nil {
		return nil, err
	}
	for token, count := range tokenCount {
		var sim float64
		allngrams := float64(len(inputNgrams))
		matchngrams := float64(count)
		if ngram.warp == 1.0 {
			sim = matchngrams / allngrams
		} else {
			diffngrams := allngrams - matchngrams
			sim = math.Pow(allngrams, ngram.warp) - math.Pow(diffngrams, ngram.warp)
			sim /= math.Pow(allngrams, ngram.warp)
		}
		if sim >= tval {
			res := SearchResult{Similarity: sim, TokenID: token}
			output = append(output, res)
		}
	}
	return output, nil
}

// Search for matches between query string (input) and indexed strings.
// First parameter - threshold is optional and can be used to set minimal similarity
// between input string and matching string. You can pass only one threshold value.
// Results is an unordered array of 'SearchResult' structs. This struct contains similarity
// value (float32 value from threshold to 1.0) and token-id.
func (ngram *NGramIndex) Search(input string, threshold ...float64) ([]SearchResult, error) {
	tval, err := validateThresholdValues(threshold)
	if err != nil {
		return nil, err
	}
	return ngram.match(input, tval)
}

// BestMatch is the same as Search except that it's returning only one best result instead of all.
func (ngram *NGramIndex) BestMatch(input string, threshold ...float64) (*SearchResult, error) {
	tval, err := validateThresholdValues(threshold)
	if err != nil {
		return nil, err
	}
	variants, err := ngram.match(input, tval)
	if err != nil {
		return nil, err
	}
	if len(variants) == 0 {
		return nil, errors.New("no matches found")
	}
	var result SearchResult
	maxsim := -1.0
	for _, val := range variants {
		if val.Similarity > maxsim {
			maxsim = val.Similarity
			result = val
		}
	}
	return &result, nil
}
