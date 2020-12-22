package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

var counter int = 0
var hastes map[string]string = map[string]string{}
var testServer *httptest.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && r.URL.Path == "/documents" {
		hasteBytes, _ := ioutil.ReadAll(r.Body)
		key := fmt.Sprintf("k%d", counter)
		counter++
		hastes[key] = string(hasteBytes)

		w.Write([]byte(fmt.Sprintf(`{"key": "%s"}`, key)))

		return
	}

	getRegexp := regexp.MustCompile(`^/raw/([\w\.]+)`)
	if r.Method == "GET" && getRegexp.MatchString(r.URL.Path) {
		match := getRegexp.FindStringSubmatch(r.URL.Path)
		if strings.Contains(match[1], "/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		haste := hastes[match[1]]
		if haste == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write([]byte(haste))
		return
	}

	w.WriteHeader(http.StatusNotFound)
}))

func TestCreateAndGet(t *testing.T) {
	originalHaste := "This is a test.\n🙃"
	key, err := create(originalHaste, t)
	if err != nil {
		t.Fatalf(`Error creating haste: %s`, err.Error())
	}

	haste, err := get(key, t)
	if err != nil {
		t.Fatalf(`Error reading haste: %s`, err.Error())
	}

	if haste != originalHaste {
		t.Fatalf(`Expected "%s" to be "%s"`, haste, originalHaste)
	}

}

func create(haste string, t *testing.T) (string, error) {
	input := bytes.NewBufferString(haste)
	output := bytes.NewBufferString("")

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"-s", testServer.URL})
	cmd.SetIn(input)
	cmd.SetOut(output)
	cmd.SetErr(nil)

	err := cmd.Execute()
	if err != nil {
		return "", err
	}

	key, err := ioutil.ReadAll(output)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

func get(key string, t *testing.T) (string, error) {
	output := bytes.NewBufferString("")

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"get", key})
	cmd.SetOut(output)
	cmd.SetErr(nil)

	err := cmd.Execute()
	if err != nil {
		return "", err
	}

	haste, err := ioutil.ReadAll(output)
	if err != nil {
		return "", err
	}

	return string(haste), nil
}
