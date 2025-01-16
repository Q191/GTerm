package services

import (
	"github.com/MisakaTAT/GTerm/backend/dal/model"
	"github.com/MisakaTAT/GTerm/backend/dal/query"
	"github.com/MisakaTAT/GTerm/backend/enums"
	"github.com/MisakaTAT/GTerm/backend/pkg/exec"
	"github.com/MisakaTAT/GTerm/backend/pkg/metadata"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var MetadataSrvSet = wire.NewSet(wire.Struct(new(MetadataSrv), "*"))

type MetadataSrv struct {
	Logger *zap.Logger
	Query  *query.Query
}

func (s *MetadataSrv) UpdateByHost(host *model.Host) {
	t := s.Query.Metadata

	config := &exec.Config{
		Host:     host.Host,
		Port:     host.Port,
		User:     host.Credential.Username,
		AuthType: host.Credential.AuthType,
	}
	switch host.Credential.AuthType {
	case enums.CredentialAuthTypePassword:
		config.Password = host.Credential.Password
	case enums.CredentialAuthTypePrivateKey:
		config.PrivateKey = host.Credential.PrivateKey
		config.KeyPassword = host.Credential.KeyPassword
	}
	client, err := exec.NewExec(config)
	if err != nil {
		s.Logger.Error("failed to create ssh client", zap.Error(err))
		return
	}
	defer func() {
		_ = client.Close()
	}()

	meta, err := t.Where(t.HostID.Eq(host.ID)).FirstOrInit()
	if err != nil {
		s.Logger.Error("failed to get metadata", zap.Error(err))
		return
	}

	osRelease := metadata.NewMetadata(client).GetOSRelease()
	if osRelease != nil {
		meta.OS = osRelease.PrettyName
	}

	if err = t.Save(meta); err != nil {
		s.Logger.Error("failed to update metadata", zap.Error(err))
	}
}
