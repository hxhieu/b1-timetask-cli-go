package console

import (
	"io"
	"os"
	"strings"
	"testing"
)

type printFn = func(msg string)

func captureStdOut(fn printFn, msg string) string {
	storeStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn(msg)

	w.Close()
	out, _ := io.ReadAll(r)
	// restore the stdout
	os.Stdout = storeStdout

	return strings.TrimSpace(string(out))
}

func TestError(t *testing.T) {
	want := "Hello world!!!"
	got := captureStdOut(Error, want)
	if want != got {
		t.Errorf("want '%s', but got '%s'", want, got)
	}
}

func TestErrorLn(t *testing.T) {
	want := "Hello world!!!"
	got := captureStdOut(ErrorLn, want)
	if want != got {
		t.Errorf("want '%s', but got '%s'", want, got)
	}
}

func TestSuccess(t *testing.T) {
	want := "Hello world!!!"
	got := captureStdOut(Success, want)
	if want != got {
		t.Errorf("want '%s', but got '%s'", want, got)
	}
}

func TestSuccessLn(t *testing.T) {
	want := "Hello world!!!"
	got := captureStdOut(SuccessLn, want)
	if want != got {
		t.Errorf("want '%s', but got '%s'", want, got)
	}
}

func TestInfo(t *testing.T) {
	want := "Hello world!!!"
	got := captureStdOut(Info, want)
	if want != got {
		t.Errorf("want '%s', but got '%s'", want, got)
	}
}

func TestInfoLn(t *testing.T) {
	want := "Hello world!!!"
	got := captureStdOut(InfoLn, want)
	if want != got {
		t.Errorf("want '%s', but got '%s'", want, got)
	}
}

func TestWarn(t *testing.T) {
	want := "Hello world!!!"
	got := captureStdOut(Warn, want)
	if want != got {
		t.Errorf("want '%s', but got '%s'", want, got)
	}
}

func TestWarnLn(t *testing.T) {
	want := "Hello world!!!"
	got := captureStdOut(WarnLn, want)
	if want != got {
		t.Errorf("want '%s', but got '%s'", want, got)
	}
}

func TestHeader(t *testing.T) {
	want := "Hello world!!!"
	got := captureStdOut(Header, want)
	if want != got {
		t.Errorf("want '%s', but got '%s'", want, got)
	}
}
