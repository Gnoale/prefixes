package repository

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/stretchr/testify/assert"
)

var words = []string{"ab", "abc", "azer", "azerty", "turt", "turtle", "turtle"}

func TestBsearch(t *testing.T) {
	sort.Strings(words)
	assert.Equal(t, 2, bsearch("az", words))
	assert.Equal(t, -1, bsearch("nope", words))
	assert.Equal(t, 0, bsearch("ab", words))
	assert.Equal(t, 4, bsearch("tur", words))
}

func TestByPrefixes(t *testing.T) {
	store := NewMemRepository()
	ctx := context.Background()
	for _, word := range words {
		store.Insert(ctx, word)
	}
	v, err := store.GetByPrefix(ctx, "tur")
	assert.NoError(t, err)
	assert.Equal(t, "turtle", v.Word)
	assert.Equal(t, 2, v.Count)

	store = NewMemRepository()
	store.Insert(ctx, "abcd")
	v, err = store.GetByPrefix(ctx, "abc")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", v.Word)
	store.Insert(ctx, "def")
	v, err = store.GetByPrefix(ctx, "de")
	assert.NoError(t, err)
	assert.Equal(t, "def", v.Word)

	store = NewMemRepository()
	store.Insert(ctx, "def")
	store.Insert(ctx, "ijk")
	store.Insert(ctx, "ijm")
	v, err = store.GetByPrefix(ctx, "ij")
	assert.NoError(t, err)
	assert.Equal(t, "ijk", v.Word)
}

func TestConcurentPrefixes(t *testing.T) {
	store := NewMemRepository()
	ctx := context.Background()
	i := 0
	wg := sync.WaitGroup{}
	for i < 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			j := 0
			for j < 100 {
				store.Insert(ctx, petname.Generate(2, ""))
				j++
			}
		}()
		i++
	}
	i = 0
	for i < 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			j := 0
			for j < 100 {
				_, err := store.GetByPrefix(ctx, petname.Adjective())
				if err != nil {
					assert.True(t, errors.Is(err, ErrNotFound))
				}
				j++
			}
		}()
		i++
	}
	wg.Wait()
	v, _ := store.List(ctx)
	fmt.Println(len(v))
}
