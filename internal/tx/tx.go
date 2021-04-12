package tx

import (
	"context"
	"errors"
)

var (
	// ErrTxNotFound error
	ErrTxNotFound = errors.New("tx not found")
)

var (
	txKey = ctxKey("tx")
)

type ctxKey string

// Preparer interface
type Preparer interface {
	PrepareCommit() error
}

// Committer interface
type Committer interface {
	Commit() error
}

// Rollbacker interface
type Rollbacker interface {
	Rollback() error
}

// CommitRollbacker interface
type CommitRollbacker interface {
	Committer
	Rollbacker
}

type uow struct {
	crs         map[interface{}]CommitRollbacker
	preparer    map[interface{}]Preparer
	committers  map[interface{}]Committer
	rollbackers map[interface{}]Rollbacker
}

// Get commit rollbacker
func (u *uow) Get(id interface{}, fn func() (CommitRollbacker, error)) (CommitRollbacker, error) {
	cr, ok := u.crs[id]
	if ok {
		return cr, nil
	}

	cr, err := fn()
	if err != nil {
		return nil, err
	}

	u.crs[id] = cr
	u.committers[id] = cr
	u.rollbackers[id] = cr

	creq, ok := cr.(Preparer)
	if ok {
		u.preparer[id] = creq
	}

	return cr, nil
}

// Commit tx
func (u *uow) Commit() error {
	for _, creq := range u.preparer {
		if err := creq.PrepareCommit(); err != nil {
			return err
		}
	}

	for _, committer := range u.committers {
		if err := committer.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// Rollback tx
func (u *uow) Rollback() error {
	var result error

	for _, rollbacker := range u.rollbackers {
		err := rollbacker.Rollback()
		if result == nil && err != nil {
			result = err
		}
	}

	return result
}

// Acquire tx
func Acquire(ctx context.Context, id interface{}, fn func() (CommitRollbacker, error)) (CommitRollbacker, error) {
	if ctx == nil {
		return nil, ErrTxNotFound
	}

	tx, ok := ctx.Value(txKey).(*uow)
	if !ok {
		return nil, ErrTxNotFound
	}

	t, err := tx.Get(id, fn)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Run tx
func Run(ctx context.Context, fn func(context.Context) error) error {
	tx := &uow{
		crs:         map[interface{}]CommitRollbacker{},
		preparer:    map[interface{}]Preparer{},
		committers:  map[interface{}]Committer{},
		rollbackers: map[interface{}]Rollbacker{},
	}
	ctx = context.WithValue(ctx, txKey, tx)

	if err := fn(ctx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
