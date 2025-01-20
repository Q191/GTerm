package services

import (
	"github.com/MisakaTAT/GTerm/backend/dal/model"
	"github.com/MisakaTAT/GTerm/backend/dal/query"
	"github.com/MisakaTAT/GTerm/backend/utils/resp"
	"github.com/google/uuid"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var HostSrvSet = wire.NewSet(wire.Struct(new(HostSrv), "*"))

type HostSrv struct {
	Logger *zap.Logger
	Query  *query.Query
}

func (s *HostSrv) CreateHost(host *model.Host) *resp.Resp {
	if err := s.Query.Transaction(func(tx *query.Query) error {
		if host.CredentialID == nil && host.Credential != nil {
			host.Credential.Name = uuid.New().String()
			if err := tx.Credential.Create(host.Credential); err != nil {
				return err
			}
			host.CredentialID = &host.Credential.ID
		}
		if err := tx.Host.Create(host); err != nil {
			return err
		}
		return nil
	}); err != nil {
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
	if err := s.Query.Transaction(func(tx *query.Query) error {
		host, err := tx.Host.Where(tx.Host.ID.Eq(id)).First()
		if err != nil {
			return err
		}
		if !host.IsCommonCredential && host.CredentialID != nil {
			if _, err = tx.Credential.Where(tx.Credential.ID.Eq(*host.CredentialID)).Unscoped().Delete(); err != nil {
				return err
			}
		}
		if _, err = tx.Host.Select(tx.Host.Metadata.Field()).Unscoped().Delete(host); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *HostSrv) ListHost() *resp.Resp {
	t := s.Query.Host
	hosts, err := t.Preload(t.Metadata, t.Credential).Find()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithData(hosts)
}
