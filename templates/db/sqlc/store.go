package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(q *Queries) error) error
}

type SQLStore struct {
	connPoll *pgxpool.Pool
	*Queries
}

func NewStore(connPoll *pgxpool.Pool) Store {
	return &SQLStore{
		connPoll: connPoll,
		Queries:  New(connPoll),
	}
}

func (store *SQLStore) ExecTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := store.connPoll.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p) // re-throw panic after Rollback
		}
	}()

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
