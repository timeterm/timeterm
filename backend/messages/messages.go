package messages

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"gitlab.com/timeterm/timeterm/backend/database"
	"gitlab.com/timeterm/timeterm/backend/secrets"
)

type Wrapper struct {
	log  logr.Logger
	dbw  *database.Wrapper
	secr *secrets.Wrapper
}

func NewWrapper(log logr.Logger, dbw *database.Wrapper, secr *secrets.Wrapper) *Wrapper {
	return &Wrapper{
		log:  log.WithName("MessagesWrapper"),
		dbw:  dbw,
		secr: secr,
	}
}

type Severity int

const (
	SeverityError Severity = iota
	SeverityInfo
)

func convertSeverity(s Severity) database.AdminMessageSeverity {
	switch s {
	case SeverityError:
		return database.AdminMessageSeverityError
	default:
		fallthrough
	case SeverityInfo:
		return database.AdminMessageSeverityInfo
	}
}

func (w *Wrapper) Start(organizationID uuid.UUID) Entry {
	return Entry{
		w:              w,
		organizationID: organizationID,
		severity:       SeverityInfo,
		verbosity:      0,
	}
}

type Fields map[string]interface{}

type Entry struct {
	w              *Wrapper
	organizationID uuid.UUID
	severity       Severity
	verbosity      int
	data           encryptedData
}

func (e Entry) V(v int) Entry {
	e.verbosity += v
	return e
}

func (e Entry) Error() Entry {
	e.severity = SeverityError
	return e
}

func (e Entry) Info() Entry {
	e.severity = SeverityInfo
	return e
}

func (e Entry) Summaryf(f string, a ...interface{}) Entry {
	e.data.Summary = fmt.Sprintf(f, a...)
	return e
}

func (e Entry) Messagef(f string, a ...interface{}) Entry {
	e.data.Message = fmt.Sprintf(f, a...)
	return e
}

func (e Entry) WithField(key string, value interface{}) Entry {
	return e.WithFields(Fields{key: value})
}

func (e Entry) WithFields(f Fields) Entry {
	newFields := make(Fields, len(e.data.Fields)+len(f))
	for k, v := range e.data.Fields {
		newFields[k] = v
	}
	for k, v := range f {
		newFields[k] = v
	}
	e.data.Fields = newFields
	return e
}

func (e Entry) Log() {
	if err := e.log(); err != nil {
		e.w.log.Error(err, "could not log entry")
	}
}

func (e Entry) log() error {
	logsKey, err := e.w.secr.GetOrganizationLogsKeySecret(e.organizationID)
	if err != nil {
		return fmt.Errorf("could not get organization logs secret: %w", err)
	}

	nonce, encrypted, err := e.data.encrypt(logsKey)
	if err != nil {
		return fmt.Errorf("could not encrypt message data: %w", err)
	}

	const createTimeout = time.Second * 5
	ctx, cancel := context.WithTimeout(context.Background(), createTimeout)
	defer cancel()

	err = e.w.dbw.CreateAdminMessage(ctx, database.AdminMessage{
		OrganizationID: e.organizationID,
		LoggedAt:       time.Now(),
		Severity:       convertSeverity(e.severity),
		Verbosity:      e.verbosity,
		Nonce:          nonce,
		Data:           encrypted,
	})
	if err != nil {
		return fmt.Errorf("could not create admin message (in database): %w", err)
	}

	return nil
}

type encryptedData struct {
	Summary string                 `json:"summary"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields"`
}

func (d encryptedData) encrypt(key []byte) ([]byte, []byte, error) {
	bytes, err := json.Marshal(&d)
	if err != nil {
		return nil, nil, fmt.Errorf("could not marshal message data as JSON: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create AES cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create AES GCM cipher: %w", err)
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, fmt.Errorf("could not generate nonce: %w", err)
	}

	return nonce, aesgcm.Seal(nil, nonce, bytes, nil), nil
}
