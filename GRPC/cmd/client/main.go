package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"mygrpc/pkg/client"
)

func main() {
	// Define CLI flags
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createName := createCmd.String("name", "", "User name for creating a user")
	createEmail := createCmd.String("email", "", "User email for creating a user")

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getId := getCmd.String("id", "", "User ID to retrieve")

	// Parse command
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if *createName == "" || *createEmail == "" {
			fmt.Println("Error: --name and --email flags are required for create")
			os.Exit(1)
		}
		handleCreate(*createName, *createEmail)

	case "get":
		getCmd.Parse(os.Args[2:])
		if *getId == "" {
			fmt.Println("Error: --id flag is required for get")
			os.Exit(1)
		}
		handleGet(*getId)

	case "interactive":
		handleInteractive()

	default:
		printUsage()
		os.Exit(1)
	}
}

func handleCreate(name, email string) {
	fmt.Printf("Creating user with name=%s, email=%s\n", name, email)
	res, err := client.CreateUserRequest(name, email)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("POST Success: ID=%s, Name=%s, Email=%s\n", res.GetId(), res.GetName(), res.GetEmail())
}

func handleGet(id string) {
	fmt.Printf("Fetching user with ID=%s\n", id)
	res, err := client.GetUserRequest(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("GET Success: ID=%s, Name=%s, Email=%s\n", res.GetId(), res.GetName(), res.GetEmail())
}

func handleInteractive() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Interactive mode. Commands: create <name> <email> | get <id> | exit")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "create":
			if len(parts) < 3 {
				fmt.Println("Usage: create <name> <email>")
				continue
			}
			handleCreate(parts[1], parts[2])

		case "get":
			if len(parts) < 2 {
				fmt.Println("Usage: get <id>")
				continue
			}
			handleGet(parts[1])

		case "exit":
			fmt.Println("Goodbye!")
			os.Exit(0)

		default:
			fmt.Println("Unknown command. Use: create <name> <email> | get <id> | exit")
		}
	}
}

func printUsage() {
	fmt.Println(`Usage:
  go run cmd/client/main.go create --name <name> --email <email>
  go run cmd/client/main.go get --id <id>
  go run cmd/client/main.go interactive

Examples:
  go run cmd/client/main.go create --name Alice --email alice@example.com
  go run cmd/client/main.go get --id usr_1
  go run cmd/client/main.go interactive`)
}
