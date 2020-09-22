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
	"fyne.io/fyne/theme"
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
	"runtime/debug"
)

type BCFyne struct {
	App    fyne.App
	Window fyne.Window
	Dialog dialog.Dialog
}

func (f *BCFyne) ExistingNode(client *bcclientgo.BCClient, alias string, password []byte, callback func(*bcgo.Node)) {
	rootDir, err := client.GetRoot()
	if err != nil {
		f.ShowError(err)
		return
	}
	// Get key store
	keystore, err := bcgo.GetKeyDirectory(rootDir)
	if err != nil {
		f.ShowError(err)
		return
	}
	// Get private key
	key, err := cryptogo.GetRSAPrivateKey(keystore, alias, password)
	if err != nil {
		f.ShowError(err)
		return
	}
	cache, err := client.GetCache()
	if err != nil {
		f.ShowError(err)
		return
	}
	network, err := client.GetNetwork()
	if err != nil {
		f.ShowError(err)
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

func (f *BCFyne) GetNode(client *bcclientgo.BCClient) (*bcgo.Node, error) {
	if client.Node == nil {
		nc := make(chan *bcgo.Node, 1)
		go f.ShowAccessDialog(client, func(n *bcgo.Node) {
			nc <- n
		})
		client.Node = <-nc
	}
	return client.Node, nil
}

func (f *BCFyne) GetLogo() fyne.CanvasObject {
	return &canvas.Image{
		Resource: data.NewPrimaryThemedResource(data.Logo),
		//FillMode: canvas.ImageFillContain,
		FillMode: canvas.ImageFillOriginal,
	}
}

func (f *BCFyne) NewNode(client *bcclientgo.BCClient, alias string, password []byte, callback func(*bcgo.Node)) {
	rootDir, err := client.GetRoot()
	if err != nil {
		f.ShowError(err)
		return
	}
	// Get key store
	keystore, err := bcgo.GetKeyDirectory(rootDir)
	if err != nil {
		f.ShowError(err)
		return
	}
	// Create private key
	key, err := cryptogo.CreateRSAPrivateKey(keystore, alias, password)
	if err != nil {
		f.ShowError(err)
		return
	}
	cache, err := client.GetCache()
	if err != nil {
		f.ShowError(err)
		return
	}
	network, err := client.GetNetwork()
	if err != nil {
		f.ShowError(err)
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

	// Create Progress Dialog
	progress := dialog.NewProgress("Registering", "Registering "+alias, f.Window)
	progress.Show()
	defer progress.Hide()
	listener := &ui.ProgressMiningListener{Func: progress.SetValue}

	// Register Alias
	if err := aliasgo.Register(node, listener); err != nil {
		f.ShowError(err)
		return
	}

	callback(node)
}

func (f *BCFyne) ShowAccessDialog(client *bcclientgo.BCClient, callback func(*bcgo.Node)) {
	signIn := account.NewSignIn()
	importKey := account.NewImportKey()
	signUp := account.NewSignUp()
	if f.Dialog != nil {
		f.Dialog.Hide()
	}
	f.Dialog = dialog.NewCustom("Account Access", "Cancel",
		widget.NewAccordionContainer(
			&widget.AccordionItem{Title: "Sign In", Detail: signIn.CanvasObject(), Open: true},
			widget.NewAccordionItem("Import Keys", importKey.CanvasObject()),
			widget.NewAccordionItem("Sign Up", signUp.CanvasObject()),
		), f.Window)

	if alias, err := bcgo.GetAlias(); err == nil {
		signIn.Alias.SetText(alias)
		importKey.Alias.SetText(alias)
		signUp.Alias.SetText(alias)
	}
	if pwd, ok := os.LookupEnv("PASSWORD"); ok {
		signIn.Password.SetText(pwd)
		// TODO if alias was also set auto click
	}
	signIn.SignInButton.OnTapped = func() {
		f.Dialog.Hide()
		log.Println("Sign In Tapped")
		alias := signIn.Alias.Text
		password := []byte(signIn.Password.Text)
		if len(password) < cryptogo.MIN_PASSWORD {
			f.ShowError(errors.New(fmt.Sprintf(cryptogo.ERROR_PASSWORD_TOO_SHORT, len(password), cryptogo.MIN_PASSWORD)))
			return
		}
		f.ExistingNode(client, alias, password, callback)
	}
	importKey.ImportKeyButton.OnTapped = func() {
		f.Dialog.Hide()
		log.Println("Import Key Tapped")
		host := bcgo.GetBCWebsite()
		alias := importKey.Alias.Text
		access := importKey.Access.Text
		if err := client.ImportKeys(host, alias, access); err != nil {
			f.ShowError(err)
		}
	}
	signUp.SignUpButton.OnTapped = func() {
		f.Dialog.Hide()
		log.Println("Sign Up Tapped")
		alias := signUp.Alias.Text
		password := []byte(signUp.Password.Text)
		confirm := []byte(signUp.Confirm.Text)

		err := aliasgo.ValidateAlias(alias)
		if err != nil {
			f.ShowError(err)
			return
		}

		if len(password) < cryptogo.MIN_PASSWORD {
			f.ShowError(errors.New(fmt.Sprintf(cryptogo.ERROR_PASSWORD_TOO_SHORT, len(password), cryptogo.MIN_PASSWORD)))
			return
		}
		if !bytes.Equal(password, confirm) {
			f.ShowError(errors.New(cryptogo.ERROR_PASSWORDS_DO_NOT_MATCH))
			return
		}
		f.NewNode(client, alias, password, callback)
	}
	f.Dialog.Show()
}

func (f *BCFyne) ShowAccount(client *bcclientgo.BCClient) {
	log.Println("ShowAccount")
	node, err := f.GetNode(client)
	if err != nil {
		f.ShowError(err)
		return
	}
	box := widget.NewVBox(
		f.nodeUI(node),
		widget.NewButton("Export Keys", func() {
			access, err := client.ExportKeys(bcgo.GetBCWebsite(), node.Alias)
			if err != nil {
				f.ShowError(err)
				return
			}
			info := fmt.Sprintf("Alias: %s\nAccess Code: %s", node.Alias, access)
			dialog := dialog.NewInformation("Export Keys", info, f.Window)
			dialog.Show()
		}),
		widget.NewButton("Switch Keys", func() {
			// TODO
			f.ShowError(fmt.Errorf("Not yet implemented: %s", "BCFyne.Account.SwitchKeys"))
		}),
		widget.NewButton("Delete Keys", func() {
			// TODO
			f.ShowError(fmt.Errorf("Not yet implemented: %s", "BCFyne.Account.DeleteKeys"))
		}),
	)
	if f.Dialog != nil {
		f.Dialog.Hide()
	}
	f.Dialog = dialog.NewCustom("Account", "OK", box, f.Window)
	f.Dialog.Show()
}

func (f *BCFyne) ShowError(err error) {
	log.Println("Error:", err)
	debug.PrintStack()
	if f.Dialog != nil {
		f.Dialog.Hide()
	}
	f.Dialog = dialog.NewError(err, f.Window)
	f.Dialog.Show()
}

func (f *BCFyne) ShowNode(node *bcgo.Node) {
	log.Println("ShowNode")
	form := f.nodeUI(node)
	if f.Dialog != nil {
		f.Dialog.Hide()
	}
	f.Dialog = dialog.NewCustom("Node", "OK", form, f.Window)
	f.Dialog.Show()
}

func (f *BCFyne) nodeUI(node *bcgo.Node) fyne.CanvasObject {
	publicKeyBytes, err := cryptogo.RSAPublicKeyToPKIXBytes(&node.Key.PublicKey)
	if err != nil {
		f.ShowError(err)
		return nil
	}

	publicKeyBase64 := base64.RawURLEncoding.EncodeToString(publicKeyBytes)
	var publicKeyRunes []rune
	for i, r := range []rune(publicKeyBase64) {
		if i > 0 && i%64 == 0 {
			publicKeyRunes = append(publicKeyRunes, '\n')
		}
		publicKeyRunes = append(publicKeyRunes, r)
	}

	aliasScroller := widget.NewHScrollContainer(widget.NewLabel(node.Alias))
	publicKeyScroller := widget.NewHScrollContainer(widget.NewLabelWithStyle(string(publicKeyRunes), fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}))
	publicKeyScroller.SetMinSize(fyne.NewSize(10*theme.TextSize(), 0))

	form := widget.NewForm(
		widget.NewFormItem(
			"Alias",
			aliasScroller,
		),
		widget.NewFormItem(
			"Public Key",
			publicKeyScroller,
		),
	)
	return form
}
