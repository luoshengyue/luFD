package merger

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/pkg/errors"
	"io"
	"luFD/internal/tool"
	"os"
	"sort"
)

func MergeFiles(files []string, out string) error {
	// merge files in order
	sort.Strings(files)

	bar := new(pb.ProgressBar)
	bar.ShowBar = false

	if tool.DisappearProgressBar() {
		fmt.Printf("Start joining \n")
		bar = pb.StartNew(len(files))
	}
	resFile, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, f := range files {
		if err = luFDCopy(f, resFile); err != nil {
			return errors.WithStack(err)
		}
		if tool.DisappearProgressBar() {
			bar.Increment()
		}
	}
	if tool.DisappearProgressBar() {
		bar.Finish()
	}
	return resFile.Close()
}

func luFDCopy(from string, to io.Writer) error {
	f, err := os.OpenFile(from, os.O_RDONLY, 0600)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = io.Copy(to, f)
	if err != nil {
		return errors.WithStack(err)
	}
	return f.Close()
}
