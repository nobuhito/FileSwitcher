package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CfgFile string
var Target string
var TargetArg string

func main() {

	var RootCmd = &cobra.Command{
		Use:  "FileSwitcher",
		Long: `FileSwitcher can switch the file of the specified path with hard link.`,
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize",
		Long:  `Prepare to create a hard link of the target file in the current directory and start using it.`,
		Run: func(cmd *cobra.Command, args []string) {
			if isFirstUse(Target, "./") {
				ext := filepath.Ext(Target)
				err := os.Link(Target, "default"+ext)
				if err != nil {
					fmt.Println(err)
				}
			}
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List the target file",
		Long:  `List files to be switched. Listed files can be used with the set command.`,
		Run: func(cmd *cobra.Command, args []string) {
			ext := filepath.Ext(Target)

			fileInfos, err := ioutil.ReadDir(filepath.FromSlash("./"))
			if err != nil {
				fmt.Println(err)
			}

			for _, fileInfo := range fileInfos {
				name := fileInfo.Name()
				if isSameExt(ext, name) {
					if isSameFile(Target, name) {
						fmt.Print(" * ")
					} else {
						fmt.Print("   ")
					}
					fmt.Printf("%s\n", strings.Replace(name, ext, "", 1))
				}
			}
		},
	}

	var setCmd = &cobra.Command{
		Use:   "set",
		Short: "Switch the target file",
		Long:  `Switch the target file to the specified file. Please specify the file displayed by list command.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {

				ext := filepath.Ext(Target)

				if hasOrginalFile(Target, "./") {
					err := os.Remove(Target)
					if err != nil {
						fmt.Println(err)
					}

					err = os.Link(args[0]+ext, Target)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		},
	}

	RootCmd.PersistentFlags().StringVarP(&CfgFile, "config", "c", "FileSwitcher.yaml", "config file")
	RootCmd.PersistentFlags().StringVarP(&TargetArg, "target", "t", "", "target file")

	cobra.OnInitialize(initConfig)

	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(setCmd)

	RootCmd.Execute()
}

func initConfig() {
	_, err := os.Stat(CfgFile)
	if err == nil {

		viper.SetConfigFile(CfgFile)
		// viper.SetConfigName("FileSwitcher") // name of config file (without extension)
		viper.AddConfigPath(".") // adding home directory as first search path
		viper.AutomaticEnv()     // read in environment variables that match

		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		_target := viper.GetString("target")
		if _target != "" {
			Target = _target
		}
	}

	if TargetArg != "" {
		Target = TargetArg
	}

	Target = filepath.FromSlash(Target)

	if Target == "" {
		fmt.Println("target file is not set.")
		os.Exit(1)
	}
}

func isFirstUse(path string, dir string) bool {
	dirname := normarizeDir(dir)
	ext := filepath.Ext(path)
	fileInfos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return false
	}

	ret := true
	for _, fileInfo := range fileInfos {
		name := dirname + fileInfo.Name()
		if isSameExt(ext, name) && isSameFile(path, name) {
			ret = false
		}
	}

	return ret
}

func hasOrginalFile(Target string, dir string) bool {
	dirname := normarizeDir(dir)
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	for _, fileInfo := range fileInfos {
		name := dirname + fileInfo.Name()
		if isSameFile(Target, name) {
			return true
		}
	}

	return false
}

func isSameFile(source, dist string) bool {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return false
	}

	distInfo, err := os.Stat(dist)
	if err != nil {
		return false
	}
	return os.SameFile(sourceInfo, distInfo)
}

func isSameExt(_ext, filename string) bool {
	ext := filepath.Ext(filename)
	if ext == _ext {
		return true
	} else {
		return false
	}
}

func normarizeDir(dir string) string {
	re, err := regexp.Compile("/$")
	if err != nil {
		fmt.Println("could not compile regexp")
	} else {
		dir = re.ReplaceAllString(dir, "") + "/"
	}
	return filepath.FromSlash(dir)
}
