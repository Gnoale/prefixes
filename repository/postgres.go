package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepostory struct {
	db *Queries
}

func NewPGRepository(connString string) (Store, error) {
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &pgRepostory{New(dbpool)}, nil

}

func (r *pgRepostory) Insert(ctx context.Context, word string) error {
	return r.db.InsertWord(ctx, word)
}

// GetByPrefix implements Store.
func (r *pgRepostory) GetByPrefix(ctx context.Context, prefix string) (*Result, error) {
	row, err := r.db.GetByPrefix(ctx, prefix+"%")
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &Result{
		Word:  row.Word,
		Count: int(row.Count),
	}, nil

}

// List implements Store.
func (r *pgRepostory) List(ctx context.Context) ([]Result, error) {
	rows, err := r.db.List(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]Result, len(rows))
	for i := range rows {
		results[i] = Result{
			Word:  rows[i].Word,
			Count: int(rows[i].Count),
		}
	}
	return results, nil
}
