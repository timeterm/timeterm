package manager

import (
	"context"
	"database/sql"
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
	"gitlab.com/timeterm/timeterm/nats-manager/secrets"
)

type (
	// An AccountClaimsEditor edits account claims.
	AccountClaimsEditor func(c *jwt.AccountClaims)
	// A UserClaimsEditor edits user claims.
	UserClaimsEditor func(c *jwt.UserClaims)
	// A OperatorClaimsEditor edits operator claims.
	OperatorClaimsEditor func(c *jwt.OperatorClaims)
)

// Manager manages operators, accounts, users and the tokens of these entities.
type Manager struct {
	secrets *secrets.Store
	log     logr.Logger
	dbw     *database.Wrapper

	operator OperatorConfig
}

// OperatorConfig contains configuration about the operator
// used by the manager.
type OperatorConfig struct {
	Name        string
	AccountName string
	UserName    string
	ServiceURLs []string
}

// DefaultOperatorConfig loads the default (unvalidated) OperatorConfig
// from the environment.
func DefaultOperatorConfig() OperatorConfig {
	return OperatorConfig{
		Name:        os.Getenv("OPERATOR_NAME"),
		AccountName: os.Getenv("OPERATOR_ACCOUNT_NAME"),
		UserName:    os.Getenv("OPERATOR_USER_NAME"),
		ServiceURLs: strings.Split(os.Getenv("OPERATOR_SERVICE_URLS"), ","),
	}
}

// Validate validates the OperatorConfig.
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

// New creates a new Manager. All parameters must be non-nil and oc must be valid.
func New(log logr.Logger, store *secrets.Store, dbw *database.Wrapper, oc OperatorConfig) (*Manager, error) {
	if err := oc.Validate(); err != nil {
		return nil, fmt.Errorf("error validating operator config: %w", err)
	}

	return &Manager{
		log:      log.WithName("Manager"),
		secrets:  store,
		dbw:      dbw,
		operator: oc,
	}, nil
}

// Init initializes the manager. Only has to run on the first run of the program (ever), as it configures
// the keys necessary for issuing other accounts and users.
func (m *Manager) Init(ctx context.Context) error {
	return m.InitKeys(ctx)
}

// CheckJWTs checks all JWTs for validity and writes to the log for those that have issues.
func (m *Manager) CheckJWTs(ctx context.Context) {
	m.log.Info("checking JWTs")

	if err := m.dbw.WalkJWTs(ctx, func(subj string) bool {
		if err := m.checkJWT(subj); err != nil {
			m.log.Error(err, "error checking JWT", "subject", subj)
		}
		return true
	}); err != nil {
		m.log.Error(err, "error walking JWTs")
	}

	m.log.Info("JWTs checked")
}

// checkJWT checks a single JWT and returns an error if validation fails.
// A JWT with validation errors does not count as a failure, but the issues with the JWT are logged.
func (m *Manager) checkJWT(subj string) error {
	token, err := m.secrets.ReadJWT(subj)
	if err != nil {
		return fmt.Errorf("could not read JWT: %w", err)
	}

	var vr jwt.ValidationResults
	token.Validate(&vr)

	for _, res := range vr.Issues {
		// Don't count expiration as an error.
		if !res.TimeCheck {
			if res.Blocking {
				m.log.Error(res, "JWT has a blocking validation issue (error)", "subject", subj)
			} else {
				m.log.Info("JWT has a validation issue", "subject", subj, "warning", res.Description)
			}
		}
	}

	return nil
}

// newAccountKeys creates a new key pair for an account and writes it to Vault.
func (m *Manager) newAccountKeys() (nkeys.KeyPair, error) {
	kp, err := nkeys.CreateAccount()
	if err != nil {
		return kp, fmt.Errorf("could not generate account keys: %w", err)
	}

	if err = m.secrets.WriteAccountSeed(kp); err != nil {
		return kp, fmt.Errorf("could not write account seed: %w", err)
	}

	return kp, nil
}

// newOperatorKeys creates a new key pair for an operator and writes it to Vault.
func (m *Manager) newOperatorKeys() (nkeys.KeyPair, error) {
	kp, err := nkeys.CreateOperator()
	if err != nil {
		return kp, fmt.Errorf("could not generate operator keys: %w", err)
	}

	if err = m.secrets.WriteOperatorSeed(kp); err != nil {
		return kp, fmt.Errorf("could not write operator seed: %w", err)
	}

	return kp, nil
}

// newUserKeys creates a new key pair for a user and writes it to Vault.
func (m *Manager) newUserKeys() (nkeys.KeyPair, error) {
	kp, err := nkeys.CreateUser()
	if err != nil {
		return kp, fmt.Errorf("could not generate user keys: %w", err)
	}

	if err = m.secrets.WriteUserSeed(kp); err != nil {
		return kp, fmt.Errorf("could not write user seed: %w", err)
	}

	return kp, nil
}

// NewOperator creates a new operator. It requires the public key of the system account.
// Additional changes to the claims can be made using editors. It automatically creates
// the key pair and registers the required information in the database and Vault.
// It returns the public key (subject) of the operator.
// The name is decided by the Manager's configuration and defaults to the value
// provided by the environment variable OPERATOR_NAME.
func (m *Manager) NewOperator(
	ctx context.Context,
	systemAccountPubKey string,
	editors ...OperatorClaimsEditor,
) (string, error) {
	kp, err := m.newOperatorKeys()
	if err != nil {
		return "", fmt.Errorf("could not create operator keys: %w", err)
	}
	defer kp.Wipe()

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
	claims.OperatorServiceURLs = m.operator.ServiceURLs
	for _, edit := range editors {
		edit(claims)
	}

	err = m.secrets.WriteOperatorJWT(claims, pk)
	if err != nil {
		return pk, fmt.Errorf("could not write operator JWT: %w", err)
	}

	return pk, nil
}

// newAccount creates a new account, automatically creating a key pair.
// It requires the public key of the operator and a name.
// newAccount automatically registers the required information in the database and Vault.
// It returns the public key (subject) of the newly created account.
func (m *Manager) newAccount(
	ctx context.Context,
	name, operatorPubKey string,
	editors ...AccountClaimsEditor,
) (string, error) {
	kp, err := m.newAccountKeys()
	if err != nil {
		return "", fmt.Errorf("could not create account keys: %w", err)
	}
	defer kp.Wipe()

	pk, err := kp.PublicKey()
	if err != nil {
		return "", fmt.Errorf("could not create account public key: %w", err)
	}

	if err = m.newAccountWithPubKey(ctx, name, pk, operatorPubKey, editors...); err != nil {
		return "", err
	}

	return pk, nil
}

// NewAccount creates a new account, automatically creating a key pair.
// The name for the account must be provided. The required information for validation
// is automatically registered in the database and Vault.
// It returns the public key (subject) of the newly created account.
func (m *Manager) NewAccount(ctx context.Context, name string, editors ...AccountClaimsEditor) (string, error) {
	pk, err := m.dbw.GetOperatorSubject(ctx, m.operator.Name)
	if err != nil {
		return "", fmt.Errorf("could not fetch operator public key: %w", err)
	}
	return m.newAccount(ctx, name, pk, editors...)
}

// newAccountWithPubKey creates a new account with a known name, public key
// and public key for the operator.
func (m *Manager) newAccountWithPubKey(
	ctx context.Context,
	name, pubKey, operatorPubKey string,
	editors ...AccountClaimsEditor,
) error {
	if err := m.dbw.CreateAccount(ctx, name, pubKey, operatorPubKey); err != nil {
		return fmt.Errorf("could not create account in database: %w", err)
	}

	claims := jwt.NewAccountClaims(pubKey)
	claims.Name = name
	claims.Issuer = operatorPubKey
	claims.IssuedAt = time.Now().Unix()
	for _, edit := range editors {
		edit(claims)
	}

	err := m.secrets.WriteAccountJWT(claims, operatorPubKey)
	if err != nil {
		return fmt.Errorf("could not write account JWT: %w", err)
	}

	return nil
}

// newSystemAccount creates a new system account.
// It does the same as newAccountWithPubKey but disables JetStream.
func (m *Manager) newSystemAccount(
	ctx context.Context,
	name, pubKey, operatorPubKey string,
	editors ...AccountClaimsEditor,
) error {
	return m.newAccountWithPubKey(
		ctx, name, pubKey, operatorPubKey,
		append(
			[]AccountClaimsEditor{
				func(c *jwt.AccountClaims) {
					// A system account cannot have JetStream configured.
					c.Limits.JetStreamLimits = jwt.JetStreamLimits{}
				},
			},
			editors...,
		)...,
	)
}

// newUser creates a new user, issued by an account with a known public key (accountPubKey).
// The default claims can be edited with editors.
func (m *Manager) newUser(
	ctx context.Context,
	userName, accountPubKey string,
	editors ...UserClaimsEditor,
) (string, error) {
	kp, err := m.newUserKeys()
	if err != nil {
		return "", fmt.Errorf("could not create user keys: %w", err)
	}
	defer kp.Wipe()

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
	// Always allow listening for responses
	claims.Sub.Allow = []string{"INBOX.>"}
	for _, edit := range editors {
		edit(claims)
	}

	err = m.secrets.WriteUserJWT(claims, accountPubKey)
	if err != nil {
		return pk, fmt.Errorf("could not write user JWT: %w", err)
	}

	return pk, nil
}

// NewUser creates a new user issued by an existing account. The default claims can be edited with editors.
func (m *Manager) NewUser(ctx context.Context, name, accountName string, editors ...UserClaimsEditor) (string, error) {
	acpk, err := m.dbw.GetAccountSubject(ctx, accountName, m.operator.Name)
	if err != nil {
		return "", fmt.Errorf("could not fetch account public key: %w", err)
	}

	pk, err := m.newUser(ctx, name, acpk, editors...)
	if err != nil {
		return "", err
	}
	return pk, m.SaveAppCreds(ctx, name, accountName)
}

func (m *Manager) UpdateUser(ctx context.Context, name, accountName string, editors ...UserClaimsEditor) error {
	pk, err := m.dbw.GetUserSubject(ctx, name, accountName, m.operator.Name)
	if err != nil {
		return fmt.Errorf("could not get user subject: %w", err)
	}

	apk, err := m.dbw.GetAccountSubject(ctx, accountName, m.operator.Name)
	if err != nil {
		return fmt.Errorf("could not fetch user public key: %w", err)
	}

	akp, err := m.secrets.ReadAccountSeed(apk)
	if err != nil {
		return fmt.Errorf("could not fetch account key pair: %w", err)
	}
	defer akp.Wipe()

	tok, err := m.secrets.ReadUserJWT(pk)
	if err != nil {
		return err
	}
	for _, edit := range editors {
		edit(tok)
	}

	if err = m.secrets.WriteUserJWT(tok, apk); err != nil {
		return fmt.Errorf("could not write user JWT: %w", err)
	}

	if err = m.SaveAppCreds(ctx, name, accountName); err != nil {
		return fmt.Errorf("could not save app credentials: %w", err)
	}

	return nil
}

func (m *Manager) SaveAppCreds(ctx context.Context, userName, accountName string) error {
	if app, ok := GetAppByUser(userName, accountName); ok {
		return m.saveAppCreds(ctx, app, userName, accountName)
	}
	return nil
}

func wipeBytes(bs []byte) {
	for i := range bs {
		bs[i] = 'X'
	}
}

func (m *Manager) saveAppCreds(ctx context.Context, appName, userName, accountName string) error {
	creds, err := m.GenerateUserCredentials(ctx, userName, accountName)
	if err != nil {
		return err
	}
	defer wipeBytes(creds)

	return m.secrets.WriteAppCreds(appName, creds)
}

func (m *Manager) operatorExists(ctx context.Context, name string) (bool, error) {
	if _, err := m.dbw.GetOperatorSubject(ctx, name); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (m *Manager) accountExists(ctx context.Context, name, operatorName string) (bool, error) {
	if _, err := m.dbw.GetAccountSubject(ctx, name, operatorName); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (m *Manager) userExists(ctx context.Context, name, accountName, operatorName string) (bool, error) {
	if _, err := m.dbw.GetUserSubject(ctx, name, accountName, operatorName); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// The Manager is initialized if the operator, the operator account, the operator user,
// the system account and system account all exist (and are valid). When none of these exist,
// the Manager is free to initialize. Otherwise, the user should decide on what to do
// (it would likely not be a very good idea to clear the database for the user).
func (m *Manager) requiresKeysInit(ctx context.Context) (bool, error) {
	operatorExists, err := m.operatorExists(ctx, m.operator.Name)
	if err != nil {
		return false, err
	}

	operatorAccountExists, err := m.accountExists(ctx, m.operator.AccountName, m.operator.Name)
	if err != nil {
		return false, err
	}

	operatorUserExists, err := m.userExists(ctx, m.operator.UserName, m.operator.AccountName, m.operator.Name)
	if err != nil {
		return false, err
	}

	sysAccountExists, err := m.accountExists(ctx, "SYS", m.operator.Name)
	if err != nil {
		return false, err
	}

	sysUserExists, err := m.userExists(ctx, "sys", "SYS", m.operator.Name)
	if err != nil {
		return false, err
	}

	if !boolsEqual(operatorExists, operatorAccountExists, operatorUserExists, sysAccountExists, sysUserExists) {
		return false, fmt.Errorf("partially initialized")
	}

	return !operatorExists, nil
}

func boolsEqual(b0 bool, bs ...bool) bool {
	for _, b := range bs {
		if b != b0 {
			return false
		}
		b0 = b
	}
	return true
}

// InitKeys initializes the system account and the operator.
func (m *Manager) InitKeys(ctx context.Context) error {
	required, err := m.requiresKeysInit(ctx)
	if err != nil {
		return err
	}
	if !required {
		return nil
	}

	m.log.Info("initializing keys")

	// Create the keys for the system account
	sakp, err := m.newAccountKeys()
	if err != nil {
		return fmt.Errorf("could not create system account keys: %w", err)
	}
	defer sakp.Wipe()

	// Create the public key for the system account
	sapk, err := sakp.PublicKey()
	if err != nil {
		return fmt.Errorf("could not create system account public key: %w", err)
	}

	// Create a new operator and configure the system account in its claims
	oppk, err := m.NewOperator(ctx, sapk)
	if err != nil {
		return fmt.Errorf("could not create operator: %w", err)
	}

	// Create an account for the operator
	oapk, err := m.newAccount(ctx, m.operator.AccountName, oppk)
	if err != nil {
		return fmt.Errorf("could not create operator account: %w", err)
	}

	// Create a user for the operator
	if _, err = m.newUser(ctx, m.operator.UserName, oapk); err != nil {
		return fmt.Errorf("could not create operator user: %w", err)
	}

	// Create a new system account, issued by the new operator
	if err = m.newSystemAccount(ctx, "SYS", sapk, oppk); err != nil {
		return fmt.Errorf("could not create system account: %w", err)
	}

	// Create a new user for the system account
	if _, err = m.newUser(ctx, "sys", sapk); err != nil {
		return fmt.Errorf("could not create system user: %w", err)
	}

	m.log.Info("keys initialized")

	return nil
}

// getOperatorPubKey retrieves the public key for the default operator.
func (m *Manager) getOperatorPubKey(ctx context.Context) (string, error) {
	pk, err := m.dbw.GetOperatorSubject(ctx, m.operator.Name)
	if err != nil {
		return "", fmt.Errorf("could not get operator subject: %w", err)
	}
	return pk, nil
}

// GetSystemAccountPubKey retrieves the public key for the system account.
func (m *Manager) GetSystemAccountPubKey(ctx context.Context) (string, error) {
	pk, err := m.dbw.GetAccountSubject(ctx, "SYS", m.operator.Name)
	if err != nil {
		return "", fmt.Errorf("could not get system account subject: %w", err)
	}
	return pk, nil
}

// ProvisionNewDevice provision a new device with an account and user.
// The ID for the device must be provided.
func (m *Manager) ProvisionNewDevice(ctx context.Context, id uuid.UUID) error {
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

// GenerateDeviceCredentials generates new NATS credentials for a device with a known ID.
func (m *Manager) GenerateDeviceCredentials(ctx context.Context, id uuid.UUID) ([]byte, error) {
	return m.GenerateUserCredentials(ctx, deviceUserName(id), deviceAccountName(id))
}

// GenerateUserCredentials generates new NATS credentials for a user with a known name and issuer (account).
func (m *Manager) GenerateUserCredentials(ctx context.Context, userName, accountName string) ([]byte, error) {
	// Get the subject for the user
	pk, err := m.dbw.GetUserSubject(ctx, userName, accountName, m.operator.Name)
	if err != nil {
		return nil, err
	}

	// Read the key pair for the user, the seed is part of the credentials file
	kp, err := m.secrets.ReadUserSeed(pk)
	if err != nil {
		return nil, err
	}
	defer kp.Wipe()

	// Extract the seed from the key pair
	seed, err := kp.Seed()
	if err != nil {
		return nil, err
	}

	// Read the user's JWT
	token, err := m.secrets.ReadJWTLiteral(pk)
	if err != nil {
		return nil, err
	}

	// Create the config
	cfg, err := jwt.FormatUserConfig(token, seed)
	if err != nil {
		return nil, fmt.Errorf("could not format user config: %w", err)
	}
	return cfg, nil
}

// AccountExists checks if an account with a known name exists. It returns false if the account
// doesn't exist and true if it does.
func (m *Manager) AccountExists(ctx context.Context, name string) (bool, error) {
	subj, err := m.dbw.GetAccountSubject(ctx, name, m.operator.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	if _, err = m.secrets.ReadAccountJWT(subj); err != nil {
		return false, err
	}
	return true, nil
}

// UserExists checks if a user with a known name and issuer (account) exists. It returns false if the user
// doesn't exist and true if it does. If the name of the user is known but the account name is not,
// it will still return false.
func (m *Manager) UserExists(ctx context.Context, name, accountName string) (bool, error) {
	subj, err := m.dbw.GetUserSubject(ctx, name, accountName, m.operator.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	if _, err = m.secrets.ReadUserJWT(subj); err != nil {
		return false, err
	}
	return true, nil
}

// GetOperatorJWT retrieves the operator JWT.
func (m *Manager) GetOperatorJWT(ctx context.Context) (string, error) {
	oppk, err := m.dbw.GetOperatorSubject(ctx, m.operator.Name)
	if err != nil {
		return "", err
	}

	claims, err := m.secrets.ReadJWTLiteral(oppk)
	if err != nil {
		return "", err
	}

	if _, err = jwt.DecodeOperatorClaims(claims); err != nil {
		return "", err
	}

	return claims, nil
}

func (m *Manager) GetUserKeyPair(ctx context.Context, name, accountName string) (nkeys.KeyPair, error) {
	subj, err := m.dbw.GetUserSubject(ctx, name, accountName, m.operator.Name)
	if err != nil {
		return nil, err
	}
	return m.secrets.ReadUserSeed(subj)
}

func (m *Manager) GetUserJWT(ctx context.Context, name, accountName string) (string, error) {
	subj, err := m.dbw.GetUserSubject(ctx, name, accountName, m.operator.Name)
	if err != nil {
		return "", err
	}
	return m.secrets.ReadJWTLiteral(subj)
}

// deviceAccountName generates a new name for an (embedded) device account.
func deviceAccountName(id uuid.UUID) string {
	return fmt.Sprintf("EMDEV-%s", id)
}

// deviceUserName generates a new name for an (embedded) device user.
func deviceUserName(id uuid.UUID) string {
	return fmt.Sprintf("emdev-%s", id)
}
