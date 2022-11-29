package ngram

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func newIndex(t *testing.T) *NGramIndex {
	const path = "/tmp/ngram_badger_test"
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatal(err)
	}
	st, err := NewBadgerStorage(path)
	require.NoError(t, err)
	index, err := NewNGramIndex(st)
	require.NoError(t, err)
	return index
}

func TestIndexBasics(t *testing.T) {
	index := newIndex(t)

	id1 := TokenID(uuid.NewV4())
	err := index.Add("hello", id1)
	require.NoError(t, err)

	id2 := TokenID(uuid.NewV4())
	err = index.Add("world", id2)
	require.NoError(t, err)

	results, err := index.Search("hello", 0.0)
	require.NoError(t, err)
	require.Len(t, results, 1, "len(results) != 1")

	if results[0].Similarity != 1.0 && results[0].TokenID != id1 {
		t.Error("Bad result")
	}
}

//
//func TestIndexInitialization(t *testing.T) {
//	index, error := NewNGramIndex()
//	if error != nil {
//		t.Error(error)
//	}
//	if index.n != defaultN {
//		t.Error("n is not set to default value")
//	}
//	if index.pad != defaultPad {
//		t.Error("pad is not set to default value")
//	}
//	index, error = NewNGramIndex(SetN(4))
//	if error != nil {
//		t.Error(error)
//	}
//	if index.n != 4 {
//		t.Error("n is not set to 4")
//	}
//	index, error = NewNGramIndex(SetPad('@'))
//	if error != nil {
//		t.Error(error)
//	}
//	if index.pad != "@" {
//		t.Error("pad is not set to @")
//	}
//	// check off limits
//	index, error = NewNGramIndex(SetN(1))
//	if error == nil {
//		t.Error("Error not set (1)")
//	}
//	index, error = NewNGramIndex(SetN(maxN + 1))
//	if error == nil {
//		t.Error("Error not set (2)")
//	}
//}
//
//func BenchmarkAdd(b *testing.B) {
//	b.StopTimer()
//	// init
//	index, _ := NewNGramIndex()
//	var arr []string
//	for i := 0; i < 10000; i++ {
//		str := fmt.Sprintf("%x", i)
//		arr = append(arr, str)
//	}
//	b.StartTimer()
//	for _, hexstr := range arr {
//		index.Add(hexstr)
//	}
//}
//
func TestSearch(t *testing.T) {

	// init
	index := newIndex(t)
	var arr []string
	for i := 0; i < 10000; i++ {
		str := fmt.Sprintf("%000x", i)
		arr = append(arr, str)
	}
	for _, hexstr := range arr {
		index.Add(hexstr, TokenID(uuid.NewV4()))
	}

	for i := 0; i < 10000; i += 4 {
		index.Search(arr[i], 0.5)
	}
}

var pattern = "long_longer_tag_%d"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestSearch2(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var err error

	ti := newIndex(t)

	id1 := TokenID(uuid.NewV4())
	id2 := TokenID(uuid.NewV4())
	id3 := TokenID(uuid.NewV4())

	err = ti.Add("Code is my life", id1) //doc 1
	require.NoError(t, err)
	err = ti.Add("Search", id2) //doc 2
	require.NoError(t, err)
	err = ti.Add("I write a lot of Codes", id3) //doc 3
	require.NoError(t, err)

	for i := 0; i < 100000; i++ {
		err = ti.Add(randSeq(20), TokenID(uuid.NewV4()))
		require.NoError(t, err)
	}

	results, err := ti.Search("Code", 0.7)
	require.NoError(t, err)
	require.Equal(t, id1, results[0].TokenID)
}
