package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/pkg/errors"
)

func runCommand() error {
	if isInputFromPipe() {
		print("data is from pipe")
		return jsonifyText(os.Stdin, os.Stdout)
	} else {
		file, e := getFile()
		if e != nil {
			return e
		}
		defer file.Close()
		return jsonifyText(file, os.Stdout)
	}
}

func isInputFromPipe() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice == 0
}

func getFile() (*os.File, error) {
	if flags.filepath == "" {
		return nil, errors.New("please input a file")
	}
	if !fileExists(flags.filepath) {
		return nil, errors.New("the file provided does not exist")
	}
	file, e := os.Open(flags.filepath)
	if e != nil {
		return nil, errors.Wrapf(e,
			"unable to read the file %s", flags.filepath)
	}
	return file, nil
}

func jsonifyText(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(bufio.NewReader(r))
	re := regexp.MustCompile(`{\"log\".*`)
	for scanner.Scan() {
		line := scanner.Text()
		_, e := fmt.Fprintln(
			w, re.FindAllString(line, -1)[0])
		if e != nil {
			return e
		}
	}
	return nil
}

func fileExists(filepath string) bool {
	info, e := os.Stat(filepath)
	if os.IsNotExist(e) {
		return false
	}
	return !info.IsDir()
}
