package secrets

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"

	"gitlab.com/timeterm/timeterm/nats-manager/database"
)

type Manager struct {
	safe *VaultClient
	dbw  *database.Wrapper

	operatorName string
}

func NewManager(c *VaultClient, dbw *database.Wrapper, operatorName string) *Manager {
	return &Manager{
		safe:         c,
		dbw:          dbw,
		operatorName: operatorName,
	}
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

	if err = m.dbw.CreateOperator(ctx, m.operatorName, pk); err != nil {
		return "", fmt.Errorf("could not create operator in database: %w", err)
	}

	claims := jwt.NewOperatorClaims(pk)
	claims.Name = m.operatorName
	claims.Issuer = pk
	claims.IssuedAt = time.Now().Unix()
	claims.SystemAccount = systemAccountPubKey
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

func (m *Manager) InitKeys(ctx context.Context) error {
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

	oapk, err := m.newAccount(ctx, m.operatorName, oppk)
	if err != nil {
		return fmt.Errorf("could not create operator account: %w", err)
	}

	if _, err = m.newUser(ctx, m.operatorName, oapk); err != nil {
		return fmt.Errorf("could not create operator user: %w", err)
	}

	if err = m.newAccountWithPubKey(ctx, "SYS", sapk, oppk); err != nil {
		return fmt.Errorf("could not create system account: %w", err)
	}

	if _, err = m.newUser(ctx, "sys", sapk); err != nil {
		return fmt.Errorf("could not create system user: %w", err)
	}

	return nil
}
