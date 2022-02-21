/*
Copyright © 2022 rjbrown57

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
	rootCmd.PersistentFlags().StringVarP(&siteTemplate, "siteTempalte", "s", "config/site.yaml", "site tempalte see https://github.com/rjbrown57/lp/blob/main/config/site.yaml")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
