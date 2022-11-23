package cmd

import (
	"github.com/spf13/cobra"

	lp "github.com/rjbrown57/lp/pkg"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate html from yaml",
	Long:  "generate html pages from yaml descriptions",
	Run: func(cmd *cobra.Command, args []string) {
		lp.Lp("generate", genFollow, lpConfig, siteTemplate)
	},
}
