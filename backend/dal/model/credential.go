package model

import (
	"github.com/MisakaTAT/GTerm/backend/enums"
	"github.com/MisakaTAT/GTerm/backend/utils/encrypt"
	"gorm.io/gorm"
)

type Credential struct {
	Common
	Label                string           `json:"label" gorm:"uniqueIndex;not null"`
	Username             string           `json:"username"`
	IsCommonCredential   bool             `json:"isCommonCredential"`
	AuthMethod           enums.AuthMethod `json:"authMethod"`
	Password             string           `json:"password" gorm:"-"`
	PasswordCiphertext   string
	PasswordSalt         string
	PrivateKey           string `json:"privateKey" gorm:"-"`
	PrivateKeyCiphertext string
	PrivateKeySalt       string
	Passphrase           string `json:"passphrase" gorm:"-"`
	PassphraseCiphertext string
	PassphraseSalt       string
}

func (c *Credential) TableName() string {
	return "credentials"
}

func (c *Credential) BeforeCreate(tx *gorm.DB) error {
	cred, err := encrypt.NewCredential()
	if err != nil {
		return err
	}
	fields := []*encrypt.CredentialField{
		{Plaintext: c.Password, Ciphertext: &c.PasswordCiphertext, Salt: &c.PasswordSalt},
		{Plaintext: c.PrivateKey, Ciphertext: &c.PrivateKeyCiphertext, Salt: &c.PrivateKeySalt},
		{Plaintext: c.Passphrase, Ciphertext: &c.PassphraseCiphertext, Salt: &c.PassphraseSalt},
	}
	for _, field := range fields {
		if field.Plaintext == "" {
			continue
		}
		encrypted, err := cred.EncryptPassword(field.Plaintext)
		if err != nil {
			return err
		}
		*field.Ciphertext = encrypted.Ciphertext
		*field.Salt = encrypted.Salt
	}
	return nil
}

func (c *Credential) AfterFind(tx *gorm.DB) error {
	cred, err := encrypt.NewCredential()
	if err != nil {
		return err
	}
	fields := []*encrypt.CredentialField{
		{Plaintext: c.Password, Ciphertext: &c.PasswordCiphertext, Salt: &c.PasswordSalt},
		{Plaintext: c.PrivateKey, Ciphertext: &c.PrivateKeyCiphertext, Salt: &c.PrivateKeySalt},
		{Plaintext: c.Passphrase, Ciphertext: &c.PassphraseCiphertext, Salt: &c.PassphraseSalt},
	}
	for i, field := range fields {
		if *field.Ciphertext == "" || *field.Salt == "" {
			continue
		}
		plaintext, err := cred.DecryptPassword(*field.Ciphertext, *field.Salt)
		if err != nil {
			return err
		}
		switch i {
		case 0:
			c.Password = plaintext
		case 1:
			c.PrivateKey = plaintext
		case 2:
			c.Passphrase = plaintext
		}
	}
	return nil
}
