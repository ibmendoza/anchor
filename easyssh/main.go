package main

import (
	"flag"
	"fmt"
	"github.com/ibmendoza/easyssh"
	"log"
)

//For example, pass -user root -server 192.168.56.101 -keypath path/to/id_rsa -cmd pwd

var user = flag.String("user", "", "username")
var server = flag.String("server", "", "server name or IP address")
var pswd = flag.String("pswd", "", "password")
var keypath = flag.String("keypath", "", "keypath")
var scp = flag.String("scp", "", "optional file to upload")
var cmd = flag.String("cmd", "", "optional command to run")

func main() {

	flag.Parse()

	ssh := &easyssh.MakeConfig{
		User:     *user,
		Server:   *server,
		Key:      *keypath,
		Password: *pswd,
		Port:     "22",
	}

	// Call Run method with command you want to run on remote server.
	if *cmd != "" {
		response, err := ssh.Run(*cmd)

		if err != nil {
			log.Fatal("Can't run remote command: " + err.Error())
		} else {
			fmt.Println(response)
		}
	}

	if *scp != "" {
		// Call Scp method with file you want to upload to remote server.
		err := ssh.Scp(*scp)

		if err != nil {
			log.Fatal("Can't run remote command: " + err.Error())
		} else {
			fmt.Println("success")

			response, _ := ssh.Run("ls -al " + *scp)

			fmt.Println(response)
		}
	}
}
