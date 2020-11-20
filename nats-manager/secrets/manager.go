package secrets

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"

	"gitlab.com/timeterm/timeterm/nats-manager/database"
)

type Manager struct {
	safe *VaultClient
	log  logr.Logger
	dbw  *database.Wrapper

	operator OperatorConfig
}

type OperatorConfig struct {
	Name        string
	AccountName string
	UserName    string
	ServiceURLS []string
}

func DefaultOperatorConfig() OperatorConfig {
	return OperatorConfig{
		Name:        os.Getenv("OPERATOR_NAME"),
		AccountName: os.Getenv("OPERATOR_ACCOUNT_NAME"),
		UserName:    os.Getenv("OPERATOR_USER_NAME"),
		ServiceURLS: strings.Split(os.Getenv("OPERATOR_SERVICE_URLS"), ","),
	}
}

func (c OperatorConfig) Validate() error {
	var err error

	if c.Name == "" {
		err = multierror.Append(err, errors.New("operator name is not set"))
	}
	if c.AccountName == "" {
		err = multierror.Append(err, errors.New("operator account name is not set"))
	}
	if c.UserName == "" {
		err = multierror.Append(err, errors.New("operator user name is not set"))
	}

	return err
}

func NewManager(log logr.Logger, c *VaultClient, dbw *database.Wrapper, oc OperatorConfig) (*Manager, error) {
	if err := oc.Validate(); err != nil {
		return nil, fmt.Errorf("error validating operator config: %w", err)
	}

	return &Manager{
		log:      log,
		safe:     c,
		dbw:      dbw,
		operator: oc,
	}, nil
}

func (m *Manager) Init(ctx context.Context) error {
	return m.InitKeys(ctx)
}

func (m *Manager) newAccountKeys() (nkeys.KeyPair, error) {
	kp, err := nkeys.CreateAccount()
	if err != nil {
		return kp, fmt.Errorf("could not generate account keys: %w", err)
	}

	if err = m.safe.WriteAccountSeed(kp); err != nil {
		return kp, fmt.Errorf("could not write account seed: %w", err)
	}

	return kp, nil
}

func (m *Manager) newOperatorKeys() (nkeys.KeyPair, error) {
	kp, err := nkeys.CreateOperator()
	if err != nil {
		return kp, fmt.Errorf("could not generate operator keys: %w", err)
	}

	if err = m.safe.WriteOperatorSeed(kp); err != nil {
		return kp, fmt.Errorf("could not write operator seed: %w", err)
	}

	return kp, nil
}

func (m *Manager) newUserKeys() (nkeys.KeyPair, error) {
	kp, err := nkeys.CreateUser()
	if err != nil {
		return kp, fmt.Errorf("could not generate user keys: %w", err)
	}

	if err = m.safe.WriteUserSeed(kp); err != nil {
		return kp, fmt.Errorf("could not write user seed: %w", err)
	}

	return kp, nil
}

func (m *Manager) newOperator(ctx context.Context, systemAccountPubKey string) (string, error) {
	kp, err := m.newOperatorKeys()
	if err != nil {
		return "", fmt.Errorf("could not create operator keys: %w", err)
	}

	pk, err := kp.PublicKey()
	if err != nil {
		return "", fmt.Errorf("could not create operator public key: %w", err)
	}

	if err = m.dbw.CreateOperator(ctx, m.operator.Name, pk); err != nil {
		return "", fmt.Errorf("could not create operator in database: %w", err)
	}

	claims := jwt.NewOperatorClaims(pk)
	claims.Name = m.operator.Name
	claims.Issuer = pk
	claims.IssuedAt = time.Now().Unix()
	claims.SystemAccount = systemAccountPubKey
	claims.OperatorServiceURLs = m.operator.ServiceURLS
	// TODO(rutgerbrf): set some more claims to the correct values.

	err = m.safe.WriteOperatorJWT(claims, pk)
	if err != nil {
		return pk, fmt.Errorf("could not write operator JWT: %w", err)
	}

	return pk, nil
}

func (m *Manager) newAccount(ctx context.Context, name, operatorPubKey string) (string, error) {
	kp, err := m.newAccountKeys()
	if err != nil {
		return "", fmt.Errorf("could not create account keys: %w", err)
	}

	pk, err := kp.PublicKey()
	if err != nil {
		return "", fmt.Errorf("could not create account public key: %w", err)
	}

	if err = m.newAccountWithPubKey(ctx, name, pk, operatorPubKey); err != nil {
		return "", err
	}

	return pk, nil
}

func (m *Manager) NewAccount(ctx context.Context, name string) (string, error) {
	pk, err := m.dbw.GetOperatorSubject(ctx, m.operator.Name)
	if err != nil {
		return "", fmt.Errorf("could not fetch operator public key: %w", err)
	}
	return m.newAccount(ctx, name, pk)
}

func (m *Manager) newAccountWithPubKey(ctx context.Context, name, pubKey, operatorPubKey string) error {
	if err := m.dbw.CreateAccount(ctx, name, pubKey, operatorPubKey); err != nil {
		return fmt.Errorf("could not create account in database: %w", err)
	}

	claims := jwt.NewAccountClaims(pubKey)
	claims.Name = name
	claims.Issuer = operatorPubKey
	claims.IssuedAt = time.Now().Unix()
	// TODO(rutgerbrf): set some more claims to the correct values.

	err := m.safe.WriteAccountJWT(claims, operatorPubKey)
	if err != nil {
		return fmt.Errorf("could not write account JWT: %w", err)
	}

	return nil
}

func (m *Manager) newSystemAccount(ctx context.Context, name, pubKey, operatorPubKey string) error {
	if err := m.dbw.CreateAccount(ctx, name, pubKey, operatorPubKey); err != nil {
		return fmt.Errorf("could not create account in database: %w", err)
	}

	claims := jwt.NewAccountClaims(pubKey)
	claims.Name = name
	claims.Issuer = operatorPubKey
	claims.IssuedAt = time.Now().Unix()
	// System accounts can not have JetStream enabled
	claims.Limits.JetStreamLimits = jwt.JetStreamLimits{}
	// TODO(rutgerbrf): set some more claims to the correct values.

	err := m.safe.WriteAccountJWT(claims, operatorPubKey)
	if err != nil {
		return fmt.Errorf("could not write account JWT: %w", err)
	}

	return nil
}

func (m *Manager) newUser(ctx context.Context, userName, accountPubKey string) (string, error) {
	kp, err := m.newUserKeys()
	if err != nil {
		return "", fmt.Errorf("could not create user keys: %w", err)
	}

	pk, err := kp.PublicKey()
	if err != nil {
		return "", fmt.Errorf("could not create user public key: %w", err)
	}

	if err = m.dbw.CreateUser(ctx, userName, pk, accountPubKey); err != nil {
		return "", fmt.Errorf("could not create user in database: %w", err)
	}

	claims := jwt.NewUserClaims(pk)
	claims.Name = userName
	claims.Issuer = accountPubKey
	claims.IssuedAt = time.Now().Unix()
	// TODO(rutgerbrf): set some more claims to the correct values.

	err = m.safe.WriteUserJWT(claims, accountPubKey)
	if err != nil {
		return pk, fmt.Errorf("could not write user JWT: %w", err)
	}

	return pk, nil
}

func (m *Manager) NewUser(ctx context.Context, name, accountName string) (string, error) {
	pk, err := m.dbw.GetAccountSubject(ctx, accountName, m.operator.Name)
	if err != nil {
		return "", fmt.Errorf("could not fetch account public key: %w", err)
	}
	return m.newUser(ctx, name, pk)
}

func (m *Manager) InitKeys(ctx context.Context) error {
	m.log.Info("Initializing keys")

	sakp, err := m.newAccountKeys()
	if err != nil {
		return fmt.Errorf("could not create system account keys: %w", err)
	}

	sapk, err := sakp.PublicKey()
	if err != nil {
		return fmt.Errorf("could not create system account public key: %w", err)
	}

	oppk, err := m.newOperator(ctx, sapk)
	if err != nil {
		return fmt.Errorf("could not create operator: %w", err)
	}

	oapk, err := m.newAccount(ctx, m.operator.AccountName, oppk)
	if err != nil {
		return fmt.Errorf("could not create operator account: %w", err)
	}

	if _, err = m.newUser(ctx, m.operator.UserName, oapk); err != nil {
		return fmt.Errorf("could not create operator user: %w", err)
	}

	if err = m.newSystemAccount(ctx, "SYS", sapk, oppk); err != nil {
		return fmt.Errorf("could not create system account: %w", err)
	}

	if _, err = m.newUser(ctx, "sys", sapk); err != nil {
		return fmt.Errorf("could not create system user: %w", err)
	}

	optok, err := m.safe.ReadJWTLiteral(oppk)
	if err != nil {
		return fmt.Errorf("could not read operator JWT: %w", err)
	}

	m.log.Info("keys initialized")
	m.log.Info("system account public key", "pubKey", sapk)
	m.log.Info("operator JWT", "jwt", optok)

	return nil
}

func (m *Manager) getOperatorPubKey(ctx context.Context) (string, error) {
	pk, err := m.dbw.GetOperatorSubject(ctx, m.operator.Name)
	if err != nil {
		return "", fmt.Errorf("could not get operator subject: %w", err)
	}
	return pk, nil
}

func (m *Manager) CreateNewDeviceUser(ctx context.Context, id uuid.UUID) error {
	oppk, err := m.getOperatorPubKey(ctx)
	if err != nil {
		return err
	}

	dapk, err := m.newAccount(ctx, deviceAccountName(id), oppk)
	if err != nil {
		return err
	}

	_, err = m.newUser(ctx, deviceUserName(id), dapk)
	return err
}

func (m *Manager) GenerateDeviceCredentials(ctx context.Context, id uuid.UUID) (string, error) {
	return m.GenerateUserCredentials(ctx, deviceUserName(id), deviceAccountName(id))
}

func (m *Manager) GenerateUserCredentials(ctx context.Context, userName, accountName string) (string, error) {
	pk, err := m.dbw.GetUserSubject(ctx, userName, accountName, m.operator.Name)
	if err != nil {
		return "", err
	}

	kp, err := m.safe.ReadUserSeed(pk)
	if err != nil {
		return "", err
	}

	seed, err := kp.Seed()
	if err != nil {
		return "", err
	}

	token, err := m.safe.ReadJWTLiteral(pk)
	if err != nil {
		return "", err
	}

	// Should be valid
	if _, err = jwt.DecodeUserClaims(token); err != nil {
		return "", err
	}

	cfg, err := jwt.FormatUserConfig(token, seed)
	if err != nil {
		return "", fmt.Errorf("could not format user config: %w", err)
	}
	return string(cfg), nil
}

func deviceAccountName(id uuid.UUID) string {
	return fmt.Sprintf("EMDEV-%s", id)
}

func deviceUserName(id uuid.UUID) string {
	return fmt.Sprintf("emdev-%s", id)
}
