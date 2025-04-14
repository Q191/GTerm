package ssh

import (
	"errors"
	"fmt"
	"time"

	"github.com/MisakaTAT/GTerm/backend/types"

	"net"
	"os"
	"path/filepath"

	"github.com/MisakaTAT/GTerm/backend/enums"
	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/skeema/knownhosts"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	Host                string
	Port                uint
	User                string
	AuthMethod          enums.AuthMethod
	Password            string
	PrivateKey          string
	Passphrase          string
	TrustUnknownHost    bool
	Timeout             time.Duration
	Ciphers             []string
	KeyExchanges        []string
	MACs                []string
	HostKeyAlgorithms   []string
	PublicKeyAlgorithms []string
}

func NewSSHClient(c *Config, logger initialize.Logger) (*ssh.Client, error) {
	if c == nil {
		return nil, errors.New("config is not set")
	}
	if logger == nil {
		return nil, errors.New("logger is not set")
	}

	logger.Info("Connecting to SSH server %s:%d", c.Host, c.Port)
	host := fmt.Sprintf("%s:%d", c.Host, c.Port)

	var auth []ssh.AuthMethod
	switch c.AuthMethod {
	case enums.Password:
		auth = append(auth, ssh.Password(c.Password))
		auth = append(auth, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
			answers := make([]string, len(questions))
			if len(questions) == 1 {
				answers[0] = c.Password
			}
			return answers, nil
		}))
		logger.Info("Using password authentication")
	case enums.PrivateKey:
		var signer ssh.Signer
		var err error

		if c.Passphrase == "" {
			signer, err = ssh.ParsePrivateKey([]byte(c.PrivateKey))
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(c.PrivateKey), []byte(c.Passphrase))
		}

		if err != nil {
			logger.Error("Failed to parse private key: %v", err)
			return nil, err
		}

		auth = append(auth, ssh.PublicKeys(signer))
		logger.Info("Using private key authentication")
	default:
		return nil, errors.New("unsupported authentication method")
	}

	var hostKeyCallback ssh.HostKeyCallback
	var hostKeyAlgorithms []string

	if c.TrustUnknownHost {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
		logger.Warn("Configured to trust all unknown hosts, this may pose a security risk")
	} else {
		knownHostsFile := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
		logger.Debug("Using known_hosts file: %s", knownHostsFile)

		db, err := knownhosts.NewDB(knownHostsFile)
		if err != nil {
			logger.Error("Failed to load known_hosts database: %v", err)
			hostKeyCallback = ssh.InsecureIgnoreHostKey()
		} else {
			hostKeyCallback = db.HostKeyCallback()
			hostKeyAlgorithms = db.HostKeyAlgorithms(host)
			logger.Info("Using known_hosts for host key verification")
		}
	}

	timeout := 10 * time.Second
	if c.Timeout > 0 {
		timeout = c.Timeout
	}

	clientConfig := &ssh.ClientConfig{
		User:            c.User,
		Auth:            auth,
		Timeout:         timeout,
		HostKeyCallback: hostKeyCallback,
	}

	if len(hostKeyAlgorithms) > 0 {
		clientConfig.HostKeyAlgorithms = hostKeyAlgorithms
		logger.Info("Using host key algorithms from known_hosts: %v", hostKeyAlgorithms)
	} else if len(c.HostKeyAlgorithms) > 0 {
		clientConfig.HostKeyAlgorithms = c.HostKeyAlgorithms
		logger.Info("Using custom host key algorithms: %v", c.HostKeyAlgorithms)
	} else {
		clientConfig.HostKeyAlgorithms = []string{
			ssh.KeyAlgoED25519,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoRSASHA512,
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
		}
		logger.Info("Using default host key algorithm list")
	}

	if len(c.Ciphers) > 0 {
		clientConfig.Ciphers = c.Ciphers
		logger.Info("Using custom cipher list: %v", c.Ciphers)
	}

	if len(c.KeyExchanges) > 0 {
		clientConfig.KeyExchanges = c.KeyExchanges
		logger.Info("Using custom key exchange list: %v", c.KeyExchanges)
	}

	if len(c.MACs) > 0 {
		clientConfig.MACs = c.MACs
		logger.Info("Using custom MAC list: %v", c.MACs)
	}

	logger.Info("Starting SSH connection to server, %s@%s", c.User, host)

	client, err := ssh.Dial("tcp", host, clientConfig)
	if err != nil {
		if knownhosts.IsHostUnknown(err) {
			logger.Info("Unknown host, attempting to get host key: %s", host)
			key, keyErr := getHostKey(c, logger)
			if keyErr == nil && key != nil {
				fingerprint := ssh.FingerprintSHA256(key)
				logger.Info("Successfully obtained unknown host key, fingerprint: %s", fingerprint)
				return nil, &types.FingerprintError{
					Host:        host,
					Fingerprint: fingerprint,
				}
			}
		} else if knownhosts.IsHostKeyChanged(err) {
			logger.Warn("Host key has changed! This may indicate a MitM attack, host: %s", host)
		}
		logger.Error("SSH connection failed: %v", err)
		return nil, err
	}

	logger.Info("SSH connection successful, %s@%s", c.User, host)

	return client, nil
}

func getHostKey(c *Config, logger initialize.Logger) (hostKey ssh.PublicKey, err error) {
	host := fmt.Sprintf("%s:%d", c.Host, c.Port)

	timeout := 10 * time.Second
	if c.Timeout > 0 {
		timeout = c.Timeout
	}

	clientConfig := &ssh.ClientConfig{
		User:    c.User,
		Timeout: timeout,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			hostKey = key
			return nil
		},
	}

	conn, err := ssh.Dial("tcp", host, clientConfig)
	if err != nil {
		if hostKey != nil {
			logger.Info("Successfully obtained host key, fingerprint: %s", ssh.FingerprintSHA256(hostKey))
			return hostKey, nil
		}
		return nil, err
	}
	defer func(conn *ssh.Client) {
		if err = conn.Close(); err != nil {
			logger.Error("Failed to close SSH connection: %v", err)
		}
	}(conn)
	if hostKey == nil {
		return nil, errors.New("unable to obtain host key")
	}

	logger.Info("Successfully obtained host key, fingerprint: %s", ssh.FingerprintSHA256(hostKey))
	return hostKey, nil
}

func AddFingerprint(conf *Config, host, fingerprint string, logger initialize.Logger) error {
	logger.Info("Adding host fingerprint, host: %s, fingerprint: %s", host, fingerprint)
	knownHostsFile := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	key, err := getHostKey(conf, logger)
	if err != nil {
		return err
	}

	actualFingerprint := ssh.FingerprintSHA256(key)
	logger.Info("Comparing fingerprints, expected: %s, actual: %s", fingerprint, actualFingerprint)

	if actualFingerprint != fingerprint {
		return errors.New("fingerprint mismatch")
	}

	f, err := os.OpenFile(knownHostsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		logger.Error("Failed to open known_hosts file: %v", err)
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if err = knownhosts.WriteKnownHost(f, host, nil, key); err != nil {
		logger.Error("Failed to write to known_hosts file: %v", err)
		return err
	}
	logger.Info("Successfully added fingerprint to known_hosts, fingerprint: %s", fingerprint)
	return nil
}
