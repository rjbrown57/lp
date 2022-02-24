package cmd

import (
	"github.com/spf13/cobra"

	lp "github.com/rjbrown57/lp/pkg"
)

var cfgFile string
var lpConfig, siteTemplate string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lp",
	Short: "Host a yaml based static link page",
	Long:  `A yaml based static link page for every day work use.`,
	Run: func(cmd *cobra.Command, args []string) {
		lp.Lp(lpConfig, siteTemplate)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&lpConfig, "lpConfig", "l", "config/lp.yaml", "base config for lp see https://github.com/rjbrown57/lp/blob/main/config/lp.yaml")
	rootCmd.PersistentFlags().StringVarP(&siteTemplate, "siteTemplate", "s", "config/site.yaml", "site template see https://github.com/rjbrown57/lp/blob/main/config/site.yaml")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
