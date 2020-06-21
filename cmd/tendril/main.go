package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var verbose int

var rootCmd = &cobra.Command{
	Use:              "tendril",
	TraverseChildren: true,
}

func getHelp(file string) (shortHelp, longHelp string, err error) {
	if verbose > 0 {
		log.Printf("Getting help for %s\n", file)
	}
	shortHelp = "short"
	longHelp = "long"
	err = nil
	return
}

func getDynamicCobraCommands(dir string) map[string]*cobra.Command {
	if verbose > 0 {
		log.Printf("Loading dynamic commands from %s\n", dir)
	}
	var commands map[string]*cobra.Command
	commands = make(map[string]*cobra.Command)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		name := f.Name()
		if strings.Contains(name, ".") {
			nameArray := strings.Split(f.Name(), ".")
			name = strings.Join(nameArray[0:len(nameArray)-1], ".")
		}
		fullPath := dir + "/" + f.Name()
		if verbose > 0 {
			log.Infof("Considering %s\n", fullPath)
		}
		if f.IsDir() {
			if verbose > 0 {
				log.Infof("Recursing down %s\n", fullPath)
			}
			nextLevelCommands := getDynamicCobraCommands(fullPath)
			command := &cobra.Command{
				Use: name,
			}
			for _, nextComm := range nextLevelCommands {
				command.AddCommand(nextComm)
			}
			commands[name] = command
		} else {
			if verbose > 0 {
				log.Infof("Added command: %s\n", name)
			}
			shortHelp, longHelp, err := getHelp(fullPath)
			if err == nil {
				var command = &cobra.Command{
					Use:   name,
					Short: shortHelp,
					Long:  longHelp,
					Run: func(cmd *cobra.Command, args []string) {
						if verbose > 0 {
							fmt.Printf("Running: %s\n", fullPath)
						}
						c := exec.Command(fullPath, strings.Join(args, " "))
						c.Stdout = os.Stdout
						c.Stderr = os.Stderr
						// var waitStatus syscall.WaitStatus

						if err := c.Run(); err != nil {
							// Did the command fail because of an unsuccessful exit code
							if exitError, ok := err.(*exec.ExitError); ok {
								waitStatus := exitError.Sys().(syscall.WaitStatus)
								os.Exit(waitStatus.ExitStatus())
							}
						}

					},
				}
				if verbose > 1 {
					fmt.Printf("commands[%s] = %#v\n", name, command)
				}
				commands[name] = command
			}
		}
	}

	return commands
}

func main() {
	rootCmd.Flags().CountVarP(&verbose, "verbose", "v", "verbose output")
	rootCmd.Execute()

	commands := getDynamicCobraCommands("./tendril/commands")

	for _, cmd := range commands {
		if verbose > 0 {
			log.Infof("Added command: %s\n", cmd.Name())
		}
		rootCmd.AddCommand(cmd)
	}
	rootCmd.Execute()

}
