package services

import (
	"github.com/OpenToolkitLab/GTerm/backend/dal/model"
	"github.com/OpenToolkitLab/GTerm/backend/dal/query"
	"github.com/OpenToolkitLab/GTerm/backend/pkg/resp"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var HostSrvSet = wire.NewSet(wire.Struct(new(HostSrv), "*"))

type HostSrv struct {
	Logger *zap.Logger
	Query  *query.Query
}

func (s *GroupSrv) CreateHost(host *model.Host) *resp.Resp {
	t := s.Query.Host
	if err := t.Create(host); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *GroupSrv) UpdateHost(host *model.Host) *resp.Resp {
	t := s.Query.Host
	if _, err := t.Where(t.ID.Eq(host.ID)).Updates(host); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *GroupSrv) DeleteHost(id uint) *resp.Resp {
	t := s.Query.Host
	_, err := t.Where(t.ID.Eq(id)).Delete()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *GroupSrv) ListHost() *resp.Resp {
	t := s.Query.Host
	hosts, err := t.Find()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithData(hosts)
}
