package tool

import (
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func GetIPv4(ips []net.IP) []string {
	var ret = []string{}
	for _, ip := range ips {
		if ip.To4() != nil {
			ret = append(ret, ip.String())
		}
	}
	return ret
}

// DisappearProgressBar check if the current environment is terminal
func DisappearProgressBar() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

func IsFolderExisted(folder string) bool {
	_, err := os.Stat(folder)
	return err == nil
}

// GetFolderFrom prevent path traversal attack.
func GetFolderFrom(url string) (string, error) {
	// path is equal to `echo $HOME/Downloads/luFDTemp`
	path := filepath.Join(os.Getenv("HOME"), SaveFolder)
	// if the url is too long, we only use the first 15 characters,
	// because the file name is too long, it will cause the file to
	// be saved in the wrong folder, and the file type always at the end.
	baseUrl := strings.Split(filepath.Base(url), ".")
	if len(baseUrl[0]) > 15 {
		baseUrl[0] = baseUrl[0][:15]
	}
	// absPath is equal to `echo $HOME/Downloads/luFDTemp/$(basename $url)`
	absPath, err := filepath.Abs(filepath.Join(os.Getenv("HOME"), SaveFolder, filepath.Base(baseUrl[0])))

	if err != nil {
		return "", errors.WithStack(err)
	}
	relative, err := filepath.Rel(path, absPath)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if strings.Contains(relative, "..") {
		return "", errors.WithStack(errors.New("path traversal attack"))
	}
	return absPath, nil
}

func Mkdir(folder string) error {
	if _, err := os.Stat(folder); err != nil {
		if err = os.MkdirAll(folder, 0700); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
