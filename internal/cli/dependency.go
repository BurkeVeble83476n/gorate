package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gorate/internal/policy"
)

// NewDependencyCmd returns the root dependency command with subcommands.
func NewDependencyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dependency",
		Short: "Manage policy dependencies",
	}
	cmd.AddCommand(newDependencyAddCmd())
	cmd.AddCommand(newDependencyGetCmd())
	cmd.AddCommand(newDependencyGraphCmd())
	return cmd
}

func newDependencyAddCmd() *cobra.Command {
	var file, from, to string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a dependency from one policy to another",
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file flag is required")
			}
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return err
			}
			updated, err := policy.AddDependency(policies, from, to)
			if err != nil {
				return err
			}
			return writeFile(file, updated)
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "Path to policy YAML file")
	cmd.Flags().StringVar(&from, "from", "", "Name of the dependent policy")
	cmd.Flags().StringVar(&to, "to", "", "Name of the policy to depend on")
	return cmd
}

func newDependencyGetCmd() *cobra.Command {
	var file, name string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get dependencies for a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file flag is required")
			}
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return err
			}
			deps, err := policy.GetDependencies(policies, name)
			if err != nil {
				return err
			}
			if len(deps) == 0 {
				fmt.Println("no dependencies")
				return nil
			}
			for _, d := range deps {
				fmt.Println(d)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "Path to policy YAML file")
	cmd.Flags().StringVar(&name, "name", "", "Policy name")
	return cmd
}

func newDependencyGraphCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Output the dependency graph as JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file flag is required")
			}
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return err
			}
			graph := policy.BuildGraph(policies)
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(graph)
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "Path to policy YAML file")
	return cmd
}
