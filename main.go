//Anchor CM: 1.0.0
//DSL Spec: 1.0.0

//License: MIT
//Author: Gani Mendoza (itjumpstart.wordpress.com)

/*

Anchor Configuration Management Workflow

Note:

Anchor CM is designed for configuring one machine only
Configuring multiple machines is up to you (your workflow) OR
you can delegate it to a separate orchestration tool

Anchor CM is Data Plane

- Anchor CM is data (it is just a sequence of instructions)
- No conditionals
- No loop

Orchestration is Control Plane

- Control plane is sequence
- Control plane has conditional
- Control plane has loop

Assumption:

Your SSH key has been uploaded and configured accordingly at remote machine

Definitions:

control machine - SSH client
remote machine - machine to be configured (must be SSH server)
cmdfile - your virtual script to be executed by Anchor CM


1. Prepare cmdfile and other optional files on your control machine

2. SSH to remote machine (using SSH agent recommended)

3. SCP cmdfile and other optional files to remote machine

4. From your local control machine, execute cmdfile at remote machine via SSH

anchor /path/to/cmdfile

Note:

There is only one cmdfile for execution but you can include multiple cmdfiles

5. Debug error displayed on your control machine (if any)

6. Done!

*/

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/hoisie/mustache"
	"github.com/ibmendoza/go-ini"
	anko_core "github.com/mattn/anko/builtins"
	anko_encoding "github.com/mattn/anko/builtins/encoding"
	anko_flag "github.com/mattn/anko/builtins/flag"
	anko_io "github.com/mattn/anko/builtins/io"
	anko_math "github.com/mattn/anko/builtins/math"
	anko_net "github.com/mattn/anko/builtins/net"
	anko_os "github.com/mattn/anko/builtins/os"
	anko_path "github.com/mattn/anko/builtins/path"
	anko_regexp "github.com/mattn/anko/builtins/regexp"
	anko_sort "github.com/mattn/anko/builtins/sort"
	anko_strings "github.com/mattn/anko/builtins/strings"
	anko_term "github.com/mattn/anko/builtins/term"
	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

//global variables
var shell, flag string

func printError(err error) {

	if err != nil {
		ct.ChangeColor(ct.Red, false, ct.None, false)
		os.Stderr.WriteString(fmt.Sprintf("==> ERROR: %s\n", err.Error()))
		ct.ResetColor()
	}

}

func printOutput(outs []byte) {

	if len(outs) > 0 {
		ct.ChangeColor(ct.Green, false, ct.None, false)
		fmt.Printf("==> OUTPUT: %s\n", string(outs))
		ct.ResetColor()
	}
}

func getShellAndFlag() (string, string) {
	if runtime.GOOS == "windows" {
		return "cmd", "/C"
	} else {
		return "/bin/sh", "-c"
	}
}

func runCmd(args string) error {
	fmt.Println(args)

	splitSpace := strings.Split(args, " ")

	var err error
	switch splitSpace[0] {
	case "mkdir":
		dir := splitSpace[1]

		//make dir only if not existing
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			cmd := exec.Command(shell, flag, args)

			output, err := cmd.CombinedOutput()
			printError(err)
			printOutput(output)
		}
	default:
		cmd := exec.Command(shell, flag, args)

		output, err := cmd.CombinedOutput()
		printError(err)
		printOutput(output)
	}
	return err
}

func runflagCmd(args string, cmdfile ini.File) error {

	/*
		- Flag sections must only include key/value entries
		- To simplify parsing, flag must be either single dash or double dashes
		- As a consequence of above, do not mix % and @ in one line
		- E.g. docker -d %section1 @section2 (wrong)

		https://docs.docker.com/reference/commandline/cli/

		Case 1: single dash flags

		docker run -d -m 100m -e DEVELOPMENT=1 \
		-e BRANCH=example-code \
		-v $(pwd):/app/bin:ro \
		--name app appserver

		can be rewritten as

		[docker]

		m = 100m
		e = DEVELOPMENT=1
		e = BRANCH=example-code
		v = $(pwd):/app/bin:ro

		[code]

		RUNFLAG sudo docker run -d --name app appserver %docker


		Case 2: double dash flags

		VBoxManage modifyvm tklinux --nic1 bridged --nic2 hostonly

		can be rewritten as

		[network]
		nic1 = bridged
		nic2 = hostonly

		[code]

		RUNFLAG VBoxManage modifyvm tklinux @network

	*/

	var err error

	fmt.Println(args)

	//for docker cli, it's better to use --flags for readability

	if strings.Contains(args, "%") && strings.Contains(args, "@") {
		err = errors.New("RUNFLAG cannot contain % and @ in the same line")
		printError(err)
		return err
	}

	var strSymbol, strDash string
	//address case 1
	if strings.Contains(args, "%") {
		strSymbol = "%"
		strDash = " -"
	}

	//address case 2
	if strings.Contains(args, "@") {
		strSymbol = "@"
		strDash = " --"
	}

	if strings.Contains(args, strSymbol) {
		slcArgs := strings.Split(args, strSymbol)

		//RUNFLAG VBoxManage modifyvm tklinux @network

		// RUNFLAG VBoxManage modifyvm tklinux
		args1 := slcArgs[0]

		// ::network
		section := slcArgs[1]

		var flags string

		for key, value := range cmdfile[section] {

			//[network]
			//nic1 = bridged
			//nic2 = hostonly

			val, err := eval(value)

			if err != nil {
				printError(err)
				break
				return err

			} else {

				flags = flags + strDash + key + " " + val
			}
		}

		fmt.Println(args1 + flags)

		args = args1 + flags
	}

	cmd := exec.Command(shell, flag, args)

	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)

	return err
}

//returns itself if plain string otherwise executes the command
func eval(v string) (string, error) {
	/*
		[network]
		nic1 = bridged
		nic2 = hostonly

		[consul]

		keygen = RUN consul keygen

		[security]

		# no need to enclose in single or double quote

		pem = READFILE /etc/ssl/certs/cert.pem

	*/

	var err error

	//READFILE and RUN are reserved words in cmdfile
	isAt := strings.Contains(v, "READFILE") || strings.Contains(v, "RUN")

	//nic1 = bridged
	if !isAt {
		return v, nil
	}

	if strings.Contains(v, "READFILE") {
		slcStr := strings.Split(v, "READFILE")
		v = strings.TrimSpace(slcStr[1])

		_, err = os.Stat(v)
		if err == nil {
			//return content of file
			var fileContent []byte

			fileContent, err = ioutil.ReadFile(v)

			if err != nil {
				return "", err
			} else {
				return string(fileContent), nil
			}
		}
	}

	if strings.Contains(v, "RUN") {

		v := strings.Replace(v, "RUN", "", -1)
		v = strings.TrimSpace(v)

		cmd := exec.Command(shell, flag, v)

		var output []byte

		output, err = cmd.CombinedOutput()

		if err != nil {
			return "", err
		} else {
			return string(output), nil
		}
	}

	//compiler insists a return
	return "", err
}

func chdirCmd(dir string) error {
	fmt.Println("chdir " + dir)

	err := os.Chdir(dir)

	if err != nil {
		printError(err)
	} else {
		printOutput([]byte("chdir to " + dir))
	}
	return err
}

func getenvCmd(key string) error {
	fmt.Println("getenv " + key)

	result := os.Getenv(key)

	if len(result) == 0 {
		err := errors.New("No environment variable named " + key)
		printError(err)
		return err
	} else {
		printOutput([]byte("getenv " + key + "=" + result))
		return nil
	}
}

func setenvCmd(key, value string) error {
	if key == "" || value == "" {
		return errors.New("Error in ENV. Key or value is blank")
	}

	fmt.Println("ENV " + key + " " + value)

	err := os.Setenv(key, value)

	if err != nil {
		printError(err)
		return err
	} else {
		printOutput([]byte("ENV " + key + "=" + value))
		return nil
	}
}

func hostenvCmd(key string) error {
	fmt.Println("hostenv " + key)

	slc := os.Environ()

	found := false

	for _, v := range slc {
		//fmt.Println(slc[k])
		pair := strings.Split(v, "=")

		if pair[0] == key {

			printOutput([]byte("hostenv: " + key + "=" + v))

			found = true
			break
		}
	}

	if !found {
		err := errors.New("No host environment variable named " + key)
		printError(err)
		return err
	} else {
		return nil
	}
}

func goCmd(argCommand string, args []string) error {
	var err error
	switch argCommand {

	case "chdir":
		if len(args) != 1 {
			err = errors.New("GO chdir. Directory not specified")
			printError(err)
			return err
		} else {
			err = chdirCmd(args[0])
		}

	case "getenv":
		if len(args) != 1 {
			err = errors.New("GO getenv. Key is blank")
			printError(err)
			return err
		} else {
			err = getenvCmd(args[0])
		}

	case "hostenv":
		if len(args) != 1 {
			err = errors.New("GO setenv. Key is blank")
			printError(err)
			return err
		} else {
			err = hostenvCmd(args[0])
		}

	case "hostname":
		var str string
		str, err = os.Hostname()
		if err != nil {
			printError(err)
			return err
		} else {
			printOutput([]byte("hostname: " + str))
		}

	}

	return err
}

func ankoCmd(filename string) error {
	fmt.Println("ANKO " + filename)

	if len(filename) == 0 {
		err := errors.New("Please specify an Anko script file")
		printError(err)
		return err
	}

	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		printError(err)
		return err
	}

	env := vm.NewEnv()

	anko_core.Import(env)
	anko_flag.Import(env)
	anko_net.Import(env)
	anko_encoding.Import(env)
	anko_os.Import(env)
	anko_io.Import(env)
	anko_math.Import(env)
	anko_path.Import(env)
	anko_regexp.Import(env)
	anko_sort.Import(env)
	anko_strings.Import(env)
	anko_term.Import(env)

	var ln, code string

	lnScanner := bufio.NewScanner(file)
	for lnScanner.Scan() {

		ln = lnScanner.Text()

		code = code + ln + "\n"

		if err != nil {
			break
			printError(err)
			return err
		}
	}

	scanner := new(parser.Scanner)

	scanner.Init(code)

	stmts, err := parser.Parse(scanner)

	if err != nil {
		printError(err)
		return err
	}
	_, err = vm.Run(stmts, env)
	if err != nil {
		printError(err)
		return err
	}

	return err
}

func includeCmd(filename string) error {

	//can include nested cmdfile
	fmt.Println("INCLUDE " + filename)

	err := parseCmdfile(filename)

	if err != nil {
		printError(err)
		return err
	}

	return err
}

//USAGE: TEMPLATE template file, json file, destination file, shell script file
//Uses Mustache to render template + json into destination file
//Optional shell script file to run after completion
//WARNING: Destination file will be overwritten
//TEST: TEMPLATE test.tmpl test.json testconfig.txt testecho.bat
func templateCmd(args string) (err error) {

	slcStr := strings.Split(args, " ")

	if len(slcStr) < 3 || len(slcStr) > 4 {

		err = errors.New("TEMPLATE expects in order: template, json, " +
			"destination file and optional shell script")

		printError(err)

		return err
	}

	strTemplate := slcStr[0]
	strJson := slcStr[1]
	strDestination := slcStr[2]

	var strShell string
	if len(slcStr) == 4 {
		strShell = slcStr[3]
	}

	var fileContent []byte
	var jsonData map[string]interface{}
	var template string

	//read template file
	fileContent, err = ioutil.ReadFile(strTemplate)

	if err != nil {
		err = errors.New("Error reading template file named " + strTemplate)
		printError(err)
		return err
	} else {
		template = string(fileContent)
	}

	//read json file
	fileContent, err = ioutil.ReadFile(strJson)
	if err != nil {
		err = errors.New("Error reading json file named " + strJson)
		printError(err)
		return err
	} else {

		err := json.Unmarshal(fileContent, &jsonData)

		if err != nil {
			err = errors.New("Error parsing JSON data file named " + strJson)
			printError(err)
			return err
		}
	}

	//render mustache template + JSON data
	rendered := mustache.Render(template, jsonData)

	//write rendered template to file
	var file *os.File
	file, err = os.Create(strDestination)
	defer file.Close()
	if err != nil {
		err = errors.New("Error writing to destination file named " + strDestination)
		printError(err)
		return err
	} else {
		w := bufio.NewWriter(file)
		_, err := w.WriteString(rendered)
		if err != nil {
			err = errors.New("Error writing rendered template to destination file named " + strDestination)
			printError(err)
			return err
		} else {
			w.Flush()
		}
	}

	//execute shell file if any
	if strShell > "" {
		_, err = os.Stat(strShell)

		if err != nil {
			err = errors.New("Shell script file named " + strShell + " does not exist")
			return err
		} else {
			cmd := exec.Command(shell, flag, strShell)

			output, err := cmd.CombinedOutput()
			printError(err)
			printOutput(output)
		}
	}

	return err
}

func processCmd(command string, cmdfile ini.File) error {
	var err error

	s := strings.TrimSpace(command)

	slcStr := strings.Split(s, " ")

	args := []string{}

	var argCommand string

	cmd := strings.ToUpper(slcStr[0])

	if !strings.Contains(cmd, "FROM") || !strings.Contains(cmd, "MAINTAINER") ||
		!strings.Contains(cmd, "LICENSE") {
		fmt.Println(cmd)
	}

	for i, _ := range slcStr {

		if i == 1 {
			argCommand = slcStr[i]
		} else if i > 1 {
			args = append(args, slcStr[i])
		}
	}

	switch cmd {

	case "RUN":
		err = runCmd(strings.Join(slcStr[1:], " "))

	case "GO":
		err = goCmd(argCommand, args)

	case "ENV":
		err = setenvCmd(argCommand, args[0])

	case "ANKO":
		ankoFilename := argCommand
		err = ankoCmd(ankoFilename)

	case "RUNFLAG":
		err = runflagCmd(strings.Join(slcStr[1:], " "), cmdfile)

	case "INCLUDE":
		filename := argCommand
		err = includeCmd(filename)

	case "TEMPLATE":
		err = templateCmd(strings.Join(slcStr[1:], " "))

	}

	return err
}

func parseIniSections(filename string) string {
	file, err := os.Open(filename)

	if err != nil {
		printError(err)
		log.Fatal(err)
	}

	var line string
	scanner := bufio.NewScanner(file)

	hasReachedCodeSection := false
	strForIniParsing := ""

	for scanner.Scan() {

		line = scanner.Text()

		line = strings.TrimSpace(line)

		//skip blank line
		if len(line) != 0 {

			if strings.Contains(line, "[code]") {
				hasReachedCodeSection = true
				continue
			}

			if !hasReachedCodeSection {

				//build string for ini parsing
				if strings.Contains(line, "[") ||
					strings.Contains(line, "]") ||
					strings.Contains(line, "=") {

					strForIniParsing = strForIniParsing + line + "\n"
				}
			}
		}

		if err != nil {
			break
			log.Fatal(err)
		}
	}

	file.Close()

	return strForIniParsing
}

func parseCode(filename string, cmdfile ini.File) error {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		printError(err)
		log.Fatal(err)
	}

	var line string
	scanner := bufio.NewScanner(file)

	hasReachedCodeSection := false

	for scanner.Scan() {

		line = scanner.Text()

		line = strings.TrimSpace(line)

		//skip blank line
		if len(line) != 0 {

			if strings.Contains(line, "[code]") {
				hasReachedCodeSection = true
				continue
			}

			if hasReachedCodeSection {
				if !strings.Contains(line, "#") &&
					!strings.Contains(line, "[") &&
					!strings.Contains(line, "]") &&
					!strings.Contains(line, "=") &&
					!strings.Contains(line, ";") {
					err = processCmd(line, cmdfile)
				}
			}
		}

		if err != nil {
			break

			return err
		}
	}

	return err
}

func parseCmdfile(filename string) error {

	str := parseIniSections(filename)

	//convert string to io.Reader
	input := bytes.NewBufferString(str)

	cmdfile, err := ini.Load(input)

	if err != nil {
		return errors.New("Stop parsing ini section(s) of " + filename)
	}

	err = parseCode(filename, cmdfile)
	if err != nil {
		return errors.New("Stop parsing code section of " + filename)
	}

	return err
}

func main() {

	filename := os.Args[1]

	if len(os.Args) != 2 {
		printError(errors.New("Please specify a cmdfile"))
		os.Exit(1)
	}

	shell, flag = getShellAndFlag()

	err := parseCmdfile(filename)

	if err != nil {
		printError(err)
	}

	fmt.Println("")
	fmt.Println("If any error appears, cmdfile run is not completed")
}
