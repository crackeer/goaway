package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	signList []string
	rootCmd  = &cobra.Command{
		Use:   "goaway-cli",
		Short: "A generator for Cobra based Applications",
		Long:  `goaway-cli generates`,
		Run:   compileGoaway,
	}
)

func main() {
	rootCmd.PersistentFlags().StringSliceVarP(&signList, "signature", "", []string{}, "sign")
	rootCmd.Execute()
}

func compileGoaway(rootCmd *cobra.Command, args []string) {
	fmt.Println(signList)
}
