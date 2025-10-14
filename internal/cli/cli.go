package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/LSariol/Cove/internal/database"
	"github.com/LSariol/Cove/internal/server"
)

type CLI struct {
	DB *database.Database
}

func NewCLI(db *database.Database) *CLI {
	return &CLI{
		DB: db,
	}
}

func (c *CLI) StartCLI(ctx context.Context) {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Cove CLI> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		c.parseCLI(ctx, strings.Fields(input))
	}
}

func (c *CLI) parseCLI(ctx context.Context, args []string) {

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

		res, err := c.DB.GetSecret(ctx, args[1])
		if err != nil {
			redLog(fmt.Sprintf("error getting %q: %w", args[1], err))
			return
		}

		greenLog(fmt.Sprintf("%s : %s", res.Key, res.Value))

	case "add", "a":

		if len(args) != 3 {
			yellowLog("add requires 3 total arguments.")
			yellowLog("add <secretName> <value>")
			return
		}

		var newSecret database.Secret = database.Secret{
			Key:   args[1],
			Value: args[2],
		}

		secret, err := c.DB.CreateSecret(ctx, newSecret)
		if err != nil {
			redLog(err.Error())
		}

		greenLog(fmt.Sprintf("%s has been created at %q", secret.Key, secret.DateAdded))

	case "remove", "r", "delete", "d":

		if len(args) != 2 {
			yellowLog("Get requires 2 total arguments.")
			yellowLog("remove <secret>")
			return
		}

		err := c.DB.DeleteSecret(ctx, args[1])
		if err != nil {
			redLog(err.Error())
		}

		greenLog("Secret has been removed")

	case "update", "u":

		if len(args) != 3 {
			yellowLog("Update requires 3 total arguments.")
			yellowLog("update <secretName> <newValue>")
			return
		}

		var newSecret database.Secret = database.Secret{
			Key:   args[1],
			Value: args[2],
		}

		err := c.DB.UpdateSecret(ctx, newSecret)
		if err != nil {
			redLog(err.Error())
		}

		greenLog("Secret has been updated.")

	case "list", "l":

		if len(args) != 1 {
			yellowLog("List requires 1 argument.")
			yellowLog("list")
			return
		}
		c.displayPublicVault(ctx)

	case "bootstrap", "b":

		if len(args) != 2 {
			yellowLog("bootstrap requires 2 arguments.")
			yellowLog("bootstrap <clear/lock>")
			return
		}

		if args[1] == strings.ToLower("clear") {
			if err := server.DeleteBootstrapMarker(); err != nil {
				redLog(fmt.Sprintf("Error clearing marker: %w", err))
			}
		}

		if args[1] == strings.ToLower("lock") {
			if err := server.CreateBootstrapMarker(); err != nil {
				redLog(fmt.Sprintf("Error clearing marker: %w", err))
			}
		}
	}
}

func (c *CLI) displayPublicVault(ctx context.Context) {
	publicVault, err := c.DB.GetAllKeys(ctx)
	if err != nil {
		redLog(fmt.Sprintf("GetAllKeys: %w", err))
	}

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
