// Copyright (C) 2017 go-nebulas authors
//
// This file is part of the go-nebulas library.
//
// the go-nebulas library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-nebulas library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.
//

package main

import (
	"fmt"
	"io/ioutil"

	"github.com/nebulasio/go-nebulas/cmd/console"
	"github.com/nebulasio/go-nebulas/core"
	"github.com/urfave/cli"
)

var (
	accountCommand = cli.Command{
		Name:     "account",
		Usage:    "Manage accounts",
		Category: "ACCOUNT COMMANDS",
		Description: `
Manage accounts, list all existing accounts, import a private key into a new
account, create a new account or update an existing account.`,

		Subcommands: []cli.Command{
			{
				Name:   "new",
				Usage:  "Create a new account",
				Action: accountCreate,
				Description: `
    geth account new

Creates a new account and prints the address.

The account is saved in encrypted format, you are prompted for a passphrase.

You must remember this passphrase to unlock your account in the future.`,
			},
			{
				Name:   "list",
				Usage:  "Print summary of existing addresses",
				Action: accountList,
				Description: `
Print a short summary of all accounts`,
			},
			{
				Name:      "update",
				Usage:     "Update an existing account",
				Action:    accountUpdate,
				ArgsUsage: "<address>",
				Description: `
    geth account update <address>

Update an existing account.

The account is saved in the newest version in encrypted format, you are prompted
for a passphrase to unlock the account and another to save the updated file.`,
			},
			{
				Name:      "import",
				Usage:     "Import a private key into a new account",
				Action:    accountImport,
				ArgsUsage: "<keyFile>",
				Description: `
    geth account import <keyfile>

Imports an encrypted private key from <keyfile> and creates a new account.`,
			},
		},
	}
)

// accountList list account
func accountList(ctx *cli.Context) error {
	neb := makeNeb(ctx)
	for index, addr := range neb.AccountManager().Accounts() {
		fmt.Printf("Account #%d: %s\n", index, addr.ToHex())
		index++
	}
	return nil
}

// accountCreate creates a new account into the keystore
func accountCreate(ctx *cli.Context) error {
	neb := makeNeb(ctx)
	passphrase := getPassPhrase("Your new account is locked with a passphrase. Please give a passphrase. Do not forget this passphrase.", true)

	addr, err := neb.AccountManager().NewAccount([]byte(passphrase))
	fmt.Printf("Address: %s\n", addr.ToHex())
	return err
}

// accountUpdate update
func accountUpdate(ctx *cli.Context) error {
	if len(ctx.Args()) == 0 {
		FatalF("No accounts specified to update")
	}
	neb := makeNeb(ctx)

	for _, address := range ctx.Args() {
		addr, err := core.AddressParse(address)
		if err != nil {
			FatalF("address parse failed:%s,%s", address, err)
		}
		oldPassphrase := getPassPhrase("Please input current passhprase", false)
		newPassword := getPassPhrase("Please give a new password. Do not forget this password.", true)

		err = neb.AccountManager().Update(addr, []byte(oldPassphrase), []byte(newPassword))
		if err != nil {
			FatalF("account update failed:%s,%s", address, err)
		}
		fmt.Printf("Updated address: %s\n", addr.ToHex())
	}
	return nil
}

// accountImport import keyfile
func accountImport(ctx *cli.Context) error {
	keyfile := ctx.Args().First()
	if len(keyfile) == 0 {
		FatalF("keyfile must be given as argument")
	}
	keyJSON, err := ioutil.ReadFile(keyfile)
	if err != nil {
		FatalF("file read failed:%s", err)
	}

	neb := makeNeb(ctx)
	passphrase := getPassPhrase("", false)

	addr, err := neb.AccountManager().Import([]byte(keyJSON), []byte(passphrase))
	if err != nil {
		FatalF("key import failed:%s", err)
	}
	fmt.Printf("Import address: %s\n", addr.ToHex())
	return nil
}

// getPassPhrase get passphrase from consle
func getPassPhrase(prompt string, confirmation bool) string {
	if prompt != "" {
		fmt.Println(prompt)
	}
	passphrase, err := console.Stdin.PromptPassphrase("Passphrase: ")
	if err != nil {
		FatalF("Failed to read passphrase: %v", err)
	}
	if confirmation {
		confirm, err := console.Stdin.PromptPassphrase("Repeat passphrase: ")
		if err != nil {
			FatalF("Failed to read passphrase confirmation: %v", err)
		}
		if passphrase != confirm {
			FatalF("Passphrases do not match")
		}
	}
	return passphrase
}
