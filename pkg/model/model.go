package model

import (
	"sync"
	"time"
)

type MatchTasksTable struct {
	table    map[string]time.Time // assoc matchId -> addedAt
	tableMtx sync.RWMutex
}

func (mtt MatchTasksTable) IsNew(matchId string) bool {
	mtt.tableMtx.RLock()
	_, exists := mtt.table[matchId]
	mtt.tableMtx.RUnlock()
	if exists {
		return false
	}

	mtt.tableMtx.Lock()
	mtt.table[matchId] = time.Now()
	mtt.tableMtx.Unlock()

	return true
}

func (mtt MatchTasksTable) Clear(before time.Time) {
	for matchId, addedAt := range mtt.table {
		if addedAt.Before(before) {
			mtt.tableMtx.Lock()
			delete(mtt.table, matchId)
			mtt.tableMtx.Unlock()
		}
	}
}
