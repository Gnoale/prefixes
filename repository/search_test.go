package repository

import (
	"fmt"
	"sort"
	"sync"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/stretchr/testify/assert"
)

var words = []string{"ab", "abc", "azer", "azerty", "turtle", "turt", "turtles", "turtle", "a", "aladin", "abra", "abracadabra"}

func TestBsearch(t *testing.T) {
	sort.Strings(words)
	assert.Equal(t, 6, bsearch("az", words))
	assert.Equal(t, 3, bsearch("abra", words))
	assert.Equal(t, 4, bsearch("abrac", words))
	assert.Equal(t, -1, bsearch("nope", words))
	assert.Equal(t, 1, bsearch("ab", words))
	assert.Equal(t, 8, bsearch("tur", words))
	assert.Equal(t, 5, bsearch("ala", words))
}

func TestFindPrefixes(t *testing.T) {
	store := NewMemRepository()
	for _, word := range words {
		store.Insert(word)
	}
	v, _ := store.GetByPrefix("tur")
	assert.Equal(t, "turtle", v.Word)
}

func TestFindSingle(t *testing.T) {
	store := NewMemRepository()

	store.Insert("abcd")
	store.Insert("efg")
	v, err := store.GetByPrefix("abc")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", v.Word)
}

func TestInsertPrefixes(t *testing.T) {
	store := NewMemRepository()
	i := 0
	wg := sync.WaitGroup{}
	for i < 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			j := 0
			for j < 1000 {
				store.Insert(petname.Generate(2, ""))
				j++
			}
		}()
		i++
	}
	wg.Wait()
	v, err := store.GetByPrefix("w")
	assert.NoError(t, err)
	assert.NotEqual(t, nil, v)
	fmt.Println(v.Count, v.Word)
	v, err = store.GetByPrefix("r")
	assert.NoError(t, err)
	assert.NotEqual(t, nil, v)
	fmt.Println(v.Count, v.Word)
}
