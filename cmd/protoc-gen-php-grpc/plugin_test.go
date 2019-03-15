package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// When the environment variable RUN_AS_PROTOC_GEN_PHP is set, we skip running
// tests and instead act as protoc-gen-php. This allows the test binary to
// pass itself to protoc.
func init() {
	if os.Getenv("RUN_AS_PROTOC_GEN_PHP") != "" {
		main()
		os.Exit(0)
	}
}

func Test_Simple(t *testing.T) {
	workdir, _ := os.Getwd()
	tmpdir, err := ioutil.TempDir("", "proto-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	args := []string{"-Itestdata", "--php-grpc_out=" + tmpdir}
	args = append(args, "simple/simple.proto")
	protoc(t, args)

	assert.FileExists(t, tmpdir+"/TestSimple/SimpleServiceInterface.php")

	originalFile, err := ioutil.ReadFile(workdir + "/testdata/simple/TestSimple/SimpleServiceInterface.php")
	if err != nil {
		t.Fatal("Can't find original file for comparison")
	}

	generatedFile, err := ioutil.ReadFile(tmpdir + "/TestSimple/SimpleServiceInterface.php")
	if err != nil {
		t.Fatal("Can't find generated file for comparison")
	}

	assert.Equal(t, originalFile, generatedFile)
}

func Test_PhpNamespaceOption(t *testing.T) {
	workdir, _ := os.Getwd()
	tmpdir, err := ioutil.TempDir("", "proto-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	args := []string{"-Itestdata", "--php-grpc_out=" + tmpdir}
	args = append(args, "php_namespace/service.proto")
	protoc(t, args)

	assert.FileExists(t, tmpdir+"/Test/CustomNamespace/ServiceInterface.php")

	originalFile, err := ioutil.ReadFile(workdir + "/testdata/php_namespace/Test/CustomNamespace/ServiceInterface.php")
	if err != nil {
		t.Fatal("Can't find original file for comparison")
	}

	generatedFile, err := ioutil.ReadFile(tmpdir + "/Test/CustomNamespace/ServiceInterface.php")
	if err != nil {
		t.Fatal("Can't find generated file for comparison")
	}

	assert.Equal(t, string(originalFile), string(generatedFile))
}

func Test_UseImportedMessage(t *testing.T) {
	workdir, _ := os.Getwd()
	tmpdir, err := ioutil.TempDir("", "proto-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	args := []string{"-Itestdata", "--php-grpc_out=" + tmpdir}
	args = append(args, "import/service.proto")
	protoc(t, args)

	assert.FileExists(t, tmpdir+"/Import/ServiceInterface.php")

	originalFile, err := ioutil.ReadFile(workdir + "/testdata/import/Import/ServiceInterface.php")
	if err != nil {
		t.Fatal("Can't find original file for comparison")
	}

	generatedFile, err := ioutil.ReadFile(tmpdir + "/Import/ServiceInterface.php")
	if err != nil {
		t.Fatal("Can't find generated file for comparison")
	}

	assert.Equal(t, string(originalFile), string(generatedFile))
}

func protoc(t *testing.T, args []string) {
	cmd := exec.Command("protoc", "--plugin=protoc-gen-php-grpc=" + os.Args[0])
	cmd.Args = append(cmd.Args, args...)
	cmd.Env = append(os.Environ(), "RUN_AS_PROTOC_GEN_PHP=1")
	out, err := cmd.CombinedOutput()
	if len(out) > 0 || err != nil {
		t.Log("RUNNING: ", strings.Join(cmd.Args, " "))
	}
	if len(out) > 0 {
		t.Log(string(out))
	}
	if err != nil {
		t.Fatalf("protoc: %v", err)
	}
}
