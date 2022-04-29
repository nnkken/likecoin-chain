package staking

import (
	dbm "github.com/tendermint/tm-db"
)

func (q *Querier) batchFlusher() {
	for {
		b := <-q.batchFlushQueue
		err := b.Write()
		if err != nil {
			panic(err)
		}
		err = b.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (q *Querier) getCurrentBatch() dbm.Batch {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.batch
}

func (q *Querier) newBatch() dbm.Batch {
	q.lock.Lock()
	defer q.lock.Unlock()
	b := q.batch
	q.batch = q.indexingDB.NewBatch()
	return b
}

func (q *Querier) queueBatch() dbm.Batch {
	b := q.newBatch()
	q.batchFlushQueue <- b
	return b
}

func (q *Querier) batchSet(key, value []byte) {
	err := q.getCurrentBatch().Set(key, value)
	if err != nil {
		panic(err)
	}
}

func (q *Querier) batchDelete(key []byte) {
	err := q.getCurrentBatch().Delete(key)
	if err != nil {
		panic(err)
	}
}
