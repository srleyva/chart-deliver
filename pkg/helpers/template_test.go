package helpers_test

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	. "github.com/srleyva/chart-deliver/pkg/helpers"
)

// MockRunner is designed to test that helm is called as expected
type MockRunner struct {
	cmd  string
	args []string
}

func (r *MockRunner) Run(command string, args ...string) ([]byte, error) {
	r.cmd = command
	r.args = args
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	out, err := cmd.CombinedOutput()
	return out, err
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	fmt.Println("testing helper process")
}

var mockRunner *MockRunner = &MockRunner{}

var template Template = Template{
	Runner:      mockRunner,
	ReleaseName: "test",
	ChartName:   "tester",
	Version:     "v0.0.1",
}

func TestGenerateMetadata(t *testing.T) {
	meta, err := template.GenerateMetadata()
	if err != nil {
		t.Errorf("err returned where not expected: %s", err)
	}

	// TODO Make this dry using table tests
	if !strings.Contains(meta, "name: tester") {
		t.Errorf("Name not templated correctly: \n%s", meta)
	}

	if !strings.Contains(meta, "version: v0.0.1") {
		t.Errorf("Version not templated correctly: \n%s", meta)
	}

	if !strings.Contains(meta, "appVersion: test") {
		t.Errorf("ReleaseName not templated correctly: \n%s", meta)
	}
}

func TestGenerateHelmChart(t *testing.T) {
	if err := template.GenerateHelmChart(); err != nil {
		t.Errorf("err returned where not expected: %s", err)
	}
	defer os.RemoveAll(template.ChartName)

	// TODO test correction of yaml files
	if _, err := os.Stat(fmt.Sprintf("%s/templates", template.ChartName)); os.IsNotExist(err) {
		t.Errorf("directories were not created")
	}
}

func TestPrintTemplate(t *testing.T) {
	defer os.RemoveAll(template.ChartName)

	t.Run("test print with no values file passed", func(t *testing.T) {
		template.Path = ""
		_, err := template.PrintHelmTemplate()
		if err != nil {
			t.Errorf("err returned where not expected: %s", err)
		}

		expectedArgs := []string{"template", "tester"}

		if mockRunner.cmd != "helm" {
			t.Errorf("Expected Call: helm \n Actual Call: %s", mockRunner.cmd)
		}

		if !reflect.DeepEqual(expectedArgs, mockRunner.args) {
			t.Errorf("wrong args passed: \n Expected: %s\n Actual: %s\n", expectedArgs, mockRunner.args)
		}
	})

	t.Run("test print with values passed", func(t *testing.T) {
		template.Path = ""
		template.Values = "values.yaml"
		_, err := template.PrintHelmTemplate()
		if err != nil {
			t.Errorf("err returned where not expected: %s", err)
		}

		expectedArgs := []string{"template", "tester", "--values", "values.yaml"}

		if mockRunner.cmd != "helm" {
			t.Errorf("Expected Call: helm \n Actual Call: %s", mockRunner.cmd)
		}

		if !reflect.DeepEqual(expectedArgs, mockRunner.args) {
			t.Errorf("wrong args passed: \n Expected: %s\n Actual: %s\n", expectedArgs, mockRunner.args)
		}
	})

	t.Run("test print with image and tag passed", func(t *testing.T) {
		template.Path = ""
		template.Values = ""
		template.Image = "tutum/hello-world"
		template.Tag = "latest"
		_, err := template.PrintHelmTemplate()
		if err != nil {
			t.Errorf("err returned where not expected: %s", err)
		}

		expectedArgs := []string{
			"template",
			"tester",
			"--set",
			"image.repository=tutum/hello-world",
			"--set",
			"image.tag=latest"}

		if mockRunner.cmd != "helm" {
			t.Errorf("Expected Call: helm \n Actual Call: %s", mockRunner.cmd)
		}

		if !reflect.DeepEqual(expectedArgs, mockRunner.args) {
			t.Errorf("wrong args passed: \n Expected: %s\n Actual: %s\n", expectedArgs, mockRunner.args)
		}
		template.Tag = ""
		template.Image = ""
	})
}

func TestInstallTemplate(t *testing.T) {
	t.Run("test install with no values file", func(t *testing.T) {
		template.Path = ""
		template.Values = ""
		_, err := template.InstallTemplate()
		if err != nil {
			t.Errorf("err returned where not expected: %s", err)
		}

		expectedArgs := []string{"upgrade", "--install", "test", "tester"}

		if mockRunner.cmd != "helm" {
			t.Errorf("Expected Call: helm \n Actual Call: %s", mockRunner.cmd)
		}

		if !reflect.DeepEqual(expectedArgs, mockRunner.args) {
			t.Errorf("wrong args passed: \n Expected: %s\n Actual: %s\n", expectedArgs, mockRunner.args)
		}
	})

	t.Run("test install with a values file", func(t *testing.T) {
		template.Path = ""
		template.Values = "values.yaml"
		_, err := template.InstallTemplate()
		if err != nil {
			t.Errorf("err returned where not expected: %s", err)
		}

		expectedArgs := []string{"upgrade", "--install", "test", "tester", "--values", "values.yaml"}

		if mockRunner.cmd != "helm" {
			t.Errorf("Expected Call: helm \n Actual Call: %s", mockRunner.cmd)
		}

		if !reflect.DeepEqual(expectedArgs, mockRunner.args) {
			t.Errorf("wrong args passed: \n Expected: %s\n Actual: %s\n", expectedArgs, mockRunner.args)
		}
	})
}
