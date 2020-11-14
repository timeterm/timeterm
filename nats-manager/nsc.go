package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

const nscOperatorAccount = "TIMETERM"

type topicOp int

const (
	topicOpNone topicOp = 0
	topicOpPub  topicOp = 1 << (iota - 1)
	topicOpSub
)

type streamConsumer struct {
	stream, consumer string
}

type userConfig struct {
	streams []streamConsumer
	other   []aclEntry
}

type aclEntry struct {
	op    topicOp
	topic string
}

func (e aclEntry) asNscAddUserArgs() []string {
	flag := "--allow-"
	if e.op == topicOpNone {
		return nil
	}
	if e.op&topicOpPub != 0 {
		flag += "pub"
	}
	if e.op&topicOpSub != 0 {
		flag += "sub"
	}

	return []string{flag, e.topic}
}

func streamConsumerACLs(c streamConsumer) []aclEntry {
	// See https://github.com/nats-io/jetstream#acls
	return []aclEntry{
		{topicOpPub, fmt.Sprintf("$JS.API.CONSUMER.MSG.NEXT.%s.%s", c.stream, c.consumer)},
		{topicOpPub, fmt.Sprintf("$JS.ACK.%s.%s.>", c.stream, c.consumer)},
	}
}

type nsc struct {
	dataDir string
	nscPath string
}

func newNsc(log logr.Logger, dataDir, nscPath string) (*nsc, error) {
	n := nsc{
		dataDir: dataDir,
		nscPath: nscPath,
	}

	needsInit, err := needsInit(n.dataDir)
	if err != nil {
		return nil, fmt.Errorf("could not check if already initialized: %w", err)
	}

	if needsInit {
		log.Info("nsc initialization required, initializing")

		err = n.runCmd(n.initCmd())
		if err != nil {
			return nil, fmt.Errorf("could not init nsc: %w", err)
		}

		log.Info("nsc initialized")
	}

	return &n, nil
}

func (n *nsc) storeDir() string {
	return path.Join(n.dataDir, "store")
}

func (n *nsc) nkeysPath() string {
	return path.Join(n.dataDir, "nkeys")
}

func (n *nsc) nscHome() string {
	return path.Join(n.dataDir, "nsc")
}

func (n *nsc) env() []string {
	return []string{
		"NKEYS_PATH=" + n.nkeysPath(),
		"NSC_HOME=" + n.nscHome(),
	}
}

func (n *nsc) createNewDevUser(id uuid.UUID) (string, error) {
	accountName := fmt.Sprintf("EMDEV-%s", id)
	userName := accountName

	if err := n.runCmd(n.addAccountCmd(accountName)); err != nil {
		return "", fmt.Errorf("could not create account %q: %w", accountName, err)
	}

	userConfig := createDevUserConfig(id)
	if err := n.runCmd(n.addUserCmd(accountName, userName, userConfig)); err != nil {
		return "", fmt.Errorf("could not create user %q for account %q: %w", userName, accountName, err)
	}

	return n.generateUserCreds(accountName, userName)
}

func (n *nsc) initCmd() *exec.Cmd {
	return exec.Command(n.nscPath, "init", "--name", nscOperatorAccount, "--dir", n.storeDir())
}

func (n *nsc) addAccountCmd(name string) *exec.Cmd {
	args := []string{"add", "account", "--name", name}

	return exec.Command(n.nscPath, args...)
}

func (n *nsc) addUserCmd(account, name string, cfg userConfig) *exec.Cmd {
	args := []string{"add", "user", "--name", name, "--account", account}

	acls := cfg.other
	for _, sc := range cfg.streams {
		acls = append(acls, streamConsumerACLs(sc)...)
	}

	for _, aclEntry := range acls {
		args = append(args, aclEntry.asNscAddUserArgs()...)
	}

	return exec.Command(n.nscPath, args...)
}

func (n *nsc) generateUserCreds(account, user string) (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command(n.nscPath, "generate", "creds", "--account", account, "--name", user)
	cmd.Stdout = &buf

	if err := n.runCmd(cmd); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (n *nsc) runCmd(cmd *exec.Cmd) error {
	envCmd := exec.Command(n.nscPath, "env", "--store", n.storeDir(), "--operator", nscOperatorAccount)
	envCmd.Env = append(os.Environ(), append(envCmd.Env, n.env()...)...)
	if err := runNscCmdWithLog(envCmd); err != nil {
		return err
	}

	cmd.Env = append(os.Environ(), append(cmd.Env, n.env()...)...)
	return runNscCmdWithLog(cmd)
}

type teeWriter struct {
	a, b io.Writer
}

func (t teeWriter) Write(p []byte) (n int, err error) {
	n, err = t.a.Write(p)
	if err != nil {
		n, err = t.b.Write(p)
	} else {
		_, _ = t.b.Write(p)
	}
	return
}

func makeTeeWriter(a, b io.Writer) io.Writer {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return teeWriter{a, b}
}

func runNscCmdWithLog(cmd *exec.Cmd) error {
	oldStderr, oldStdout := cmd.Stderr, cmd.Stdout

	var buf bytes.Buffer
	cmd.Stderr, cmd.Stdout = makeTeeWriter(&buf, oldStderr), makeTeeWriter(&buf, oldStdout)

	err := cmd.Run()
	if err != nil {
		return nscError{err: err, log: buf.String()}
	}
	return nil
}

type nscError struct {
	err error
	log string
}

func (e nscError) Error() string {
	return e.err.Error()
}

func (e nscError) Unwrap() error {
	return e.err
}

func (e nscError) Log() string {
	return e.log
}
