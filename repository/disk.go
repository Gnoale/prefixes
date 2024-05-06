package repository

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type memRepository struct {
	words  []string
	sorted bool
	sync.RWMutex
}

func NewMemRepository() Store {
	return &memRepository{
		words: []string{},
	}
}

func (r *memRepository) Insert(word string) error {
	r.Lock()
	defer r.Unlock()
	r.words = append(r.words, word)
	r.sorted = false
	return nil
}

func (r *memRepository) GetByPrefix(prefix string) (*Result, error) {
	r.RLock()
	defer r.RUnlock()
	if !r.sorted {
		sort.Strings(r.words)
		r.sorted = true
	}
	// lookup the index
	index := bsearch(prefix, r.words)
	if index == -1 {
		return nil, fmt.Errorf("no match")
	}
	candidates := []*Result{}
	// inspect our word list from the current index
	candidates = append(candidates, &Result{
		Word:  r.words[index],
		Count: 1,
	})
	for i := index + 1; i < len(r.words); i++ {
		// until it does not longer matches
		if !strings.HasPrefix(r.words[i], prefix) && r.words[i] != prefix {
			break
		}
		// keep count of each different words occurence
		if r.words[i] != candidates[len(candidates)-1].Word {
			candidates = append(candidates, &Result{
				Word:  r.words[i],
				Count: 1,
			})
		} else {
			candidates[len(candidates)-1].Count++
		}
	}
	// sort in decreasing order
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Count > candidates[j].Count
	})
	// return the one with much occurences
	return candidates[0], nil
}
