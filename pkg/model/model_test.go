package model

import (
	"testing"
	"time"
)

func TestMatchTasksTable_IsNew(t *testing.T) {
	type args struct {
		matchId  string
		oldTable map[string]time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "new for empty",
			args: args{matchId: "#/123456", oldTable: map[string]time.Time{}},
			want: true,
		},
		{
			name: "new for non-empty",
			args: args{matchId: "#/123456", oldTable: map[string]time.Time{
				"#/000000": time.Now(),
				"#/000001": time.Now(),
			}},
			want: true,
		},
		{
			name: "is not new for non-empty",
			args: args{matchId: "#/123456", oldTable: map[string]time.Time{
				"#/000000": time.Now(),
				"#/000001": time.Now(),
				"#/123456": time.Now(),
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mtt := NewMatchTasksTable()
			mtt.table = tt.args.oldTable
			if got := mtt.IsNew(tt.args.matchId); got != tt.want {
				t.Errorf("IsNew() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchTasksTable_Clear(t *testing.T) {
	type args struct {
		before   time.Time
		oldTable map[string]time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
	}{
		{
			name:    "empty table",
			args:    args{before: time.Now(), oldTable: map[string]time.Time{}},
			wantLen: 0,
		},
		{
			name: "clear non-empty table",
			args: args{
				before: time.Now().Add(-1 * time.Hour),
				oldTable: map[string]time.Time{
					"match-1": time.Now(),
					"match-2": time.Now(),
					"match-3": time.Now().Add(-2 * time.Hour),
				},
			},
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mtt := NewMatchTasksTable()
			mtt.table = tt.args.oldTable
			mtt.Clear(tt.args.before)

			if got := len(mtt.table); got != tt.wantLen {
				t.Errorf("Clear() = %v, want %v", got, tt.wantLen)
			}
		})
	}
}
