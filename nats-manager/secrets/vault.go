package secrets

import (
	"errors"
	"fmt"
	"path"

	vault "github.com/hashicorp/vault/api"
	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

type Store struct {
	prefix string
	vault  *vault.Client
}

func NewStore(prefix string, c *vault.Client) *Store {
	return &Store{
		prefix: prefix,
		vault:  c,
	}
}

func (s *Store) operatorSeedPath(pubKey string) string {
	return path.Join(s.prefix, "/keys/operator/", pubKey)
}

func (s *Store) accountSeedPath(pubKey string) string {
	return path.Join(s.prefix, "/keys/account/", pubKey)
}

func (s *Store) userSeedPath(pubKey string) string {
	return path.Join(s.prefix, "/keys/user/", pubKey)
}

func (s *Store) jwtPath(subject string) string {
	return path.Join(s.prefix, "/jwts/", subject)
}

func (s *Store) WriteOperatorSeed(kp nkeys.KeyPair) error {
	pubKey, err := kp.PublicKey()
	if err != nil {
		return err
	}

	pat := s.operatorSeedPath(pubKey)
	err = s.writeSeed(pat, kp)
	if err != nil {
		return fmt.Errorf("could not write operator seed: %w", err)
	}
	return nil
}

func (s *Store) WriteAccountSeed(kp nkeys.KeyPair) error {
	pubKey, err := kp.PublicKey()
	if err != nil {
		return err
	}

	pat := s.accountSeedPath(pubKey)
	err = s.writeSeed(pat, kp)
	if err != nil {
		return fmt.Errorf("could not write account seed: %w", err)
	}
	return nil
}

func (s *Store) WriteUserSeed(kp nkeys.KeyPair) error {
	pubKey, err := kp.PublicKey()
	if err != nil {
		return err
	}

	pat := s.userSeedPath(pubKey)
	err = s.writeSeed(pat, kp)
	if err != nil {
		return fmt.Errorf("could not write user seed: %w", err)
	}
	return nil
}

func (s *Store) WriteOperatorJWT(claims *jwt.OperatorClaims, operatorPubKey string) error {
	if claims.Name == "" {
		return errors.New("operator name may not be empty")
	}
	kp, err := s.ReadOperatorSeed(operatorPubKey)
	if err != nil {
		return fmt.Errorf("could not read operator seed: %w", err)
	}

	pat := s.jwtPath(claims.Subject)
	return s.writeJWT(pat, claims, kp)
}

func (s *Store) WriteUserJWT(claims *jwt.UserClaims, accountPubKey string) error {
	if claims.Name == "" {
		return errors.New("user name may not be empty")
	}
	kp, err := s.ReadAccountSeed(accountPubKey)
	if err != nil {
		return fmt.Errorf("could not read account seed: %w", err)
	}

	pat := s.jwtPath(claims.Subject)
	return s.writeJWT(pat, claims, kp)
}

func (s *Store) WriteAccountJWT(claims *jwt.AccountClaims, operatorPubKey string) error {
	if claims.Name == "" {
		return errors.New("account name may not be empty")
	}
	kp, err := s.ReadOperatorSeed(operatorPubKey)
	if err != nil {
		return fmt.Errorf("could not read operator seed: %w", err)
	}

	pat := s.jwtPath(claims.Subject)
	return s.writeJWT(pat, claims, kp)
}

func (s *Store) writeJWT(pat string, claims jwt.Claims, kp nkeys.KeyPair) error {
	encoded, err := claims.Encode(kp)
	if err != nil {
		return fmt.Errorf("could not encode operator claims: %w", err)
	}

	_, err = s.vault.Logical().Write(pat, map[string]interface{}{
		"jwt": encoded,
	})
	return err
}

func (s *Store) readJWT(pat string) (string, error) {
	secret, err := s.vault.Logical().Read(pat)
	if err != nil {
		return "", err
	}

	token, ok := secret.Data["jwt"].(string)
	if !ok {
		return "", fmt.Errorf("jwt not present in secret at path %s", pat)
	}
	return token, nil
}

func (s *Store) ReadJWT(subject string) (jwt.Claims, error) {
	token, err := s.ReadJWTLiteral(subject)
	if err != nil {
		return nil, err
	}

	claims, err := jwt.Decode(token)
	if err != nil {
		return nil, fmt.Errorf("could not decode claims: %w", err)
	}
	return claims, nil
}

func (s *Store) ReadJWTLiteral(subject string) (string, error) {
	pat := s.jwtPath(subject)
	token, err := s.readJWT(pat)
	if err != nil {
		return "", fmt.Errorf("could not read JWT at path %s from Vault: %w", pat, err)
	}
	return token, nil
}

func (s *Store) ReadOperatorJWT(subject string) (*jwt.OperatorClaims, error) {
	token, err := s.ReadJWTLiteral(subject)
	if err != nil {
		return nil, err
	}

	claims, err := jwt.DecodeOperatorClaims(token)
	if err != nil {
		return nil, fmt.Errorf("could not decode operator claims: %w", err)
	}
	return claims, nil
}

func (s *Store) ReadAccountJWT(subject string) (*jwt.AccountClaims, error) {
	pat := s.jwtPath(subject)
	token, err := s.readJWT(pat)
	if err != nil {
		return nil, err
	}

	claims, err := jwt.DecodeAccountClaims(token)
	if err != nil {
		return nil, fmt.Errorf("could not decode account claims: %w", err)
	}
	return claims, nil
}

func (s *Store) ReadUserJWT(subject string) (*jwt.UserClaims, error) {
	pat := s.jwtPath(subject)
	token, err := s.readJWT(pat)
	if err != nil {
		return nil, err
	}

	claims, err := jwt.DecodeUserClaims(token)
	if err != nil {
		return nil, fmt.Errorf("could not decode user claims: %w", err)
	}
	return claims, nil
}

func (s *Store) ReadOperatorSeed(pubKey string) (nkeys.KeyPair, error) {
	pat := s.operatorSeedPath(pubKey)
	return s.readSeed(pat)
}

func (s *Store) ReadUserSeed(pubKey string) (nkeys.KeyPair, error) {
	pat := s.userSeedPath(pubKey)
	return s.readSeed(pat)
}

func (s *Store) ReadAccountSeed(pubKey string) (nkeys.KeyPair, error) {
	pat := s.accountSeedPath(pubKey)
	return s.readSeed(pat)
}

func (s *Store) writeSeed(path string, kp nkeys.KeyPair) error {
	seed, err := kp.Seed()
	if err != nil {
		return err
	}

	// Converting seed to a string is safe (or should really be safe)
	// because it should be in Base 32.
	seedStr := string(seed)
	_, err = s.vault.Logical().Write(path, map[string]interface{}{
		"seed": seedStr,
	})
	return err
}

func (s *Store) readSeed(path string) (nkeys.KeyPair, error) {
	var kp nkeys.KeyPair

	secret, err := s.vault.Logical().Read(path)
	if err != nil {
		return kp, err
	}

	seed, ok := secret.Data["seed"].(string)
	if !ok {
		return kp, fmt.Errorf("seed not present in secret at path %s", path)
	}

	kp, err = nkeys.FromSeed([]byte(seed))
	if err != nil {
		return kp, fmt.Errorf("could not create key pair from seed: %w", err)
	}
	return kp, nil
}
