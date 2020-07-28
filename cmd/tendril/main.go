package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v1"

	"github.com/pete0emerson/spm/pkg/spm"

	"github.com/spf13/cobra"
)

var verbose int

var rootCmd = &cobra.Command{
	Use:              "tendril",
	TraverseChildren: true,
}

type Command struct {
	Short string `yaml:"short"`
	Long  string `yaml:"long"`
}

func setVerbose() {
	if verbose == 0 {
		log.SetLevel(log.FatalLevel)
	} else if verbose == 1 {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}

func loadYAMLFile(filename string) (Command, error) {
	var command Command
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return command, err
	}
	err = yaml.Unmarshal(yamlFile, &command)
	if err != nil {
		return command, fmt.Errorf("Unmarshal: %v", err)
	}
	return command, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getHelp(file string) (shortHelp, longHelp string, err error) {
	if verbose > 0 {
		log.Debugf("Getting help for %s", file)
	}
	helpFile := fmt.Sprintf("%s.yaml", file)
	log.Debugf("Looking for %s", helpFile)
	if fileExists(helpFile) {
		// Load the help file if it exists
		if verbose > 0 {
			log.Debugf("Loading help from %s", helpFile)
		}
		c, err := loadYAMLFile(helpFile)
		if err != nil {
			return "", "", err
		}
		return c.Short, c.Long, nil
	} else {
		// Otherwise we run the script, passing tendril-help to it to get the help yaml
		out, err := exec.Command(file, "tendril-help").Output()
		if err != nil {
			return "", "", err
		}
		var c Command
		err = yaml.Unmarshal(out, &c)
		if err != nil {
			return "", "", err
		}
		return c.Short, c.Long, nil
	}
}

func getDynamicCobraCommands(dir string) map[string]*cobra.Command {
	log.Infof("Loading dynamic commands from %s\n", dir)
	var commands map[string]*cobra.Command
	commands = make(map[string]*cobra.Command)

	fi, err := os.Stat(dir)
	if err != nil {
		return nil
	}
	mode := fi.Mode()
	if !mode.IsDir() {
		return nil
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".yaml") {
			continue
		}
		if strings.Contains(name, ".") {
			nameArray := strings.Split(f.Name(), ".")
			name = strings.Join(nameArray[0:len(nameArray)-1], ".")
		}
		fullPath := dir + "/" + f.Name()
		log.Debugf("Considering %s\n", fullPath)
		if f.IsDir() {
			log.Debugf("Recursing down %s\n", fullPath)
			shortHelp, longHelp, _ := getHelp(fullPath)
			nextLevelCommands := getDynamicCobraCommands(fullPath)
			command := &cobra.Command{
				Use:   name,
				Short: shortHelp,
				Long:  longHelp,
			}
			for _, nextComm := range nextLevelCommands {
				command.AddCommand(nextComm)
			}
			commands[name] = command
		} else {

			fi, err := os.Stat(fullPath)
			if err != nil {
				log.Fatal(err)
			}
			mode := fi.Mode()
			if mode&0111 == 0 {
				log.Debugf("Rejecting non-executable %s\n", fullPath)
				continue
			}

			log.Debugf("Added command: %s\n", name)
			shortHelp, longHelp, err := getHelp(fullPath)
			if err == nil {
				var command = &cobra.Command{
					Use:   name,
					Short: shortHelp,
					Long:  longHelp,
					Run: func(cmd *cobra.Command, args []string) {
						setVerbose()
						log.Infof("Running: %s\n", fullPath)
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
				commands[name] = command
			} else {
				// We may want to fail silently here? Leave it for now.
				log.Fatal(err)
			}
		}
	}

	return commands
}

var operatorForce bool

func main() {
	rootCmd.Flags().CountVarP(&verbose, "verbose", "v", "verbose output")
	rootCmd.Execute()
	setVerbose()

	commands := getDynamicCobraCommands("./tendril")

	for _, cmd := range commands {
		if verbose > 0 {
			log.Infof("Added command: %s\n", cmd.Name())
		}
		rootCmd.AddCommand(cmd)
	}

	var operatorCommand = &cobra.Command{
		Use:   "operator",
		Short: "Tendril operator commands",
		Long:  "Tendril operator commands",
	}

	var operatorInstallCommand = &cobra.Command{
		Use:   "install SOURCE DESTINATION",
		Short: "Install a tendril package",
		Long:  "Install a tendril package",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			setVerbose()
			err := spm.Install(args[0], args[1], operatorForce)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	var operatorRemoveCommand = &cobra.Command{
		Use:   "remove",
		Short: "Remove a tendril package",
		Long:  "Remove a tendril package",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			setVerbose()
			err := spm.Remove(args[0], operatorForce)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	operatorInstallCommand.Flags().BoolVarP(&operatorForce, "force", "f", false, "Force install")
	operatorRemoveCommand.Flags().BoolVarP(&operatorForce, "force", "f", false, "Force remove")
	operatorCommand.AddCommand(operatorInstallCommand)
	operatorCommand.AddCommand(operatorRemoveCommand)
	rootCmd.AddCommand(operatorCommand)
	rootCmd.Execute()
}
