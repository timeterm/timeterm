package secrets

import (
	"fmt"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

type Manager struct {
	safe *VaultClient

	operatorName string
}

func (d *Manager) Init() error {
	return d.InitKeys()
}

func (d *Manager) newAccountKeys() (nkeys.KeyPair, error) {
	kp, err := nkeys.CreateAccount()
	if err != nil {
		return kp, fmt.Errorf("could not generate account keys: %w", err)
	}

	if err = d.safe.WriteAccountSeed(kp); err != nil {
		return kp, fmt.Errorf("could not write account seed: %w", err)
	}

	return kp, nil
}

func (d *Manager) newOperatorKeys() (nkeys.KeyPair, error) {
	kp, err := nkeys.CreateOperator()
	if err != nil {
		return kp, fmt.Errorf("could not generate operator keys: %w", err)
	}

	if err = d.safe.WriteOperatorSeed(kp); err != nil {
		return kp, fmt.Errorf("could not write operator seed: %w", err)
	}

	return kp, nil
}

func (d *Manager) newUserKeys() (nkeys.KeyPair, error) {
	kp, err := nkeys.CreateUser()
	if err != nil {
		return kp, fmt.Errorf("could not generate user keys: %w", err)
	}

	if err = d.safe.WriteUserSeed(kp); err != nil {
		return kp, fmt.Errorf("could not write user seed: %w", err)
	}

	return kp, nil
}

func (d *Manager) newOperator(systemAccountPubKey string) (string, error) {
	kp, err := d.newOperatorKeys()
	if err != nil {
		return "", fmt.Errorf("could not create operator keys: %w", err)
	}

	pk, err := kp.PublicKey()
	if err != nil {
		return "", fmt.Errorf("could not create operator public key: %w", err)
	}

	claims := jwt.NewOperatorClaims(pk)
	claims.Name = d.operatorName
	claims.Issuer = pk
	claims.IssuedAt = time.Now().Unix()
	claims.SystemAccount = systemAccountPubKey
	// TODO(rutgerbrf): set some more claims to the correct values.

	err = d.safe.WriteOperatorJWT(claims, pk)
	if err != nil {
		return pk, fmt.Errorf("could not write operator JWT: %w", err)
	}

	return pk, nil
}

func (d *Manager) newAccount(name, operatorPubKey string) (string, error) {
	kp, err := d.newAccountKeys()
	if err != nil {
		return "", fmt.Errorf("could not create account keys: %w", err)
	}

	pk, err := kp.PublicKey()
	if err != nil {
		return "", fmt.Errorf("could not create account public key: %w", err)
	}

	if err = d.newAccountWithPubKey(name, pk, operatorPubKey); err != nil {
		return "", err
	}

	return pk, nil
}

func (d *Manager) newAccountWithPubKey(name, pubKey, operatorPubKey string) error {
	claims := jwt.NewAccountClaims(pubKey)
	claims.Name = name
	claims.Issuer = operatorPubKey
	claims.IssuedAt = time.Now().Unix()
	claims.Limits.JetStreamLimits = jwt.JetStreamLimits{
		Consumer:      jwt.NoLimit,
		DiskStorage:   jwt.NoLimit,
		MemoryStorage: jwt.NoLimit,
		Streams:       jwt.NoLimit,
	}

	err := d.safe.WriteAccountJWT(claims, operatorPubKey)
	if err != nil {
		return fmt.Errorf("could not write account JWT: %w", err)
	}

	return nil
}

func (d *Manager) newUser(userName, accountName, accountPubKey string) (string, error) {
	kp, err := d.newUserKeys()
	if err != nil {
		return "", fmt.Errorf("could not create user keys: %w", err)
	}

	pk, err := kp.PublicKey()
	if err != nil {
		return "", fmt.Errorf("could not create user public key: %w", err)
	}

	claims := jwt.NewUserClaims(pk)
	claims.Name = userName
	claims.Issuer = accountPubKey
	claims.IssuedAt = time.Now().Unix()

	err = d.safe.WriteUserJWT(claims, accountName, accountPubKey)
	if err != nil {
		return pk, fmt.Errorf("could not write user JWT: %w", err)
	}

	return pk, nil
}

func (d *Manager) InitKeys() error {
	sakp, err := d.newAccountKeys()
	if err != nil {
		return fmt.Errorf("could not create system account keys: %w", err)
	}

	sapk, err := sakp.PublicKey()
	if err != nil {
		return fmt.Errorf("could not create system account public key: %w", err)
	}

	oppk, err := d.newOperator(sapk)
	if err != nil {
		return fmt.Errorf("could not create operator: %w", err)
	}

	oapk, err := d.newAccount(d.operatorName, oppk)
	if err != nil {
		return fmt.Errorf("could not create operator account: %w", err)
	}

	if _, err = d.newUser(d.operatorName, d.operatorName, oapk); err != nil {
		return fmt.Errorf("could not create operator user: %w", err)
	}

	if err = d.newAccountWithPubKey("SYS", sapk, oppk); err != nil {
		return fmt.Errorf("could not create system account: %w", err)
	}

	if _, err = d.newUser("SYS", "SYS", sapk); err != nil {
		return fmt.Errorf("could not create system user: %w", err)
	}

	return nil
}
