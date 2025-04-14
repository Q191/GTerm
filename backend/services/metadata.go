package services

import (
	"github.com/MisakaTAT/GTerm/backend/dal/model"
	"github.com/MisakaTAT/GTerm/backend/dal/query"
	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/MisakaTAT/GTerm/backend/pkg/exec"
	"github.com/MisakaTAT/GTerm/backend/pkg/metadata"
	commonssh "github.com/MisakaTAT/GTerm/backend/pkg/ssh"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var MetadataSrvSet = wire.NewSet(wire.Struct(new(MetadataSrv), "*"))

type MetadataSrv struct {
	Logger initialize.Logger
	Query  *query.Query
}

func (s *MetadataSrv) UpdateByConnection(conn *model.Connection) {
	t := s.Query.Metadata

	config := &commonssh.Config{
		Host:             conn.Host,
		Port:             conn.Port,
		User:             conn.Credential.Username,
		AuthMethod:       conn.Credential.AuthMethod,
		Password:         conn.Credential.Password,
		PrivateKey:       conn.Credential.PrivateKey,
		Passphrase:       conn.Credential.Passphrase,
		TrustUnknownHost: true,
	}
	client, err := exec.NewExec(config, s.Logger)
	if err != nil {
		s.Logger.Error("failed to create ssh client", zap.Error(err))
		return
	}
	defer func() {
		_ = client.Close()
	}()

	meta, err := t.Where(t.ConnectionID.Eq(conn.ID)).FirstOrInit()
	if err != nil {
		s.Logger.Error("failed to get metadata", zap.Error(err))
		return
	}

	metaInfo := metadata.NewMetadata(client).Parser()
	if metaInfo != nil {
		meta.Vendor = metaInfo.Vendor
		meta.Type = metaInfo.Type
	}

	if err = t.Save(meta); err != nil {
		s.Logger.Error("failed to update metadata", zap.Error(err))
	}
}
