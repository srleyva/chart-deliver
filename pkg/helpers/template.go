package helpers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

// Template represent the structure to be templated into metadata for helm
type Template struct {
	ReleaseName string
	ChartName   string
	Version     string
	Path        string
	Values      string
}

// GenerateMetadata will template the metadata into a helm metadata file
func (t *Template) GenerateMetadata() (string, error) {
	tmpl, err := template.New("metadata").Parse(metadata)
	if err != nil {
		return "", err
	}
	buffer := bytes.NewBufferString("")
	err = tmpl.Execute(buffer, t)
	return buffer.String(), err
}

// GenerateHelmChart will generate a whole helm chart
func (t *Template) GenerateHelmChart() error {
	t.Path = fmt.Sprintf("%s/templates", t.ChartName)
	if err := os.MkdirAll(t.Path, 0744); err != nil {
		return err
	}

	chartFiles := []string{"values.yaml", ".helmignore"}
	templateFiles := []string{"NOTES.txt", "_helpers.tpl", "deployment.yaml", "ingress.yaml", "service.yaml"}

	// TODO make DRY...lots of duplicate logic
	for _, file := range chartFiles {
		filePath := fmt.Sprintf("%s/%s", t.ChartName, file)
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		switch file {
		case "values.yaml":
			f.WriteString(defaultValues)
		case ".helmignore":
			f.WriteString(helmIgnore)
		}
		f.Close()
	}

	for _, file := range templateFiles {
		filePath := fmt.Sprintf("%s/%s", t.Path, file)
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		switch file {
		case "NOTES.txt":
			if _, err := f.WriteString(defaultNotes); err != nil {
				return err
			}
		case "_helpers.tpl":
			if _, err := f.WriteString(defaultHelpers); err != nil {
				return err
			}
		case "deployment.yaml":
			if _, err := f.WriteString(defaultDeployment); err != nil {
				return err
			}
		case "ingress.yaml":
			if _, err := f.WriteString(defaultIngress); err != nil {
				return err
			}
		case "service.yaml":
			if _, err := f.WriteString(defaultService); err != nil {
				return err
			}
		}
		f.Close()
	}

	metadataFile, err := os.Create(fmt.Sprintf("%s/Chart.yaml", t.ChartName))
	if err != nil {
		return err
	}

	metadata, err := t.GenerateMetadata()
	if err != nil {
		return err
	}
	if _, err := metadataFile.WriteString(metadata); err != nil {
		return err
	}
	return nil
}

// PrintHelmTemplate will print the helm template to stdout
func (t *Template) PrintHelmTemplate() (string, error) {
	if t.Path == "" {
		if err := t.GenerateHelmChart(); err != nil {
			return "", err
		}
	}
	cmd := exec.Command("helm", "template", t.ChartName)
	if t.Values != "" {
		cmd.Args = append(cmd.Args, "--values", t.Values)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}
