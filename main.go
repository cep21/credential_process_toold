package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-ini/ini"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "credentials_process_toold",
		Short:   "Convert a credentials file to a older style credentials file",
		Example: "credentials_process_toold",
		Version: "0.1",
	}
	gen := &generateCommand{}
	rootCmd.AddCommand(gen.Cobra())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type generateCommand struct {
}

func (g *generateCommand) run(cmd *cobra.Command, args []string) (retErr error) {
	inFile, err := cmd.Flags().GetString("in")
	if err != nil {
		return errors.Wrapf(err, "unable to get config for file")
	}
	file, err := ini.Load(inFile)

	if err != nil {
		return errors.Wrapf(err, "unable to load file %s", inFile)
	}
	err = g.generate(file)
	if err != nil {
		return errors.Wrapf(err, "unable to generate config for file %ss", inFile)
	}

	outFilename, err := cmd.Flags().GetString("out")
	if err != nil {
		return errors.Wrapf(err, "unable to get config for file output")
	}
	var out io.Writer
	if outFilename == "-" {
		out = os.Stdout
	} else {
		outFile, err := os.Create(outFilename)
		if err != nil {
			return errors.Wrapf(err, "unable to create file %s", outFilename)
		}
		defer func() {
			if retErr == nil {
				retErr = outFile.Close()
			}
		}()
	}
	_, err = file.WriteTo(out)
	return errors.Wrap(err, "unable to write file output")
}

func (g *generateCommand) Cobra() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "generate the credentials",
		Example: "credentials_process_toold generate",
	}
	cmd.Flags().String("in", defaults.SharedCredentialsFilename(), "Location of the shared credentials file")
	cmd.Flags().String("out", "-", "Filename to write output to.  Use - for stdout")
	cmd.RunE = g.run
	return cmd
}

func (g *generateCommand) generate(file *ini.File) error {
	for _, section := range file.Sections() {
		if section.Name() == "DEFAULT" {
			continue
		}
		ses, err := session.NewSessionWithOptions(session.Options{
			Profile: section.Name(),
		})
		if err != nil {
			return errors.Wrapf(err, "unable to make session for %s", section.Name())
		}
		creds, err := ses.Config.Credentials.Get()
		if err != nil {
			return errors.Wrapf(err, "no credentials provider for %s", section.Name())
		}
		_, err = section.NewKey("aws_access_key_id", creds.AccessKeyID)
		if err != nil {
			return errors.Wrapf(err, "unable to generate section aws_access_key_id for %s", section.Name())
		}
		_, err = section.NewKey("aws_secret_access_key", creds.SecretAccessKey)
		if err != nil {
			return errors.Wrapf(err, "unable to generate section aws_secret_access_key for %s", section.Name())
		}
		_, err = section.NewKey("aws_session_token", creds.SessionToken)
		if err != nil {
			return errors.Wrapf(err, "unable to generate section aws_session_token for %s", section.Name())
		}
	}
	return nil
}