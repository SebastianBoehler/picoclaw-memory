package memory

import "context"

type Store interface {
	Save(context.Context, Memory) error
	Recall(context.Context, RecallRequest) ([]RecallResult, error)
	ListRecent(context.Context, string, int) ([]Memory, error)
}
