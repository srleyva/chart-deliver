package helpers_test

import (
	"fmt"
	. "github.com/srleyva/chart-deliver/pkg/helpers"
	"os"
	"strings"
	"testing"
)

var template Template = Template{
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

func TestPrintHelmTemplates(t *testing.T) {
	template.Path = ""
	kubernetes, err := template.PrintHelmTemplate()
	if err != nil {
		t.Errorf("err returned where not expected: %s", err)
	}
	defer os.RemoveAll(template.ChartName)

	if !strings.Contains(kubernetes, "chart: tester-v0.0.1") {
		t.Errorf("chart not generated correctly: \n%s", kubernetes)
	}
}

func TestPrintHelmTemplatesValues(t *testing.T) {
	values := `
image:
  repository: tutum/hello-world
`
	file, err := os.Create("values.yaml")
	if err != nil {
		t.Errorf("err creating test value file: %s", err)
	}
	defer file.Close()
	if _, err := file.WriteString(values); err != nil {
		t.Errorf("err creating test value file: %s", err)
	}
	defer os.Remove("values.yaml")
	defer os.RemoveAll(template.ChartName)

	template.Path = ""
	template.Values = "values.yaml"
	kubernetes, err := template.PrintHelmTemplate()
	if err != nil {
		t.Errorf("err returned where not expected: %s", err)
	}

	if !strings.Contains(kubernetes, `image: "tutum/hello-world:stable"`) {
		t.Errorf("chart not generated correctly: \n%s", kubernetes)
	}

}
