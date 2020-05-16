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

type Client struct {
	bcclientgo.Client
	Node   *bcgo.Node
	App    fyne.App
	Window fyne.Window
	Dialog dialog.Dialog
}

func (c *Client) ExistingNode(alias string, password []byte, callback func(*bcgo.Node)) {
	// Get key store
	keystore, err := bcgo.GetKeyDirectory(c.Root)
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
	// Create node
	node := &bcgo.Node{
		Alias:    alias,
		Key:      key,
		Cache:    c.Cache,
		Network:  c.Network,
		Channels: make(map[string]*bcgo.Channel),
	}

	callback(node)
}

func (c *Client) GetNode() *bcgo.Node {
	if c.Node == nil {
		nc := make(chan *bcgo.Node, 1)
		go c.ShowAccessDialog(func(n *bcgo.Node) {
			nc <- n
		})
		c.Node = <-nc
	}
	return c.Node
}

func (c *Client) GetLogo() fyne.CanvasObject {
	return &canvas.Image{
		Resource: data.NewThemedResource(data.Logo),
		//FillMode: canvas.ImageFillContain,
		FillMode: canvas.ImageFillOriginal,
	}
}

func (c *Client) NewNode(alias string, password []byte, callback func(*bcgo.Node)) {
	// Create Progress Dialog
	progress := dialog.NewProgress("Registering", "message", c.Window)
	defer progress.Hide()
	listener := &ui.ProgressMiningListener{Func: progress.SetValue}

	// Get key store
	keystore, err := bcgo.GetKeyDirectory(c.Root)
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
	// Create node
	node := &bcgo.Node{
		Alias:    alias,
		Key:      key,
		Cache:    c.Cache,
		Network:  c.Network,
		Channels: make(map[string]*bcgo.Channel),
	}

	// Register Alias
	if err := aliasgo.Register(node, listener); err != nil {
		c.ShowError(err)
		return
	}

	callback(node)
}

func (c *Client) ShowAccessDialog(callback func(*bcgo.Node)) {
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

func (c *Client) ShowAccount() {
	//
}

func (c *Client) ShowError(err error) {
	if c.Dialog != nil {
		c.Dialog.Hide()
	}
	c.Dialog = dialog.NewError(err, c.Window)
	c.Dialog.Show()
}

func (c *Client) ShowNode() {
	log.Println("ShowNode")
	node := c.GetNode()
	log.Println("Alias:", node.Alias)
}
