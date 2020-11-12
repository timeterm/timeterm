package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"

	"github.com/google/uuid"
)

type nscConfig struct {
	dataDir string
}

func (c *nscConfig) storeDir() string {
	return path.Join(c.dataDir, "store")
}

func (c *nscConfig) nkeysPath() string {
	return path.Join(c.dataDir, "nkeys")
}

func (c *nscConfig) nscHome() string {
	return path.Join(c.dataDir, "nsc")
}

func (c *nscConfig) env() []string {
	return []string{
		"NKEYS_PATH=" + c.nkeysPath(),
		"NSC_HOME=" + c.nscHome(),
	}
}

type aclEntry struct {
	op    topicOp
	topic string
}

func (e aclEntry) asAddUserArgs() []string {
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

func createNewDevUser(id uuid.UUID, c *nscConfig) (string, error) {
	accountName := fmt.Sprintf("fedev-%s", id)
	userName := accountName

	if err := runNscCmd(nscAddAccountCmd(accountName), c); err != nil {
		return "", fmt.Errorf("could not create account %q: %w", accountName, err)
	}

	userConfig := createDevUserConfig(id)
	if err := runNscCmd(nscAddUserCmd(accountName, userName, userConfig), c); err != nil {
		return "", fmt.Errorf("could not create user %q for account %q: %w", userName, accountName, err)
	}

	return nscGenerateUserCreds(accountName, userName)
}

func nscGenerateUserCreds(account, user string) (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("nsc", "generate", "creds", "--account", account, "--name", user)
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func createDevUserConfig(id uuid.UUID) userConfig {
	return userConfig{
		streams: []streamConsumer{
			{
				stream:   "EMDEV-DISOWN-TOKEN",
				consumer: fmt.Sprintf("EMDEV-%s-EMDEV-DISOWN-TOKEN", id),
			},
		},
		other: []aclEntry{
			{
				topic: fmt.Sprintf("EMDEV.%s.REBOOT", id),
				op:    topicOpSub,
			},
		},
	}
}

type streamConsumer struct {
	stream, consumer string
}

func nscAddAccountCmd(name string) *exec.Cmd {
	args := []string{"add", "account", "--name", name}

	return exec.Command("nsc", args...)
}

type topicOp int

const (
	topicOpNone topicOp = 0
	topicOpPub  topicOp = 1 << (iota - 1)
	topicOpSub
)

type userConfig struct {
	streams []streamConsumer
	other   []aclEntry
}

func nscAddUserCmd(account, name string, cfg userConfig) *exec.Cmd {
	args := []string{"add", "user", "--name", name, "--account", account}

	acls := cfg.other
	for _, sc := range cfg.streams {
		acls = append(acls, streamConsumerACLs(sc)...)
	}

	for _, aclEntry := range acls {
		args = append(args, aclEntry.asAddUserArgs()...)
	}

	return exec.Command("nsc", args...)
}

func nscInitCmd(storeDir string) *exec.Cmd {
	return exec.Command("nsc", "init", "--name", "timeterm", "--dir", storeDir)
}

func runNscCmd(cmd *exec.Cmd, c *nscConfig) error {
	envCmd := exec.Command("nsc", "env", "--store", c.storeDir(), "--operator", "timeterm")
	envCmd.Env = c.env()
	if err := envCmd.Run(); err != nil {
		return err
	}

	cmd.Env = c.env()
	return cmd.Run()
}
