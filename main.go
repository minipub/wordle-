package main

import (
	"fmt"
	"os"

	puzzle "github.com/minipub/wordle/cmd/wordle-puzzle"
	solver "github.com/minipub/wordle/cmd/wordle-solver"

	"github.com/spf13/cobra"
)

var (
	version bool
	root    cobra.Command
)

func main() {
	root = cobra.Command{
		Use: "wordle",
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				fmt.Println("Wordle version:", "v0.1.5")
				os.Exit(0)
			}
			help()
			os.Exit(0)
		},
	}

	root.AddCommand(
		puzzle.Cmd,
		solver.Cmd,
	)

	root.Flags().BoolVarP(&version, "version", "V", false, "Wordle version")

	root.Execute()
}

func help() {
	root.Help()
}
