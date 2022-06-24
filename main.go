package main

import (
	"fmt"
	"lts/executor"
	"lts/finder"
	"lts/hooks"
	"lts/logging"
	"os"
	"os/signal"
	"syscall"

	"github.com/inancgumus/screen"
)

const VERSION = "v1.2.1"

func main() {
	if len(os.Args) == 0 {
		logging.Errorf("Error: No argument was provided to the program.\n")
		os.Exit(1)
	}
	LTSExecutableName, Args := os.Args[0], os.Args[1:]
	if len(Args) == 0 {
		ShowHelp(LTSExecutableName)
		os.Exit(1)
	}
	CommandName := Args[0]
	// Check if it's internal
	if HandleBuiltInCommand(CommandName) {
		return
	}

	// Read flags
	flags := ReadFlags(Args[1:])

	// Read the command list
	list, err := finder.ReadCommandList()
	if err != nil {
		logging.Errorf("Error: Unable to read command: %v\n", err)
		os.Exit(1)
	}

	// Get command
	command, err := list.GetCommand(CommandName)
	if err != nil {
		logging.Errorf("Error: Unable to execute %v: %v\n", CommandName, err)
		os.Exit(1)
	}

	// Get hooks
	hooksForCommand := []string{}
	if !flags.NoHook {
		hooksForCommand = list.GetHooks(CommandName)
	} else {
		logging.Successf("Ignoring hooks.\n")
	}

	// Execute the command
	var currentKill *func()
	execute := func() (chan int, func(), error) {
		return executor.ExecuteShell(command, func(code int, err error) {
			if err != nil {
				logging.Errorf("%v", err)
				return
			}
			if code != 0 {
				if !flags.NoConsole {
					logging.Errorf("[%v] exit status %v\n", CommandName, code)
				}
				if len(hooksForCommand) == 0 {
					if !flags.IgnoreCode {
						os.Exit(code)
					} else {
						os.Exit(0)
					}
				}
			} else {
				if !flags.NoConsole {
					logging.Successf("Command exited with code 0.\n")
				}
			}
		})
	}

	// Register hooks
	dones := []chan int{}
	if len(hooksForCommand) > 0 {
		screen.Clear()
		screen.MoveTopLeft()
		logging.Successf("Running with hooks: %v\n", hooksForCommand)
	}

	for _, v := range hooksForCommand {
		dones = append(dones, hooks.RegisterHook(v, func() {
			if currentKill != nil {
				fn := *(currentKill)
				go fn()
			}

			_, kill, err := execute()
			*currentKill = kill

			if err != nil {
				logging.Errorf("Error: %v\n", err)
			}
		}))
	}

	// Now, execute it
	done, kill, err := execute()
	currentKill = &kill
	if err != nil {
		logging.Errorf("Error: %v\n", err)
	}

	// Cleanup child when killed
	sigKill := make(chan os.Signal, 3)
	signal.Notify(sigKill, syscall.SIGTERM, syscall.SIGABRT, os.Interrupt)
	go func() {
		<-sigKill
		logging.Errorf("\n\nCleaning up...\n")
		if currentKill != nil {
			fn := *currentKill
			go fn()
		}

		logging.Errorf("Exiting.\n")
		os.Exit(0)
	}()

	hooks.WaitForAllHooks(append(dones, done))
}

func HandleBuiltInCommand(cmd string) bool {
	switch cmd {
	case "list":
		fmt.Printf("Looking for lts.json...\n")
		path, err := finder.FindLTS()
		if err != nil {
			logging.Errorf("LTS lookup failure: %v\n", err)
			os.Exit(1)
			return true
		}
		logging.Infof("The configuration file is found at:\n")
		logging.Infof("%v\n\n", path)
		logging.Infof("Available commands...\n")
		list, err := finder.ReadCommandList()
		if err != nil {
			logging.Errorf("Parse failure: %v\n", err)
			os.Exit(1)
			return true
		}
		for name, exec := range list.Scripts {
			logging.Successf("%v\t", name)
			fmt.Printf("Hooks=")
			hooks := list.GetHooks(name)
			logging.Warnf("%v\t", hooks)
			fmt.Printf("Command=")
			fmt.Printf("%v\n", exec)
		}
		return true
	case "help":
		ShowHelp(os.Args[0])
		return true
	case "version":
		fmt.Printf("%v\n", VERSION)
		return true
	default:
		return false
	}
}

type Flags struct {
	NoHook     bool
	NoConsole  bool
	IgnoreCode bool
}

func ReadFlags(args []string) Flags {
	// Default values
	flags := Flags{
		NoHook:     false,
		NoConsole:  false,
		IgnoreCode: false,
	}
	for _, v := range args {
		switch v {
		case "--no-hook":
			flags.NoHook = true
		case "--no-console":
			flags.NoConsole = true
		case "--ignore-code":
			flags.IgnoreCode = true
		}
	}
	return flags
}

func ShowHelp(ExecutableName string) {
	var InternalCommands map[string]string = map[string]string{
		"list":    "Lists all user-defined commands and the file that defined them.",
		"help":    "Shows this help message.",
		"version": "Shows the version of the program.",
	}

	getSupportedFileExtensionsDescription := func() string {
		description := "Supported extensions: "
		for i, v := range hooks.ELIGIBLE_FILE_EXTENSIONS {
			if i == 0 {
				description += "." + v
			} else {
				description += ", ." + v
			}
		}
		return description
	}

	var Hooks map[string]string = map[string]string{
		"change":   "Listens for file changes in the current directory, and if a change is detected, kill the current command process (and child processes) then rerun the same command. " + getSupportedFileExtensionsDescription(),
		"periodic": "Runs the command periodically, every 30 seconds.",
	}

	getUsageList := func() string {
		description := ""
		i := 0
		for k, _ := range InternalCommands {
			if i == 0 {
				description += k
			} else {
				description += "|" + k
			}
			i++
		}
		return description
	}

	var Flags map[string]string = map[string]string{
		"--no-hook":     "Disables hooks.",
		"--no-console":  "Does not show a console message when the program quits, even if abnormally.",
		"--ignore-code": "Always set the exit code to be zero regardless of what the command exits with.",
	}

	getFlagsList := func() string {
		description := ""
		i := 0
		for k, _ := range Flags {
			if i == 0 {
				description += "[" + k + "]"
			} else {
				description += " [" + k + "]"
			}
			i++
		}
		return description
	}

	logging.Warnf("LTS - Scripts Service (as a part of Lincoln's Tools)\n")
	logging.Warnf("Version: %s\n\n", VERSION)
	fmt.Printf("This program will try and read lts.json from the current directory or the parent directories, then execute the script using shell.\n\n")
	logging.Infof("Usage:\n")
	fmt.Printf("\t%v <name>|%v %v\n\n", ExecutableName, getUsageList(), getFlagsList())
	logging.Infof("Built in commands:\n")
	for k, v := range InternalCommands {
		fmt.Printf("%v:\n", k)
		fmt.Printf("\t%v\n\n", v)
	}
	logging.Infof("Hooks:\n")
	for k, v := range Hooks {
		fmt.Printf("%v:\n", k)
		fmt.Printf("\t%v\n\n", v)
	}
	logging.Infof("Flags:\n")
	for k, v := range Flags {
		fmt.Printf("%v:\n", k)
		fmt.Printf("\t%v\n\n", v)
	}
}
