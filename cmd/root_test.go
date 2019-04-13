package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"testing"
)

func TestExecute(t *testing.T) {
	cmd := &cobra.Command{Run: func(cmd *cobra.Command, args []string) {}}
	var a int64
	cmd.Flags().Int64Var(&a, "test", 999, "default number is 999")
	cmd.Flags().Set("test", "7")
	cmd.Run(cmd, []string{""})
	i, err := cmd.Flags().GetInt64("test")
	fmt.Printf("value is %d %v", i, err)

}

func TestExecute2(t *testing.T) {
	rootCmd.Flags().Set("timeout", "5")
	rootCmd.Execute()
}
