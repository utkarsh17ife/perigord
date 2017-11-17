// Copyright © 2017 PolySwarm <info@polyswarm.io>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package testing

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/polyswarm/perigord/contract"
	"github.com/polyswarm/perigord/migration"
	"github.com/polyswarm/perigord/project"
)

func SetUpTest() (*migration.Network, error) {
	prj, err := project.FindProject()
	if err != nil {
		return nil, errors.New("Could not find project")
	}

	viper.SetConfigFile(filepath.Join(prj.AbsPath(), project.ProjectConfigFilename))
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	migration.InitNetworks()
	// TODO: Fix this in config
	network, err := migration.Dial("dev")
	if err != nil {
		return nil, err
	}

	if err := migration.RunMigrations(context.Background(), network); err != nil {
		return nil, err
	}

	return network, nil
}

func TearDownTest() {
	contract.Reset()
}
