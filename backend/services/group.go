package services

import (
	"context"
	"github.com/OpenToolkitLab/GTerm/backend/dal/model"
	"github.com/OpenToolkitLab/GTerm/backend/dal/query"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var GroupSrvSet = wire.NewSet(wire.Struct(new(GroupSrv), "*"))

type GroupSrv struct {
	Logger *zap.Logger
	Query  *query.Query
}

func (s *GroupSrv) Create(group *model.Group, ctx context.Context) error {
	t := s.Query.Group
	q := t.WithContext(ctx)
	return q.Create(group)
}
