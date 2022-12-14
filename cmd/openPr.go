/*
Copyright © 2022 Maurício Mussatto Scopel <ms.mauricio93@gmail.com>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
)

// openPrCmd represents the openPr command
var openPrCmd = &cobra.Command{
	Use:   "openPr",
	Short: "Open a Pull Request",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("openPr called")
		output, _ := exec.Command("ls", ".").Output()
		fmt.Println(string(output))
		fmt.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(openPrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openPrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openPrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}