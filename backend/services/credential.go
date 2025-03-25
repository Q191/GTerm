package services

import (
	"github.com/MisakaTAT/GTerm/backend/dal/model"
	"github.com/MisakaTAT/GTerm/backend/dal/query"
	"github.com/MisakaTAT/GTerm/backend/utils/resp"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var CredentialSrvSet = wire.NewSet(wire.Struct(new(CredentialSrv), "*"))

type CredentialSrv struct {
	Logger *zap.Logger
	Query  *query.Query
}

func (s *CredentialSrv) CreateCredential(cred *model.Credential) *resp.Resp {
	t := s.Query.Credential
	cred.IsCommonCredential = true
	if err := t.Create(cred); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *CredentialSrv) UpdateCredential(cred *model.Credential) *resp.Resp {
	t := s.Query.Credential
	if _, err := t.Where(t.ID.Eq(cred.ID)).Updates(cred); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *CredentialSrv) ListCredential() *resp.Resp {
	t := s.Query.Credential
	credentials, err := t.Where(t.IsCommonCredential.Is(true)).Find()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithData(credentials)
}

func (s *CredentialSrv) DeleteCredential(id uint) *resp.Resp {
	t := s.Query.Credential
	_, err := t.Where(t.ID.Eq(id)).Delete()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}
