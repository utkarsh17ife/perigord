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

package network

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"

	"github.com/utkarsh17ife/perigord/project"

	"golang.org/x/crypto/ssh/terminal"
)

type NetworkConfig struct {
	url           string
	keystore_path string
}

var networks map[string]NetworkConfig

func InitNetworks() error {
	prj, err := project.FindProject()
	if err != nil {
		return errors.New("Could not find project")
	}

	viper.SetConfigFile(filepath.Join(prj.AbsPath(), project.ProjectConfigFilename))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	configs := viper.GetStringMap("networks")
	networks = make(map[string]NetworkConfig)

	for key, _ := range configs {
		config := viper.GetStringMapString("networks." + key)
		networks[key] = NetworkConfig{
			url:           config["url"],
			keystore_path: config["keystore"],
		}
	}

	return nil
}

type Network struct {
	name       string
	rpc_client *rpc.Client
	client     *ethclient.Client
	keystore   *keystore.KeyStore
}

func Dial(name string) (*Network, error) {
	if config, ok := networks[name]; ok {
		rpc_client, err := rpc.Dial(config.url)
		if err != nil {
			return nil, err
		}

		client := ethclient.NewClient(rpc_client)

		ks := keystore.NewKeyStore(config.keystore_path, keystore.StandardScryptN, keystore.StandardScryptP)
		if len(ks.Accounts()) == 0 {
			return nil, errors.New("No accounts configured for this network, did you set the keystore path in perigord.yaml?")
		}

		ret := &Network{
			name:       name,
			rpc_client: rpc_client,
			client:     client,
			keystore:   ks,
		}

		return ret, nil
	}

	return nil, errors.New("No such network " + name)
}

func (n *Network) Name() string {
	return n.name
}

func (n *Network) Url() string {
	return networks[n.name].url
}

func (n *Network) KeystorePath() string {
	return networks[n.name].keystore_path
}

func (n *Network) RpcClient() *rpc.Client {
	return n.rpc_client
}

func (n *Network) Client() *ethclient.Client {
	return n.client
}

func (n *Network) Keystore() *keystore.KeyStore {
	return n.keystore
}

func (n *Network) Accounts() []accounts.Account {
	return n.keystore.Accounts()
}

func (n *Network) Unlock(a accounts.Account, passphrase string) error {
	return n.keystore.Unlock(a, passphrase)
}

func (n *Network) UnlockWithPrompt(a accounts.Account) error {
	// https://stackoverflow.com/questions/2137357/getpasswd-functionality-in-go
	fmt.Print("Enter passphrase for account ", a.Address.Hex(), ": ")
	passphrase, err := terminal.ReadPassword(0)
	if err != nil {
		return err
	}
	fmt.Println()

	return n.Unlock(a, string(passphrase))
}

func (n *Network) NewTransactor(a accounts.Account) *bind.TransactOpts {
	return &bind.TransactOpts{
		From: a.Address,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != a.Address {
				return nil, errors.New("not authorized to sign this account")
			}

			signature, err := n.keystore.SignHash(a, signer.Hash(tx).Bytes())
			if err != nil {
				return nil, err
			}

			return tx.WithSignature(signer, signature)
		},
	}
}
