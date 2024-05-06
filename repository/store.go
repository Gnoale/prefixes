package repository

type Result struct {
	Prefix string `json:"prefix"`
	Word   string `json:"word"`
	Count  int    `json:"count"`
}

type Store interface {
	Insert(word string) error
	GetByPrefix(prefix string) (*Result, error)
}
