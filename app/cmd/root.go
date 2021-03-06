package cmd

import (
	"fmt"
	"github.com/AndresBott/goback/app/goback"
	"github.com/AndresBott/goback/internal/clilog"
	"github.com/AndresBott/goback/internal/profile"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

// Execute is the entry point for the command line
func Execute() {
	if err := newRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

var (
	ShaVer    string // sha1 revision used to build the program
	BuildTime string // when the executable was built
	Version   = "development"
)

func newRootCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "goback",
		Short: "goback is a simple backup tool",
		Long:  "goback is an opinionated tool to backup data from different sources like file system directories, mysql databases",
	}

	profileFlag := ""
	dirflag := ""
	generateflag := false
	version := false
	quiet := false

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		printer := clilog.CliOut{
			Quiet:  quiet,
			Prefix: "-> ",
		}
		runner := goback.BackupRunner{
			Printer: printer,
		}

		if version {
			fmt.Printf("Version: %s\n", Version)
			fmt.Printf("Build date: %s\n", BuildTime)
			fmt.Printf("Commit sha: %s\n", ShaVer)
			fmt.Printf("Compiler: %s\n", runtime.Version())

			return nil
		} else if profileFlag != "" {
			// handle a single profile
			err := runner.LoadProfile(profileFlag)
			if err != nil {
				return err
			}

			err = runner.Run()
			if err != nil {
				return err
			}

			return nil

		} else if dirflag != "" {
			// handle a directory containing profiles
			err := runner.LoadProfiles(dirflag)
			if err != nil {
				return err
			}

			err = runner.Run()
			if err != nil {
				return err
			}

			return nil
		} else if generateflag {
			// print the profile boilerplate to stdout
			fmt.Println(profile.Boilerplate())
			return nil
		}

		_ = cmd.Help()
		return nil
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVarP(&profileFlag, "profile", "p", "", "single profile file to run")
	cmd.Flags().StringVarP(&dirflag, "dir", "d", "", "directory containing multiple profiles")
	cmd.Flags().BoolVarP(&generateflag, "generate", "g", false, "print a profile boilerplate")
	cmd.Flags().BoolVarP(&version, "version", "v", false, "print version and build information")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "do not print output")

	cmd.SilenceUsage = true
	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		_ = cmd.Help()
		return nil
	})

	return cmd
}
