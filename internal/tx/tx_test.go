package tx

import (
	"testing"

	"gotest.tools/assert"
)

func TestTxCommit(t *testing.T) {
	fooCr := &testCommitRollbacker{}
	barCr := &testCommitRollbacker{}
	tx := &uow{
		committers: map[interface{}]Committer{},
	}
	tx.committers["foo"] = fooCr
	tx.committers["bar"] = barCr

	err := tx.Commit()
	assert.NilError(t, err)
	assert.Equal(t, uint32(1), fooCr.commited)
	assert.Equal(t, uint32(1), barCr.commited)
}

func TestTxRollback(t *testing.T) {
	fooCr := &testCommitRollbacker{}
	barCr := &testCommitRollbacker{}
	tx := &uow{
		rollbackers: map[interface{}]Rollbacker{},
	}
	tx.rollbackers["foo"] = fooCr
	tx.rollbackers["bar"] = barCr

	err := tx.Rollback()
	assert.NilError(t, err)
	assert.Equal(t, uint32(1), fooCr.rollbacked)
	assert.Equal(t, uint32(1), barCr.rollbacked)
}

type testCommitRollbacker struct {
	commited   uint32
	rollbacked uint32
}

func (s *testCommitRollbacker) Commit() error {
	s.commited++
	return nil
}

func (s *testCommitRollbacker) Rollback() error {
	s.rollbacked++
	return nil
}

type testCommitRequestRollbacker struct{}

func (s *testCommitRequestRollbacker) CommitRequest() error {
	return nil
}

func (s *testCommitRequestRollbacker) Commit() error {
	return nil
}

func (s *testCommitRequestRollbacker) Rollback() error {
	return nil
}
