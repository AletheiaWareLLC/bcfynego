/*
 * Copyright 2020 Aletheia Ware LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bcfynego

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/aliasgo"
	"github.com/AletheiaWareLLC/bcclientgo"
	"github.com/AletheiaWareLLC/bcfynego/ui"
	"github.com/AletheiaWareLLC/bcfynego/ui/account"
	"github.com/AletheiaWareLLC/bcfynego/ui/data"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/cryptogo"
	"log"
	"os"
)

type BCFyneClient struct {
	bcclientgo.BCClient
	App    fyne.App
	Window fyne.Window
	Dialog dialog.Dialog
}

func (c *BCFyneClient) ExistingNode(alias string, password []byte, callback func(*bcgo.Node)) {
	rootDir, err := c.GetRoot()
	if err != nil {
		c.ShowError(err)
		return
	}
	// Get key store
	keystore, err := bcgo.GetKeyDirectory(rootDir)
	if err != nil {
		c.ShowError(err)
		return
	}
	// Get private key
	key, err := cryptogo.GetRSAPrivateKey(keystore, alias, password)
	if err != nil {
		c.ShowError(err)
		return
	}
	cache, err := c.GetCache()
	if err != nil {
		c.ShowError(err)
		return
	}
	network, err := c.GetNetwork()
	if err != nil {
		c.ShowError(err)
		return
	}
	// Create node
	node := &bcgo.Node{
		Alias:    alias,
		Key:      key,
		Cache:    cache,
		Network:  network,
		Channels: make(map[string]*bcgo.Channel),
	}

	callback(node)
}

func (c *BCFyneClient) GetNode() *bcgo.Node {
	if c.BCClient.Node == nil {
		nc := make(chan *bcgo.Node, 1)
		go c.ShowAccessDialog(func(n *bcgo.Node) {
			nc <- n
		})
		c.BCClient.Node = <-nc
	}
	return c.BCClient.Node
}

func (c *BCFyneClient) GetLogo() fyne.CanvasObject {
	return &canvas.Image{
		Resource: data.NewThemedResource(data.Logo),
		//FillMode: canvas.ImageFillContain,
		FillMode: canvas.ImageFillOriginal,
	}
}

func (c *BCFyneClient) NewNode(alias string, password []byte, callback func(*bcgo.Node)) {
	// Create Progress Dialog
	progress := dialog.NewProgress("Registering", "message", c.Window)
	defer progress.Hide()
	listener := &ui.ProgressMiningListener{Func: progress.SetValue}

	rootDir, err := c.GetRoot()
	if err != nil {
		c.ShowError(err)
		return
	}
	// Get key store
	keystore, err := bcgo.GetKeyDirectory(rootDir)
	if err != nil {
		c.ShowError(err)
		return
	}
	// Create private key
	key, err := cryptogo.CreateRSAPrivateKey(keystore, alias, password)
	if err != nil {
		c.ShowError(err)
		return
	}
	cache, err := c.GetCache()
	if err != nil {
		c.ShowError(err)
		return
	}
	network, err := c.GetNetwork()
	if err != nil {
		c.ShowError(err)
		return
	}
	// Create node
	node := &bcgo.Node{
		Alias:    alias,
		Key:      key,
		Cache:    cache,
		Network:  network,
		Channels: make(map[string]*bcgo.Channel),
	}

	// Register Alias
	if err := aliasgo.Register(node, listener); err != nil {
		c.ShowError(err)
		return
	}

	callback(node)
}

func (c *BCFyneClient) ShowAccessDialog(callback func(*bcgo.Node)) {
	signIn := account.NewSignIn()
	importKey := account.NewImportKey()
	signUp := account.NewSignUp()
	c.Dialog = dialog.NewCustom("Account Access", "Cancel",
		widget.NewAccordionContainer(
			&widget.AccordionItem{Title: "Sign In", Detail: signIn.CanvasObject(), Open: true},
			widget.NewAccordionItem("Import Key", importKey.CanvasObject()),
			widget.NewAccordionItem("Sign Up", signUp.CanvasObject()),
		), c.Window)

	if alias, err := bcgo.GetAlias(); err == nil {
		signIn.Alias.SetText(alias)
		importKey.Alias.SetText(alias)
	}
	if pwd, ok := os.LookupEnv("PASSWORD"); ok {
		signIn.Password.SetText(pwd)
		// TODO if alias was also set auto click
	}
	signIn.SignInButton.OnTapped = func() {
		c.Dialog.Hide()
		log.Println("Sign In Tapped")
		alias := signIn.Alias.Text
		password := []byte(signIn.Password.Text)
		if len(password) < cryptogo.MIN_PASSWORD {
			c.ShowError(errors.New(fmt.Sprintf(cryptogo.ERROR_PASSWORD_TOO_SHORT, len(password), cryptogo.MIN_PASSWORD)))
			return
		}
		c.ExistingNode(alias, password, callback)
	}
	importKey.ImportKeyButton.OnTapped = func() {
		c.Dialog.Hide()
		log.Println("Import Key Tapped")
		// TODO alias := importKey.Alias.Text
		// TODO access := importKey.Access.Text
	}
	signUp.SignUpButton.OnTapped = func() {
		c.Dialog.Hide()
		log.Println("Sign Up Tapped")
		alias := signUp.Alias.Text
		password := []byte(signUp.Password.Text)
		confirm := []byte(signUp.Confirm.Text)

		err := aliasgo.ValidateAlias(alias)
		if err != nil {
			c.ShowError(err)
			return
		}

		if len(password) < cryptogo.MIN_PASSWORD {
			c.ShowError(errors.New(fmt.Sprintf(cryptogo.ERROR_PASSWORD_TOO_SHORT, len(password), cryptogo.MIN_PASSWORD)))
			return
		}
		if !bytes.Equal(password, confirm) {
			c.ShowError(errors.New(cryptogo.ERROR_PASSWORDS_DO_NOT_MATCH))
			return
		}
		c.NewNode(alias, password, callback)
	}
	c.Dialog.Show()
}

func (c *BCFyneClient) ShowAccount() {
	//
}

func (c *BCFyneClient) ShowError(err error) {
	if c.Dialog != nil {
		c.Dialog.Hide()
	}
	c.Dialog = dialog.NewError(err, c.Window)
	c.Dialog.Show()
}

func (c *BCFyneClient) ShowNode() {
	log.Println("ShowNode")
	node := c.GetNode()
	info := fmt.Sprintf("Alias: %s\n", node.Alias)
	publicKeyBytes, err := cryptogo.RSAPublicKeyToPKIXBytes(&node.Key.PublicKey)
	if err == nil {
		info = fmt.Sprintf("%sPublicKey: %s\n", info, base64.RawURLEncoding.EncodeToString(publicKeyBytes))
	}
	c.Dialog = dialog.NewInformation("Node", info, c.Window)
	c.Dialog.Show()
}
