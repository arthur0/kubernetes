package config

import (
	"io/ioutil"
	"os"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"bytes"
)


type renameContextTest struct {
	description    string
	config         clientcmdapi.Config //initiate kubectl config
	args           []string            //kubectl rename-context args
	expected       string              //expected out
	expectedConfig clientcmdapi.Config //expect kubectl config
}


func TestRenameContext(t *testing.T) {
	conf := clientcmdapi.Config{
		Kind:       "Config",
		APIVersion: "v1",
		Contexts: map[string]*clientcmdapi.Context{
			"old-name":      {AuthInfo: "auth-info", Cluster: "cluster"},
		},
		CurrentContext: "old-context",
	}
	test := renameContextTest{
		description: "Testing for kubectl config rename-context",
		config:      conf,
		args:        []string{"old-name", "new-name"},
		expected:     "Context \"old-name\" was renamed to \"new-name\".\n",
		expectedConfig: clientcmdapi.Config{
			Kind:       "Config",
			APIVersion: "v1",
			Contexts: map[string]*clientcmdapi.Context{
				"new-name":      {AuthInfo: "auth-info", Cluster: "cluster"},
			},
			CurrentContext: "new-context",
		},
	}
	test.run(t)
}


func (test renameContextTest) run(t *testing.T) {
	fakeKubeFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.Remove(fakeKubeFile.Name())
	err = clientcmd.WriteToFile(test.config, fakeKubeFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pathOptions := clientcmd.NewDefaultPathOptions()
	pathOptions.GlobalFile = fakeKubeFile.Name()
	pathOptions.EnvVar = ""
	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdConfigRenameContext(buf, pathOptions)
	cmd.SetArgs(test.args)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error executing command: %v,kubectl set-context args: %v", err, test.args)
	}

	if len(test.expected) != 0 {
		if buf.String() != test.expected {
			t.Errorf("Failded in:%q\n expected %v\n but got %v", test.description, test.expected, buf.String())
		}
	}
}