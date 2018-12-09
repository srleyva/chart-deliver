package helpers

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// Template represent the structure to be templated into metadata for helm
type Template struct {
	Runner      Runner
	ReleaseName string
	ChartName   string
	Version     string
	Path        string
	Values      string
	Image       string
	Tag         string
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
	args := []string{"template", t.ChartName}
	args = t.buildArgs(args)
	out, err := t.Runner.Run("helm", args...)
	return string(out), err
}

// InstallTemplate Actually will install into the k8s cluster
func (t *Template) InstallTemplate() (string, error) {
	if t.Path == "" {
		if err := t.GenerateHelmChart(); err != nil {
			return "", err
		}
	}
	args := []string{"upgrade", "--install", t.ReleaseName, t.ChartName}
	args = t.buildArgs(args)
	out, err := t.Runner.Run("helm", args...)
	return string(out), err
}

func (t *Template) buildArgs(args []string) []string {
	if t.Values != "" {
		args = append(args, "--values", t.Values)
	}
	if t.Image != "" {
		args = append(args, "--set", fmt.Sprintf("image.repository=%s", t.Image))
	}
	if t.Tag != "" {
		args = append(args, "--set", fmt.Sprintf("image.tag=%s", t.Tag))
	}
	return args
}
