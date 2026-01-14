// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hello is a hello, world program, demonstrating
// how to write a simple command-line program.
//
// Usage:
//
//	hello [options] [name]
//
// The options are:
//
//	-g greeting
//		Greet with the given greeting, instead of "Hello".
//
//	-r
//		Greet in reverse.
//
// By default, hello greets the world.
// If a name is specified, hello greets that name instead.
package main

import (
	"log"

	"github.com/andersonreyes/jarvis/money"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jarvis [command] [options]",
	Short: "tools on tools",
	Long:  "Jarvis command line to manage all things for anderson",
}

func init() {
	rootCmd.AddCommand(money.MoneyCmd)
	rootCmd.AddCommand(money.FireflyCmd)
}

func main() {
	// Configure logging for a command-line program.
	log.SetFlags(0)
	log.SetPrefix("jarvis: ")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
