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

func (d *Manager) InitKeys() error {
	sakp, err := nkeys.CreateAccount()
	if err != nil {
		return fmt.Errorf("could not create system account keys: %w", err)
	}

	if err = d.safe.WriteAccountSeed(sakp); err != nil {
		return fmt.Errorf("could not write system account seed: %w", err)
	}

	sapk, err := sakp.PublicKey()
	if err != nil {
		return fmt.Errorf("could not create system account public key: %w", err)
	}

	opkp, err := nkeys.CreateOperator()
	if err != nil {
		return fmt.Errorf("could not create operator keys: %w", err)
	}

	if err = d.safe.WriteOperatorSeed(opkp); err != nil {
		return fmt.Errorf("could not write operator seed: %w", err)
	}

	oppk, err := opkp.PublicKey()
	if err != nil {
		return fmt.Errorf("could not create operator public key: %w", err)
	}

	oclaims := jwt.NewOperatorClaims(oppk)
	oclaims.Name = d.operatorName
	oclaims.Issuer = oppk
	oclaims.IssuedAt = time.Now().Unix()
	oclaims.SystemAccount = sapk

	err = d.safe.WriteOperatorJWT(oclaims, oppk)
	if err != nil {
		return fmt.Errorf("could not write operator JWT: %w", err)
	}

	oakp, err := nkeys.CreateAccount()
	if err != nil {
		return fmt.Errorf("could not create operator account keys: %w", err)
	}

	if err = d.safe.WriteAccountSeed(oakp); err != nil {
		return fmt.Errorf("could not write operator account seed: %w", err)
	}

	oapk, err := oakp.PublicKey()
	if err != nil {
		return fmt.Errorf("could not create operator account public key: %w", err)
	}

	oukp, err := nkeys.CreateUser()
	if err != nil {
		return fmt.Errorf("could not create operator user keys: %w", err)
	}

	if err = d.safe.WriteUserSeed(oukp); err != nil {
		return fmt.Errorf("could not write operator user seed: %w", err)
	}

	oupk, err := oukp.PublicKey()
	if err != nil {
		return fmt.Errorf("could not create operator user public key: %w", err)
	}

	oaclaims := jwt.NewAccountClaims(oapk)
	oaclaims.Name = d.operatorName
	oaclaims.Issuer = oppk
	oaclaims.IssuedAt = time.Now().Unix()

	err = d.safe.WriteAccountJWT(oaclaims, oppk)
	if err != nil {
		return fmt.Errorf("could not write operator account JWT: %w", err)
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

	sukp, err := nkeys.CreateUser()
	if err != nil {
		return fmt.Errorf("could not create system user keys: %w", err)
	}

	if err = d.safe.WriteUserSeed(sukp); err != nil {
		return fmt.Errorf("could not write system user seed: %w", err)
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
