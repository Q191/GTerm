package services

import (
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

func (s *CredentialSrv) ListCredential() *resp.Resp {
	t := s.Query.Credential
	credentials, err := t.Where(t.IsCommonCredential.Is(true)).Find()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithData(credentials)
}
