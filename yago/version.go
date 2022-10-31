package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ！！！新版本打 tag 后需要修改此处的 version 值，让 tag 和 version 一致
const Version = "v1.6.5"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  `Print version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("yago version", Version)
	},
}
