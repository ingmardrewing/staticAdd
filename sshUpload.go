package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func NewSshUploader(user, pass, host, port string) *SshUploader {
	su := new(SshUploader)
	su.username = user
	su.password = pass
	su.server = host
	su.port = port
	return su
}

type SshUploader struct {
	username string
	password string
	server   string
	port     string
}

func (su *SshUploader) upload(localpath, remotepath string) {

	config := &ssh.ClientConfig{
		User: su.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(su.password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	t := net.JoinHostPort(su.server, su.port)

	sshConn, err := ssh.Dial("tcp", t, config)
	if err != nil {
		fmt.Printf("Failed to connect to %v\n", t)
		fmt.Println(err)
		os.Exit(2)
	}

	session, err := sshConn.NewSession()
	if err != nil {
		fmt.Printf("Cannot create SSH session to %v\n", t)
		fmt.Println(err)
		os.Exit(2)
	}
	defer session.Close()

	go func() {

		w, err := session.StdinPipe()

		if err != nil {
			return
		}
		defer w.Close()

		fmt.Println("upload image")

		src, err := os.Open(localpath)
		panicOnError("os.Open", err)
		defer src.Close()

		srcStat, err := os.Stat(localpath)
		if err != nil {
			panic(err)
		}

		targetFile := filepath.Base(remotepath)
		_, err = fmt.Fprintln(w, "C0644", srcStat.Size(), targetFile)
		panicOnError("C0644", err)

		if srcStat.Size() > 0 {
			n, err := io.Copy(w, src)
			panicOnError("Copy", err)
			fmt.Println(n)

			_, err = fmt.Fprint(w, "\x00")
			panicOnError("\x00", err)

		} else {
			_, err = fmt.Fprint(w, "\x00")
			panicOnError("\x00", err)

		}
	}()

	fmt.Println("start upload")
	err = session.Run(fmt.Sprintf("scp -tr %s", remotepath))
	panicOnError("Run scp", err)
	session.Close()

	session22, err := sshConn.NewSession()
	if err != nil {
		fmt.Printf("Cannot create SSH session to %v\n", t)
		fmt.Println(err)
		os.Exit(2)
	}
	// Close the session when main returns
	defer session22.Close()

}

func panicOnError(msg string, err error) {
	if err != nil {
		panic(fmt.Errorf("%s: %v", msg, err))
	}
}
