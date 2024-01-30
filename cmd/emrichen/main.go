package main

import (
	"context"
	"github.com/Masterminds/sprig"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/help"
	"github.com/go-go-golems/glazed/pkg/helpers/cast"
	"github.com/go-go-golems/go-emrichen/pkg/doc"
	"github.com/go-go-golems/go-emrichen/pkg/emrichen"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type ProcessCommand struct {
	*cmds.CommandDescription
}

var _ cmds.WriterCommand = (*ProcessCommand)(nil)

type ProcessSettings struct {
	InputFiles   []*parameters.FileData `glazed.parameter:"input-files"`
	VarFile      []*parameters.FileData `glazed.parameter:"var-file"`
	Output       string                 `glazed.parameter:"output"`
	OutputFormat string                 `glazed.parameter:"output-format"`
	IncludeEnv   bool                   `glazed.parameter:"include-env"`
	Define       map[string]string      `glazed.parameter:"define"`
}

func NewProcessCommand() (*ProcessCommand, error) {
	return &ProcessCommand{
		CommandDescription: cmds.NewCommandDescription(
			"process",
			cmds.WithShort("Process files with emrichen"),
			cmds.WithArguments(
				parameters.NewParameterDefinition(
					"input-files",
					parameters.ParameterTypeFileList,
					parameters.WithHelp("Input files to process"),
					parameters.WithRequired(true),
				),
			),
			cmds.WithFlags(
				parameters.NewParameterDefinition(
					"var-file",
					parameters.ParameterTypeFileList,
					parameters.WithHelp("File list to process"),
					parameters.WithShortFlag("f"),
				),
				parameters.NewParameterDefinition(
					"output",
					parameters.ParameterTypeString,
					parameters.WithHelp("Output file"),
					parameters.WithShortFlag("o"),
				),
				parameters.NewParameterDefinition(
					"output-format",
					parameters.ParameterTypeChoice,
					parameters.WithHelp("Output format (json, yaml, pprint)"),
					parameters.WithChoices("json", "yaml", "pprint"),
				),
				parameters.NewParameterDefinition(
					"include-env",
					parameters.ParameterTypeBool,
					parameters.WithHelp("Include environment variables"),
					parameters.WithShortFlag("e"),
				),
				parameters.NewParameterDefinition(
					"define",
					parameters.ParameterTypeKeyValue,
					parameters.WithHelp("Define key-value variables"),
					parameters.WithShortFlag("D"),
				),
			),
		),
	}, nil
}

func (c *ProcessCommand) RunIntoWriter(
	ctx context.Context,
	ps *layers.ParsedLayers,
	w io.Writer,
) error {
	s := &ProcessSettings{}
	if err := ps.InitializeStruct(layers.DefaultSlug, s); err != nil {
		return err
	}

	env := map[string]interface{}{}

	for _, file := range s.VarFile {
		// if the content is a list of objects, we want to merge them into the environment
		if objs, ok := cast.CastList2[map[string]interface{}, interface{}](file.ParsedContent); ok {
			for _, obj := range objs {
				for k, v := range obj {
					env[k] = v
				}
			}
			continue
		}

		obj, ok := file.ParsedContent.(map[string]interface{})
		if ok {
			for k, v := range obj {
				env[k] = v
			}
			continue
		}

		return errors.Errorf("could not cast %s to map[string]interface{}", file.Path)
	}

	ei, err := emrichen.NewInterpreter(emrichen.WithVars(env),
		emrichen.WithFuncMap(sprig.TxtFuncMap()))
	if err != nil {
		return err
	}

	for _, file := range s.InputFiles {
		err := processFile(ei, file.Path, w)
		if err != nil {
			return err
		}
	}

	return nil
}

func processFile(interpreter *emrichen.Interpreter, filePath string, w io.Writer) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(f io.Closer) {
		_ = f.Close()
	}(f)

	decoder := yaml.NewDecoder(f)

	docCount := 0
	for {
		var document interface{}

		err = decoder.Decode(interpreter.CreateDecoder(&document))
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// skip a document that was probably used to set !Defaults
		if document == nil {
			continue
		}

		processedYAML, err := yaml.Marshal(&document)
		if err != nil {
			return err
		}

		if docCount > 0 {
			_, err = w.Write([]byte("---\n"))
			if err != nil {
				return err
			}
		}

		_, err = w.Write(processedYAML)

		docCount += 1
	}

	return nil
}

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "emrichen",
	Short: "Emrichen is a YAML preprocessor",
}

func main() {
	helpSystem := help.NewHelpSystem()
	err := doc.AddDocToHelpSystem(helpSystem)
	cobra.CheckErr(err)
	helpSystem.SetupCobraRootCommand(rootCmd)

	processCmd, err := NewProcessCommand()
	cobra.CheckErr(err)
	processCommand, err := cli.BuildCobraCommandFromWriterCommand(processCmd)
	cobra.CheckErr(err)

	rootCmd.AddCommand(processCommand)

	err = rootCmd.Execute()
	cobra.CheckErr(err)
}
