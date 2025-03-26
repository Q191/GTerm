package services

import (
	"github.com/MisakaTAT/GTerm/backend/dal/model"
	"github.com/MisakaTAT/GTerm/backend/dal/query"
	"github.com/MisakaTAT/GTerm/backend/utils/resp"
	"github.com/google/uuid"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var ConnectionSrvSet = wire.NewSet(wire.Struct(new(ConnectionSrv), "*"))

type ConnectionSrv struct {
	Logger *zap.Logger
	Query  *query.Query
}

func (s *ConnectionSrv) CreateConnection(conn *model.Connection) *resp.Resp {
	if err := s.Query.Transaction(func(tx *query.Query) error {
		if conn.CredentialID == nil && conn.Credential != nil {
			conn.Credential.IsCommonCredential = false
			conn.Credential.Label = uuid.New().String()
			if err := tx.Credential.Create(conn.Credential); err != nil {
				return err
			}
			conn.CredentialID = &conn.Credential.ID
		}
		if err := tx.Connection.Create(conn); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *ConnectionSrv) UpdateConnection(conn *model.Connection) *resp.Resp {
	if err := s.Query.Transaction(func(tx *query.Query) error {
		oldConn, err := tx.Connection.Where(tx.Connection.ID.Eq(conn.ID)).First()
		if err != nil {
			return err
		}
		if oldConn.UseCommonCredential && !conn.UseCommonCredential {
			if conn.Credential != nil {
				conn.Credential.IsCommonCredential = false
				conn.Credential.Label = uuid.New().String()
				if err = tx.Credential.Create(conn.Credential); err != nil {
					return err
				}
				conn.CredentialID = &conn.Credential.ID
			}
		}

		if !oldConn.UseCommonCredential && conn.UseCommonCredential {
			if oldConn.CredentialID != nil {
				if _, err = tx.Credential.Where(tx.Credential.ID.Eq(*oldConn.CredentialID)).Unscoped().Delete(); err != nil {
					return err
				}
			}
		}

		return tx.Connection.Where(tx.Connection.ID.Eq(conn.ID)).Save(conn)
	}); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *ConnectionSrv) FindByID(id uint) (*model.Connection, error) {
	t := s.Query.Connection
	return t.Where(t.ID.Eq(id)).Preload(t.Credential, t.Metadata).First()
}

func (s *ConnectionSrv) DeleteConnection(id uint) *resp.Resp {
	if err := s.Query.Transaction(func(tx *query.Query) error {
		conn, err := tx.Connection.Where(tx.Connection.ID.Eq(id)).First()
		if err != nil {
			return err
		}
		if !conn.UseCommonCredential && conn.CredentialID != nil {
			if _, err = tx.Credential.Where(tx.Credential.ID.Eq(*conn.CredentialID)).Unscoped().Delete(); err != nil {
				return err
			}
		}
		if _, err = tx.Connection.Select(tx.Connection.Metadata.Field()).Unscoped().Delete(conn); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.Ok()
}

func (s *ConnectionSrv) ListConnection() *resp.Resp {
	t := s.Query.Connection
	connList, err := t.Preload(t.Metadata, t.Credential).Find()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithData(connList)
}
