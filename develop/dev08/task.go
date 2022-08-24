package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func execCmd(args []string, in, out *os.File) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = out
	cmd.Stdin = in
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func execCommand(command []string) (string, error) {
	if len(command) == 0 {
		return "", nil
	}
	switch command[0] {
	case "cd":
		if len(command) < 2 {
			return "", errors.New("error length arguments")
		}
		err := os.Chdir(command[1])
		if err != nil {
			return "", err
		}
	case "pwd":
		curPath, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return fmt.Sprint("\nPath\n----\n\"", curPath, "\"\n\n"), nil
	case "echo":
		out := ""
		for i := 1; i < len(command); i++ {
			out += command[i] + "\n"
		}
		return out, nil
	case "ps":
		out, err := exec.Command("powershell", "ps").Output()
		if err != nil {
			return "", err
		}
		return string(out), nil
	case "kill":
		if len(command) < 2 {
			return "", errors.New("error length arguments")
		}
		out, err := exec.Command("powershell", "taskkill /f /IM "+command[1]).Output()
		if err != nil {
			return "", err
		}
		return string(out), nil
	case "exit":
		os.Exit(0)
	case "fork":
		if len(command) < 2 {
			return "", errors.New("error length arguments")
		}
		go execCmd(command[1:], os.Stdin, os.Stdout)
	case "exec":
		if len(command) < 2 {
			return "", errors.New("error length arguments")
		}
		err := execCmd(command[1:], os.Stdin, os.Stdout)
		if err != nil {
			return "", err
		} else {
			os.Exit(0)
		}
	default:
		err := execCmd(command, os.Stdin, os.Stdout)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func pipe(cmdString string) (string, error) {
	commands := strings.Split(cmdString, "|")
	var out io.ReadCloser
	var outStr string
	for i := 0; i < len(commands)-1; i++ {
		commands[i] = strings.TrimSpace(commands[i])
		cmdSplit := strings.Fields(commands[i])
		if len(cmdSplit) < 1 {
			return "", errors.New("error lenght arguments")
		}
		var args []string
		if outStr != "" {
			args = append(cmdSplit, outStr)
		} else {
			args = cmdSplit
		}
		cmd := exec.Command(args[0], args[1:]...)
		outStr = ""
		var err error
		out, err = cmd.StdoutPipe()
		if err != nil {
			return "", err
		}
		errStart := cmd.Start()
		if errStart != nil {
			res, errRes := execCommand(args)
			if errRes != nil {
				return "", errRes
			}
			outStr = res
			continue
		}
		data, errRead := ioutil.ReadAll(out)
		if errRead != nil {
			return "", errRead
		}
		outStr = string(data)
	}
	outStr = strings.Trim(outStr, "\n")
	cmdSplit := strings.Fields(commands[len(commands)-1])
	args := append(cmdSplit, outStr)
	cmd := exec.Command(args[0], args[1:]...)
	pipe, errPipe := cmd.StdoutPipe()
	if errPipe != nil {
		return "", errPipe
	}
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	errRun := cmd.Start()
	if errRun != nil {
		res, errRes := execCommand(args)
		if errRes != nil {
			return "", errRes
		}
		return string(res), nil
	}
	b, errReadPipe := ioutil.ReadAll(pipe)
	if errReadPipe != nil {
		return "", errReadPipe
	}
	return string(b), nil
}

func runShell() {
	reader := bufio.NewReader(os.Stdin)
	for loop := true; loop; {
		pwd, err := os.Getwd()
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}
		os.Stdout.WriteString(pwd + "> ")
		cmdString, errRead := reader.ReadString('\n')
		if errRead != nil {
			fmt.Fprintln(os.Stderr, errRead)
		}
		cmdString = strings.TrimSuffix(cmdString, "\n")
		cmdString = strings.TrimSuffix(cmdString, "\r")
		var out string
		if strings.Contains(cmdString, "|") {
			var errPipe error
			out, errPipe = pipe(cmdString)
			if errPipe != nil {
				os.Stderr.WriteString(errPipe.Error() + "\n")
				continue
			}
		} else {
			cmdSplit := strings.Fields(cmdString)
			var errExec error
			out, errExec = execCommand(cmdSplit)
			if errExec != nil {
				os.Stderr.WriteString(errExec.Error() + "\n")
				continue
			}
		}
		os.Stdout.WriteString(out)
	}
}

func main() {
	runShell()
}
