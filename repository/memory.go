package repository

import (
	"context"
	"sort"
	"strings"
	"sync"
)

type memRepository struct {
	words     []string
	reference map[string]int
	sorted    bool
	sync.RWMutex
}

func NewMemRepository() Store {
	return &memRepository{
		words:     []string{},
		reference: map[string]int{},
	}
}

func (r *memRepository) Insert(ctx context.Context, word string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.reference[word]; ok {
		r.reference[word]++
	} else {
		r.reference[word] = 1
		r.words = append(r.words, word)
		r.sorted = false
	}
	return nil
}

func (r *memRepository) GetByPrefix(ctx context.Context, prefix string) (*Result, error) {
	r.Lock()
	defer r.Unlock()
	if !r.sorted {
		sort.Strings(r.words)
		r.sorted = true
	}
	// lookup the index
	i := bsearch(prefix, r.words)
	if i == -1 {
		return nil, ErrNotFound
	}
	result := &Result{}
	// run through the word list from the first prefix index
	for ; i < len(r.words); i++ {
		// until it does not longer matches
		if !strings.HasPrefix(r.words[i], prefix) && r.words[i] != prefix {
			break
		}
		// keep best
		if count := r.reference[r.words[i]]; count > result.Count {
			result = &Result{
				Word:  r.words[i],
				Count: count,
			}
		}
	}
	return result, nil
}

func (r *memRepository) List(ctx context.Context) ([]Result, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.words) == 0 {
		return []Result{}, nil
	}
	if !r.sorted {
		sort.Strings(r.words)
		r.sorted = true
	}
	results := []Result{}
	for i := 0; i < len(r.words); i++ {
		results = append(results, Result{
			Word:  r.words[i],
			Count: r.reference[r.words[i]],
		})
	}
	return results, nil
}

// return the index of the smallest prefix match from the word list
// -1 if there are no match at all
func bsearch(prefix string, words []string) int {
	if len(words) == 0 {
		return -1
	}
	left := 0
	right := len(words) - 1
	i := (left + right) / 2
	for left <= right {
		if words[i] < prefix {
			left = i + 1
		} else if words[i] > prefix {
			right = i - 1
		} else {
			return i
		}
		i = (left + right) / 2
	}
	// no exact match, the first prefix is either at i or i+1
	if strings.HasPrefix(words[i], prefix) {
		return i
	}
	if i < len(words)-1 {
		if strings.HasPrefix(words[i+1], prefix) {
			return i + 1
		}
	}
	// not found
	return -1
}
