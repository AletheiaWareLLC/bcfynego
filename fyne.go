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
	"aletheiaware.com/aliasgo"
	"aletheiaware.com/bcclientgo"
	"aletheiaware.com/bcfynego/ui"
	"aletheiaware.com/bcfynego/ui/account"
	"aletheiaware.com/bcfynego/ui/data"
	"aletheiaware.com/bcgo"
	"aletheiaware.com/cryptogo"
	"bytes"
	"errors"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"log"
	"os"
	"runtime/debug"
)

type BCFyne struct {
	App            fyne.App
	Window         fyne.Window
	Dialog         dialog.Dialog
	OnKeysExported func(string)
	OnKeysImported func(string)
	OnSignedIn     func(*bcgo.Node)
	OnSignedUp     func(*bcgo.Node)
	OnSignedOut    func()
}

func NewBCFyne(a fyne.App, w fyne.Window) *BCFyne {
	return &BCFyne{
		App:    a,
		Window: w,
	}
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

	if c := callback; c != nil {
		c(node)
	}
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
	// Show progress dialog
	progress := dialog.NewProgressInfinite("Creating", "Creating "+alias, f.Window)
	progress.Show()
	defer progress.Hide()

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

	{
		// Show Progress Dialog
		progress := dialog.NewProgress("Registering", "Registering "+alias, f.Window)
		progress.Show()
		defer progress.Hide()
		listener := &ui.ProgressMiningListener{Func: progress.SetValue}

		// Register Alias
		if err := aliasgo.Register(node, listener); err != nil {
			f.ShowError(err)
			return
		}
	}

	if c := callback; c != nil {
		c(node)
	}
}

func (f *BCFyne) ShowAccessDialog(client *bcclientgo.BCClient, callback func(*bcgo.Node)) {
	signIn := account.NewSignIn()
	importKey := account.NewImportKey()
	signUp := account.NewSignUp()
	if d := f.Dialog; d != nil {
		d.Hide()
	}
	tos := &widget.Hyperlink{Text: "Terms of Service"}
	tos.SetURLFromString("https://aletheiaware.com/terms-of-service.html")
	pp := &widget.Hyperlink{Text: "Privacy Policy", Alignment: fyne.TextAlignTrailing}
	pp.SetURLFromString("https://aletheiaware.com/privacy-policy.html")
	f.Dialog = dialog.NewCustom("Account Access", "Cancel",
		container.NewVBox(
			widget.NewAccordionContainer(
				&widget.AccordionItem{Title: "Sign In", Detail: signIn.CanvasObject(), Open: true},
				widget.NewAccordionItem("Import Keys", importKey.CanvasObject()),
				widget.NewAccordionItem("Sign Up", signUp.CanvasObject()),
			),
			container.NewGridWithColumns(2, tos, pp),
		),
		f.Window)

	signIn.SignInButton.OnTapped = func() {
		if d := f.Dialog; d != nil {
			d.Hide()
		}
		alias := signIn.Alias.Text
		password := []byte(signIn.Password.Text)
		if len(password) < cryptogo.MIN_PASSWORD {
			f.ShowError(fmt.Errorf(cryptogo.ERROR_PASSWORD_TOO_SHORT, len(password), cryptogo.MIN_PASSWORD))
			return
		}
		f.ExistingNode(client, alias, password, func(node *bcgo.Node) {
			if c := callback; c != nil {
				c(node)
			}
			if c := f.OnSignedIn; c != nil {
				go c(node)
			}
		})
	}
	importKey.ImportKeyButton.OnTapped = func() {
		if d := f.Dialog; d != nil {
			d.Hide()
		}

		host := bcgo.GetBCWebsite()
		alias := importKey.Alias.Text
		access := importKey.Access.Text

		// Show Progress Dialog
		progress := dialog.NewProgress("Importing", fmt.Sprintf("Importing %s from %s", alias, host), f.Window)
		progress.Show()
		defer progress.Hide()

		if err := client.ImportKeys(host, alias, access); err != nil {
			f.ShowError(err)
			return
		}
		if c := f.OnKeysImported; c != nil {
			go c(alias)
		}
	}
	signUp.SignUpButton.OnTapped = func() {
		if d := f.Dialog; d != nil {
			d.Hide()
		}
		alias := signUp.Alias.Text
		password := []byte(signUp.Password.Text)
		confirm := []byte(signUp.Confirm.Text)

		err := aliasgo.ValidateAlias(alias)
		if err != nil {
			f.ShowError(err)
			return
		}

		if len(password) < cryptogo.MIN_PASSWORD {
			f.ShowError(fmt.Errorf(cryptogo.ERROR_PASSWORD_TOO_SHORT, len(password), cryptogo.MIN_PASSWORD))
			return
		}
		if !bytes.Equal(password, confirm) {
			f.ShowError(errors.New(cryptogo.ERROR_PASSWORDS_DO_NOT_MATCH))
			return
		}
		f.NewNode(client, alias, password, func(node *bcgo.Node) {
			if c := callback; c != nil {
				c(node)
			}
			if c := f.OnSignedUp; c != nil {
				go c(node)
			}
		})
	}

	rootDir, err := client.GetRoot()
	if err != nil {
		log.Println(err)
	} else {
		keystore, err := bcgo.GetKeyDirectory(rootDir)
		if err != nil {
			log.Println(err)
		} else {
			keys, err := cryptogo.ListRSAPrivateKeys(keystore)
			if err != nil {
				log.Println(err)
			} else if len(keys) > 0 {
				signIn.Alias.SetOptions(keys)
				signIn.Alias.SetText(keys[0])
				importKey.Alias.SetText(keys[0])
				signUp.Alias.SetText(keys[0])
			}
		}
	}

	if alias, ok := os.LookupEnv("ALIAS"); ok {
		signIn.Alias.SetText(alias)
		importKey.Alias.SetText(alias)
		signUp.Alias.SetText(alias)
	}

	if pwd, ok := os.LookupEnv("PASSWORD"); ok {
		signIn.Password.SetText(pwd)
	}

	f.Dialog.Show()
}

func (f *BCFyne) ShowAccount(client *bcclientgo.BCClient) {
	node, err := f.GetNode(client)
	if err != nil {
		f.ShowError(err)
		return
	}
	form, err := nodeView(node)
	if err != nil {
		f.ShowError(err)
		return
	}
	box := widget.NewVBox(
		form,
	)
	if d := f.Dialog; d != nil {
		d.Hide()
	}
	f.Dialog = dialog.NewCustom("Account", "OK", box, f.Window)
	box.Append(widget.NewButton("Export Keys", func() {
		f.ExportKeys(client, node)
	}))
	box.Append(widget.NewButton("Switch Keys", func() {
		f.Dialog.Hide()
		f.SwitchKeys(client)
		go f.ShowAccessDialog(client, nil)
	}))
	box.Append(widget.NewButton("Delete Keys", func() {
		f.Dialog.Hide()
		f.DeleteKeys(client, node)
	}))
	f.Dialog.Show()
}

func (f *BCFyne) DeleteKeys(client *bcclientgo.BCClient, node *bcgo.Node) {
	f.ShowError(fmt.Errorf("Not yet implemented: %s", "BCFyne.DeleteKeys"))
}

func (f *BCFyne) ExportKeys(client *bcclientgo.BCClient, node *bcgo.Node) {
	authentication := account.NewAuthentication(node.Alias)
	authentication.AuthenticateButton.OnTapped = func() {
		host := bcgo.GetBCWebsite()

		// Show Progress Dialog
		progress := dialog.NewProgress("Exporting", fmt.Sprintf("Exporting %s to %s", node.Alias, host), f.Window)
		progress.Show()

		password := []byte(authentication.Password.Text)
		if len(password) < cryptogo.MIN_PASSWORD {
			f.ShowError(fmt.Errorf(cryptogo.ERROR_PASSWORD_TOO_SHORT, len(password), cryptogo.MIN_PASSWORD))
			return
		}
		access, err := client.ExportKeys(host, node.Alias, password)
		if err != nil {
			f.ShowError(err)
			return
		}

		progress.Hide()

		form := widget.NewForm(
			widget.NewFormItem("Alias", widget.NewLabel(node.Alias)),
			widget.NewFormItem("Access Code", container.NewHBox(
				widget.NewLabel(access),
				widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
					f.Window.Clipboard().SetContent(access)
					dialog.ShowInformation("Copied", "Access code copied to clipboard", f.Window)
				}),
			)),
		)
		dialog.ShowCustom("Export Keys", "OK", form, f.Window)
		if c := f.OnKeysExported; c != nil {
			go c(node.Alias)
		}
	}
	dialog.ShowCustom("Account", "Cancel", authentication.CanvasObject(), f.Window)
}

func (f *BCFyne) SwitchKeys(client *bcclientgo.BCClient) {
	client.Root = ""
	client.Cache = nil
	client.Network = nil
	client.Node = nil
	if c := f.OnSignedOut; c != nil {
		go c()
	}
}

func (f *BCFyne) ShowError(err error) {
	log.Println("Error:", err)
	debug.PrintStack()
	if d := f.Dialog; d != nil {
		d.Hide()
	}
	f.Dialog = dialog.NewError(err, f.Window)
	f.Dialog.Show()
}

func (f *BCFyne) ShowNode(node *bcgo.Node) {
	form, err := nodeView(node)
	if err != nil {
		f.ShowError(err)
		return
	}
	if d := f.Dialog; d != nil {
		d.Hide()
	}
	f.Dialog = dialog.NewCustom("Node", "OK", form, f.Window)
	f.Dialog.Show()
}

func nodeView(node *bcgo.Node) (fyne.CanvasObject, error) {
	publicKeyBytes, err := cryptogo.RSAPublicKeyToPKIXBytes(&node.Key.PublicKey)
	if err != nil {
		return nil, err
	}

	aliasScroller := widget.NewHScrollContainer(ui.NewAliasView(func() string { return node.Alias }))
	publicKeyScroller := widget.NewHScrollContainer(ui.NewKeyView(func() []byte { return publicKeyBytes }))
	publicKeyScroller.SetMinSize(fyne.NewSize(10*theme.TextSize(), 0))

	return widget.NewForm(
		widget.NewFormItem(
			"Alias",
			aliasScroller,
		),
		widget.NewFormItem(
			"Public Key",
			publicKeyScroller,
		),
	), nil
}
