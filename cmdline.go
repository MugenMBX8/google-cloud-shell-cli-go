package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Commands for this program
const (
	CMD_INFO = iota
	CMD_PUTTY
	CMD_EXEC
	CMD_UPLOAD
	CMD_DOWNLOAD
)

func process_cmdline() {
	if len(os.Args) < 2 {
		config.Command = CMD_INFO
		return
	}

	for _, arg := range os.Args {
		if arg == "--help" {
			cmd_help()
			os.Exit(1)
		}
	}

	var args []string

	for x := 1; x < len(os.Args); x++ {
		arg := os.Args[x]

		if arg == "-help" || arg == "--help" {
			cmd_help()
			os.Exit(0)
		}

		if arg == "-debug" || arg == "--debug" {
			config.Debug = true
			continue
		}

		if arg == "-auth" || arg == "--auth" {
			config.Flags.Auth = true
			continue
		}

		if arg == "-login" || arg == "--login" {
			fmt.Println("index:", x)
			fmt.Println("count:", len(os.Args))

			if x == len(os.Args) - 1 {
				fmt.Println("Error: Missing email address to --login")
				os.Exit(1)
			}

			config.Flags.Login = os.Args[x + 1]
			config.Flags.Auth = true
			x++
			continue
		}

		if strings.HasPrefix(arg, "-login=") {
			p := arg[7:]
			config.Flags.Login = p
			config.Flags.Auth = true
			continue
		}

		if strings.HasPrefix(arg, "--login=") {
			p := arg[8:]
			config.Flags.Login = p
			config.Flags.Auth = true
			continue
		}

		args = append(args, arg)
	}

	// fmt.Println("Debug:", config.Debug)
	// fmt.Println("Auth:", config.Flags.Auth)
	// fmt.Println("Login:", config.Flags.Login)

	for x := 0; x < len(args); x++ {
		arg := args[x]

		switch arg {
		case "info":
			config.Command = CMD_INFO

		case "putty":
			config.Command = CMD_PUTTY

		case "exec":
			if len(args) < 2 {
				fmt.Println("Error: expected a remote command")
				os.Exit(1)
			}

			config.Command = CMD_EXEC
			config.RemoteCommand = args[x + 1]
			x++

		case "download":
			if len(args) < 2 {
				fmt.Println("Error: expected a source file name")
				os.Exit(1)
			}

			config.Command = CMD_DOWNLOAD
			config.SrcFile = strings.ReplaceAll(args[x + 1], "\\", "/")
			x++

			if len(args) >= 3 {
				config.DstFile = strings.ReplaceAll(args[x + 1], "\\", "/")
				x++
			} else {
				_, file := path.Split(config.SrcFile)

				config.DstFile = file
			}

			fmt.Println("SrcFile:", config.SrcFile)
			fmt.Println("DstFile:", config.DstFile)

		case "upload":
			if len(args) < 2 {
				fmt.Println("Error: expected a source file name")
				os.Exit(1)
			}

			path, err := filepath.Abs(args[x + 1])

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			config.Command = CMD_UPLOAD
			config.SrcFile = path
			x++

			if len(args) >= 3 {
				config.DstFile = strings.ReplaceAll(args[x + 1], "\\", "/")
				x++
			} else {
				file := strings.ReplaceAll(config.SrcFile, "\\", "/")

				_, file = filepath.Split(file)

				config.DstFile = file
			}

			fmt.Println("SrcFile:", config.SrcFile)
			fmt.Println("DstFile:", config.DstFile)

		default:
			fmt.Println("Error: expected a command (info, putty, exec, upload, download)")
			os.Exit(1)
		}
	}
}

func cmd_help() {
	fmt.Println("Usage: cloudinfo [command]")
	fmt.Println("  cloudinfo                            - display Cloud Shell information")
	fmt.Println("  cloudinfo info                       - display Cloud Shell information")
	fmt.Println("  cloudinfo putty                      - connect to Cloud Shell with Putty")
	fmt.Println("  cloudinfo exec \"command\"             - Execute remote command on Cloud Shell")
	fmt.Println("  cloudinfo upload src_file dst_file   - Upload local file to Cloud Shell")
	fmt.Println("  cloudinfo download src_file dst_file - Download from Cloud Shell to local file")
	fmt.Println("")
	fmt.Println("--debug - Turn on debug output")
	fmt.Println("--auth  - (re)Authenticate ignoring user_credentials.json")
	fmt.Println("--login - Specify an email address as a login hint")
}