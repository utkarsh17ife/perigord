// Copyright © 2017 PolySwarm <info@polyswarm.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	perigord "github.com/utkarsh17ife/perigord/perigord/cmd"
	"github.com/utkarsh17ife/perigord/project"
)

var RootCmd = &cobra.Command{
	Use:   "stub",
	Short: "Linked into perigord projects to dispatch commands from the main application",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		perigord.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	prj, err := project.FindProject()
	if prj != nil && err == nil {
		viper.SetConfigFile(filepath.Join(prj.AbsPath(), project.ProjectConfigFilename))
		if err := viper.ReadInConfig(); err != nil {
			perigord.Fatal(err)
		}
	}
}
