package delta

import (
	"sync"
	"time"
)

type matchData struct {
	scores    [2]uint
	updatedAt time.Time
}

// ScoresTableDelta helps make diff of matches scores
type ScoresTableDelta struct {
	table    map[string]matchData // matchID -> scores hash
	tableMtx sync.Mutex
	buffer   map[string][2]uint // matchID -> scores hash
}

// Append method fills buffer with new data
func (std ScoresTableDelta) Append(matchId string, homeScore, awayScore uint) {
	std.buffer[matchId] = [2]uint{homeScore, awayScore}
}

// FlushAndGetDiff method use events table and events buffer to create DELTA
func (std ScoresTableDelta) FlushAndGetDiff() map[string][2]uint {
	delta := map[string][2]uint{}
	for matchID, actualScores := range std.buffer {
		oldScores, exists := std.table[matchID]
		if !exists {
			std.tableMtx.Lock()
			std.table[matchID] = matchData{scores: actualScores, updatedAt: time.Now()}
			std.tableMtx.Unlock()
			delta[matchID] = actualScores
			continue
		}
		if actualScores[0] != oldScores.scores[0] || actualScores[1] != oldScores.scores[1] {
			std.tableMtx.Lock()
			std.table[matchID] = matchData{scores: actualScores, updatedAt: time.Now()}
			std.tableMtx.Unlock()
			delta[matchID] = actualScores
		}
	}
	return delta
}

// Clear method help remove old events
func (std ScoresTableDelta) Clear(before time.Time) {
	for matchID, data := range std.table {
		if data.updatedAt.Before(before) {
			std.tableMtx.Lock()
			delete(std.table, matchID)
			std.tableMtx.Unlock()
		}
	}
}

func NewScoresTableDelta() *ScoresTableDelta {
	return &ScoresTableDelta{
		table:    map[string]matchData{},
		tableMtx: sync.Mutex{},
		buffer:   map[string][2]uint{},
	}
}
