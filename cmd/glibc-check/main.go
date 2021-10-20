package main

import (
	"fmt"
	"os"

	glibccheck "github.com/guseggert/glibc-check"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "glibc-check",
	Short: "tools for checking glibc versions on an executable",
}

func main() {
	rootCmd.AddCommand(&cobra.Command{
		Use:           "list-versions",
		Short:         "list the glibc versions",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			versions, err := glibccheck.ParseFile(args[0])
			if err != nil {
				return err
			}
			for _, v := range versions {
				fmt.Println(v.String())
			}
			return nil
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:           "min filename",
		Short:         "print the min glibc version",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			versions, err := glibccheck.ParseFile(args[0])
			if err != nil {
				return err
			}
			if len(versions) == 0 {
				fmt.Fprintln(os.Stderr, "no glibc versions found")
				os.Exit(1)
			}
			fmt.Println(versions[0])
			return nil
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:           "max filename",
		Short:         "print the max glibc version",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			versions, err := glibccheck.ParseFile(args[0])
			if err != nil {
				return err
			}
			if len(versions) == 0 {
				fmt.Fprintln(os.Stderr, "no glibc versions found")
				os.Exit(1)
			}
			fmt.Println(versions[len(versions)-1])
			return nil
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:           "assert-all condition filename",
		Short:         "asserts that the given condition holds for every glibc version in the file",
		Long:          "The variables 'major', 'minor', and 'patch' are bound for the current version, 'patch' is bound to 0 for versions without patch versions, such as '2.32'. Example: 'glibc-check all-versions 'major == 2 && minor > 10'",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			versions, err := glibccheck.ParseFile(args[1])
			if err != nil {
				return err
			}
			if len(versions) == 0 {
				fmt.Fprintln(os.Stderr, "no glibc versions found")
				os.Exit(1)
			}
			violations, err := versions.FindViolations(args[0])
			if len(violations) > 0 {
				fmt.Fprintf(os.Stderr, "condition did not hold for versions: %s\n", violations)
				os.Exit(len(violations))
			}
			return nil
		},
	})
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
