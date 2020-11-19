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
	opkp, err := d.newOperatorKeys()
	if err != nil {
		return "", fmt.Errorf("could not create operator keys: %w", err)
	}

	oppk, err := opkp.PublicKey()
	if err != nil {
		return "", fmt.Errorf("could not create operator public key: %w", err)
	}

	oclaims := jwt.NewOperatorClaims(oppk)
	oclaims.Name = d.operatorName
	oclaims.Issuer = oppk
	oclaims.IssuedAt = time.Now().Unix()
	oclaims.SystemAccount = systemAccountPubKey

	err = d.safe.WriteOperatorJWT(oclaims, oppk)
	if err != nil {
		return oppk, fmt.Errorf("could not write operator JWT: %w", err)
	}

	return oppk, nil
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

	oakp, err := d.newAccountKeys()
	if err != nil {
		return fmt.Errorf("could not create operator account keys: %w", err)
	}

	oapk, err := oakp.PublicKey()
	if err != nil {
		return fmt.Errorf("could not create operator account public key: %w", err)
	}

	oaclaims := jwt.NewAccountClaims(oapk)
	oaclaims.Name = d.operatorName
	oaclaims.Issuer = oppk
	oaclaims.IssuedAt = time.Now().Unix()

	err = d.safe.WriteAccountJWT(oaclaims, oppk)
	if err != nil {
		return fmt.Errorf("could not write operator account JWT: %w", err)
	}

	oukp, err := d.newUserKeys()
	if err != nil {
		return fmt.Errorf("could not create operator user keys: %w", err)
	}

	oupk, err := oukp.PublicKey()
	if err != nil {
		return fmt.Errorf("could not create operator user public key: %w", err)
	}

	ouclaims := jwt.NewUserClaims(oupk)
	ouclaims.Name = d.operatorName
	ouclaims.Issuer = oapk
	ouclaims.IssuedAt = time.Now().Unix()

	err = d.safe.WriteUserJWT(ouclaims, d.operatorName, oapk)
	if err != nil {
		return fmt.Errorf("could not write operator user JWT: %w", err)
	}

	saclaims := jwt.NewAccountClaims(sapk)
	saclaims.Name = "SYS"
	saclaims.Issuer = oppk
	saclaims.IssuedAt = time.Now().Unix()

	err = d.safe.WriteAccountJWT(saclaims, oppk)
	if err != nil {
		return fmt.Errorf("could not write system account JWT: %w", err)
	}

	sukp, err := d.newUserKeys()
	if err != nil {
		return fmt.Errorf("could not create system user keys: %w", err)
	}

	supk, err := sukp.PublicKey()
	if err != nil {
		return fmt.Errorf("could not create system user public key: %w", err)
	}

	suclaims := jwt.NewUserClaims(supk)
	suclaims.Name = "SYS"
	suclaims.Issuer = sapk
	suclaims.IssuedAt = time.Now().Unix()

	err = d.safe.WriteUserJWT(suclaims, saclaims.Name, sapk)
	if err != nil {
		return fmt.Errorf("could not write system user JWT: %w", err)
	}

	return nil
}
