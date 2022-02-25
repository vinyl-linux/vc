/*
Copyright © 2022 James Condron <james@zero-internet.org.uk>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinyl-linux/vc"
)

var showLogs bool

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of either the specified script, or all scripts",
	Long:  "Show the status of either the specified script, or all scripts",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		db, err := vc.LoadDatabase(dbFile)
		if err != nil {
			return
		}

		if len(args) == 0 {
			for id, script := range db {
				printStatus(id, script)
			}

			return
		}

		config, err := vc.LoadConfig(args[0])
		if err != nil {
			return
		}

		s, ok := db[config.Sum]
		if !ok {
			fmt.Printf("The script described in %s is unknown, or has never been run", args[0])

			return nil
		}

		printStatus(config.Sum, s)

		return
	},
}

func printStatus(id string, s vc.Script) {
	fmt.Printf("Script %s (%s)\nStatus: %s\n",
		id, s.RunAt, s.RunningState)

	if s.Error != "" {
		fmt.Printf("Error: %s\n", s.Error)
	}

	if showLogs {
		fmt.Println("Logs:")
		for _, line := range s.Logs {
			fmt.Println(line)
		}
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().BoolVarP(&showLogs, "logs", "l", false, "Display logs for command runs")
}
