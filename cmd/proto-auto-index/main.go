// Copyright 2021 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"go.linka.cloud/protoautoindex"
)

var (
	version     = "dev"
	showVersion bool
	cmd         = cobra.Command{
		Use:  "proto-auto-index [file]",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Println(version)
				return nil
			}
			if len(args) == 0 {
				return errors.New("no file specified")
			}
			i := protoautoindex.New()
			if err := i.Parse(args[0]); err != nil {
				return err
			}
			if err := i.SetIndexes(); err != nil {
				return err
			}
			if err := i.Write(args[0]); err != nil {
				return err
			}
			return nil
		},
	}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "display build version")
	cmd.ExecuteContext(ctx)
}
