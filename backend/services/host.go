package services

import (
	"github.com/MisakaTAT/GTerm/backend/dal/model"
	"github.com/MisakaTAT/GTerm/backend/dal/query"
	"github.com/MisakaTAT/GTerm/backend/pkg/resp"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var HostSrvSet = wire.NewSet(wire.Struct(new(HostSrv), "*"))

type HostSrv struct {
	Logger *zap.Logger
	Query  *query.Query
}

func (s *HostSrv) CreateHost(host *model.Host) *resp.Resp {
	t := s.Query.Host
	if err := t.Create(host); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *HostSrv) UpdateHost(host *model.Host) *resp.Resp {
	t := s.Query.Host
	if _, err := t.Where(t.ID.Eq(host.ID)).Updates(host); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *HostSrv) FindByID(id uint) (*model.Host, error) {
	t := s.Query.Host
	return t.Where(t.ID.Eq(id)).Preload(t.Credential, t.Metadata).First()
}

func (s *HostSrv) DeleteHost(id uint) *resp.Resp {
	t := s.Query.Host
	_, err := t.Where(t.ID.Eq(id)).Delete()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *HostSrv) ListHost() *resp.Resp {
	t := s.Query.Host
	hosts, err := t.Preload(t.Metadata).Find()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithData(hosts)
}
