//License: MIT
//Author: Gani Mendoza (itjumpstart.wordpress.com)

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/daviddengcn/go-colortext"
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
	"github.com/vaughan0/go-ini"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

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

func runCmd(args string) error {
	fmt.Println(args)

	splitSpace := strings.Split(args, " ")

	//based on hashicorp serf
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}

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
		- As a consequence of above, do not mix single/double colon variables
		- E.g. docker -d :singledash ::doubledash (wrong)

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

		RUNFLAG sudo docker run -d --name app appserver :docker


		Case 2: double dash flags

		[network]
		nic1 = bridged
		nic2 = hostonly

		[code]

		RUNFLAG VBoxManage modifyvm tklinux ::network

	*/

	var err error

	fmt.Println(args)

	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}

	if strings.Contains(args, ":") && strings.Contains(args, "::") {
		err = errors.New("RUNFLAG cannot contain single and double colon variables")
		printError(err)
		return err
	}

	//address case 1
	if strings.Contains(args, ":") {
		slcArgs := strings.Split(args, ":")

		args1 := slcArgs[0]
		section := slcArgs[1]

		var flags string

		for key, value := range cmdfile[section] {
			flags = flags + " -" + key + " " + value
		}

		fmt.Println(args1 + flags)

		args = args1 + flags
	}

	//address case 2
	if strings.Contains(args, "::") {
		slcArgs := strings.Split(args, "::")

		args1 := slcArgs[0]
		section := slcArgs[1]

		var flags string

		for key, value := range cmdfile[section] {
			flags = flags + " --" + key + " " + value
		}

		fmt.Println(args1 + flags)

		// Output:
		// VBoxManage modifyvm tklinux --nic1 bridged --nic2 hostonly

		args = args1 + flags
	}

	cmd := exec.Command(shell, flag, args)

	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)

	return err
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

	file.Close()

	return err
}

func includeCmd(filename string) error {
	fmt.Println("INCLUDE " + filename)

	err := parseCmdfile(filename)

	if err != nil {
		printError(err)
		return err
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

	if !strings.Contains(cmd, "FROM") || !strings.Contains(cmd, "MAINTAINER") {
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

			if hasReachedCodeSection {
				if !strings.Contains(line, "#") &&
					!strings.Contains(line, "[") &&
					!strings.Contains(line, "]") &&
					!strings.Contains(line, "=") &&
					!strings.Contains(line, ";") {
				}
			} else {
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
	//the subfunctions are responsible to print error

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

	err := parseCmdfile(filename)

	if err != nil {
		printError(err)
	}

	fmt.Println("")
	fmt.Println("If any error appears, cmdfile is not completed")
}
