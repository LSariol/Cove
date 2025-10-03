package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/LSariol/Cove/internal/encryption"
)

func StartCLI() {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Cove CLI> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		parseCLI(strings.Fields(input))
	}
}

func parseCLI(args []string) {

	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "exit", "quit":
		fmt.Println("Shutting down Cove...")
		os.Exit(0)

	case "get", "g":

		if len(args) != 2 {
			yellowLog("Get requires 2 total arguments.")
			yellowLog("get <secret>")
			return
		}

		res, ok := encryption.GetSecret(args[1])
		if !ok {
			redLog(res + ": " + args[1])
			return
		}

		greenLog("Secret has been retreived: " + res)

	case "add", "a":

		if len(args) != 3 {
			yellowLog("add requires 3 total arguments.")
			yellowLog("add <secretName> <value>")
			return
		}

		res, ok := encryption.AddSecret(args[1], args[2])
		if !ok {
			redLog(res)
		}

		greenLog("Secret has been added")

	case "remove", "r", "delete", "d":

		if len(args) != 2 {
			yellowLog("Get requires 2 total arguments.")
			yellowLog("remove <secret>")
			return
		}

		res, ok := encryption.RemoveSecret(args[1])
		if !ok {
			redLog(res)
			return
		}

		greenLog("Secret has been removed")

	case "update", "u":

		if len(args) != 3 {
			yellowLog("Update requires 3 total arguments.")
			yellowLog("update <secretName> <newValue>")
			return
		}

		res, ok := encryption.UpdateSecret(args[1], args[2])
		if !ok {
			redLog(res)
			return
		}

		greenLog("Secret has been updated.")

	case "list", "l":

		if len(args) != 1 {
			yellowLog("List requires 1 argument.")
			yellowLog("list")
			return
		}
		displayPublicVault()
	}
}

func displayPublicVault() {
	publicVault := encryption.GetPublicVault()

	header := fmt.Sprintf("%-25s | %-20s | %-20s | %-10s\n", "Key", "Date Added", "Last Modified", "Version")
	divider := fmt.Sprintln(
		strings.Repeat("-", 25) + "-+-" +
			strings.Repeat("-", 20) + "-+-" +
			strings.Repeat("-", 20) + "-+-" +
			strings.Repeat("-", 10),
	)

	greenLog(header)
	greenLog(divider)

	for _, entry := range publicVault {
		row := fmt.Sprintf(
			"%-25s | %-20s | %-20s | %-10s\n",
			entry.Key,
			entry.DateAdded,
			entry.LastModified,
			entry.Version,
		)
		greenLog(row)
	}
}

func greenLog(s string) {
	fmt.Println("\033[32mCove CLI> " + s + "\033[0m")
}

func yellowLog(s string) {
	fmt.Println("\033[33mCove CLI> " + s + "\033[0m")
}

func redLog(s string) {
	fmt.Println("\033[31mCove CLI> " + s + "\033[0m")
}
