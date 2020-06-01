/**
 * @File: cmd
 * @Author: Shaw
 * @Date: 2020/5/31 8:28 PM
 * @Desc

 */

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "GZHU-Pi",
	Short: "GZHU-Pi command",
	Long:  "GZHU-Pi command",
	Run: func(cmd *cobra.Command, args []string) {
		//默认启动api服务
		webApi()
	},
}

var cfgFile string

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.toml", "config file")
}

func Execute() {

	rootCmd.AddCommand(apiCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
