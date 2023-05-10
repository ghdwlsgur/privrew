package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "privrew",
		Short: "...",
		Long:  "...",
	}
)

const (
	rootDir = "/opt/homebrew"
)

var (
	//===/opt/homebrew/Library/Taps
	//===/opt/homebrew/Library/Taps/{Owner}/{Repository}/{File} [-]
	tapDir = func(dir, filePath string) string {
		return dir + filePath
	}(rootDir, "/Library/Taps")

	//===/opt/homebrew/Cellar
	//===/opt/homebrew/Cellar/{Repository}/{Version} [d]
	installDir = func(dir, filePath string) string {
		return dir + filePath
	}(rootDir, "/Cellar")

	//===/opt/homebrew/bin
	//===/opt/homebrew/bin/{Repository} [-]
	symlinkDir = func(dir, filePath string) string {
		return dir + filePath
	}(rootDir, "/bin")
)

func panicRed(err error) {
	fmt.Println(color.RedString("[err] %s", err.Error()))
	os.Exit(1)
}

func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		panicRed(err)
	}
}

func initConfig() {
	args := os.Args[1:]
	_, _, err := rootCmd.Find(args)
	if err != nil {
		panicRed(err)
	}
}
func init() {
	cobra.OnInitialize(initConfig)
}
