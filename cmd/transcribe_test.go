package cmd_test

import (
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

	tmpfile, err := os.CreateTemp("../test_audio", "*.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up
	if _, err := tmpfile.Write(dat); err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin
	os.Stdin = tmpfile

	root := cmd.RootCmd
	root.SetIn(nil)
	root.SetArgs([]string{"transcribe"})
	err = root.Execute()
	if err != nil {
		t.FailNow()
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
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
