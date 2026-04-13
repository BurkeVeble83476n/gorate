package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/user/gorate/internal/policy"
)

// NewAnnotateCmd returns a cobra command for managing policy annotations.
func NewAnnotateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "annotate",
		Short: "Add, remove, or view annotations on a policy",
	}

	cmd.AddCommand(newAnnotateSetCmd())
	cmd.AddCommand(newAnnotateRemoveCmd())
	cmd.AddCommand(newAnnotateGetCmd())
	return cmd
}

func newAnnotateSetCmd() *cobra.Command {
	var file, name, key, value string
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set an annotation key=value on a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}
			updated, err := policy.Annotate(policies, name, key, value)
			if err != nil {
				return err
			}
			for _, p := range updated {
				if p.Name == name {
					fmt.Fprintf(cmd.OutOrStdout(), "Annotation %q=%q set on policy %q\n", key, value, name)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Policy file (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Policy name (required)")
	cmd.Flags().StringVarP(&key, "key", "k", "", "Annotation key (required)")
	cmd.Flags().StringVarP(&value, "value", "v", "", "Annotation value (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("key")
	_ = cmd.MarkFlagRequired("value")
	return cmd
}

func newAnnotateRemoveCmd() *cobra.Command {
	var file, name, key string
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an annotation key from a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}
			_, err = policy.RemoveAnnotation(policies, name, key)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Annotation %q removed from policy %q\n", key, name)
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Policy file (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Policy name (required)")
	cmd.Flags().StringVarP(&key, "key", "k", "", "Annotation key (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("key")
	return cmd
}

func newAnnotateGetCmd() *cobra.Command {
	var file, name string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display all annotations for a policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies, err := policy.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("loading policies: %w", err)
			}
			anns, err := policy.GetAnnotations(policies, name)
			if err != nil {
				return err
			}
			if len(anns) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No annotations on policy %q\n", name)
				return nil
			}
			for k, v := range anns {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Policy file (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Policy name (required)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
