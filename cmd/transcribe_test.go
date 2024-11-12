package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/prattlOrg/prattl/cmd"
)

var audio00 = filepath.Clean("../test_audio/test.mp3")
var audio01 = filepath.Clean("../test_audio/transcribing_1.mp3")

func TestStdin(t *testing.T) {
	dat, err := os.ReadFile(audio00)
	if err != nil {
		t.Fatal(err)
	}
	// t.Log(dat)

	root := cmd.RootCmd
	in := bytes.NewBuffer(dat)
	root.SetIn(in)
	root.SetArgs([]string{"transcribe"})
	err = root.Execute()
	if err != nil {
		t.FailNow()
	}
}

func TestFp(t *testing.T) {
	root := cmd.RootCmd
	root.SetArgs([]string{"transcribe", audio00, audio01})
	err := root.Execute()
	if err != nil {
		t.FailNow()
	}
}
