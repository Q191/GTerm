package services

import (
	"github.com/OpenToolkitLab/GTerm/backend/dal/model"
	"github.com/OpenToolkitLab/GTerm/backend/dal/query"
	"github.com/OpenToolkitLab/GTerm/backend/pkg/resp"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var GroupSrvSet = wire.NewSet(wire.Struct(new(GroupSrv), "*"))

type GroupSrv struct {
	Logger *zap.Logger
	Query  *query.Query
}

func (s *GroupSrv) CreateGroup(group *model.Group) *resp.Resp {
	t := s.Query.Group
	if err := t.Create(group); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *GroupSrv) DeleteGroup(id uint) *resp.Resp {
	t := s.Query.Group
	_, err := t.Where(t.ID.Eq(id)).Delete()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *GroupSrv) ListGroup() *resp.Resp {
	t := s.Query.Group
	groups, err := t.Find()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithData(groups)
}
