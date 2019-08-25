package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/n0stack/n0stack/n0core/pkg/deploy"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
)

func DeployAgent(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return fmt.Errorf("Argument usage: %s", ctx.Command.ArgsUsage)
	}
	h := strings.Split(ctx.Args()[0], "@")
	if len(h) != 2 {
		return fmt.Errorf("Argument usage: %s", ctx.Command.ArgsUsage)
	}

	user := h[0]
	host := h[1]
	keyPath := ctx.String("identity-file")
	args := strings.Join(ctx.Args()[1:], " ")
	target := ctx.String("base-directory")
	sendTo := ctx.String("send-to")

	buf, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return errors.Wrap(err, "Failed to read key file")
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return errors.Wrap(err, "Failed to parse key")
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.SetDefaults()
	conn, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return errors.Wrap(err, "Failed to dial ssh")
	}
	defer conn.Close()

	d, err := deploy.NewRemoteDeployer(conn, sendTo)
	if err != nil {
		return err
	}

	binLocation := "n0core." + version
	fmt.Printf("---> [DEPLOY] Sending self to %s...\n", binLocation)
	self, err := d.ReadSelf()
	if err != nil {
		return err
	}
	if err := d.SendFile(self, binLocation, 0755); err != nil {
		return err
	}

	cmd := fmt.Sprintf("%s install agent -base-directory %s -arguments '%s'", filepath.Join(sendTo, binLocation), target, args)
	fmt.Printf("---> [DEPLOY] Running install '%s'...\n", cmd)
	if err := d.Command(cmd, os.Stdout, os.Stderr); err != nil {
		return err
	}

	return nil
}
