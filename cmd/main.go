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

package main

import (
	"aletheiaware.com/bcclientgo"
	"aletheiaware.com/bcfynego"
	"aletheiaware.com/bcfynego/ui"
	"aletheiaware.com/bcfynego/ui/data"
	"aletheiaware.com/bcgo"
	"encoding/base64"
	"flag"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"strings"
)

var peer = flag.String("peer", "", "BC peer")

func main() {
	// Parse command line flags
	flag.Parse()

	// Set log flags
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create Application
	a := app.NewWithID("com.aletheiaware.bc")

	// Create Window
	w := a.NewWindow("BC")

	// Create BC Client
	c := bcclientgo.NewBCClient(bcgo.SplitRemoveEmpty(*peer, ",")...)

	// Create BC Fyne
	f := bcfynego.NewBCFyne(a, w)

	logo := f.GetLogo()

	block := ui.NewBlockView()
	setAddressAction := func(address string) {
		parts := strings.Split(address, "/")
		channel := parts[0]
		var blockhash []byte
		if len(parts) > 1 {
			bh, err := base64.RawURLEncoding.DecodeString(parts[1])
			if err != nil {
				f.ShowError(err)
				return
			}
			blockhash = bh
		}
		cache, err := c.GetCache()
		if err != nil {
			f.ShowError(err)
			return
		}
		network, err := c.GetNetwork()
		if err != nil {
			f.ShowError(err)
			return
		}
		if blockhash == nil {
			ref, err := bcgo.GetHeadReference(channel, cache, network)
			if err != nil {
				f.ShowError(err)
				return
			}
			blockhash = ref.BlockHash
		}
		block.SetHash(blockhash)
		b, err := bcgo.GetBlock(channel, cache, network, blockhash)
		if err != nil {
			f.ShowError(err)
			return
		}
		block.SetBlock(b)
	}

	address := widget.NewEntry()
	address.SetPlaceHolder("Channel")
	address.OnSubmitted = func(address string) {
		go setAddressAction(address)
	}

	w.SetContent(container.NewBorder(logo, nil, nil, nil, container.NewBorder(
		container.NewBorder(nil, nil, nil, container.NewHBox(
			widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
				go setAddressAction(address.Text)
			}),
			widget.NewButtonWithIcon("", theme.NewThemedResource(data.AccountIcon), func() {
				go f.ShowAccount(c)
			}),
			widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
				go settings(f, c)
			}),
		), address),
		nil,
		nil,
		nil,
		container.NewScroll(block),
	)))
	w.Resize(ui.WindowSize)
	w.CenterOnScreen()
	w.ShowAndRun()
}

func settings(f *bcfynego.BCFyne, c *bcclientgo.BCClient) {
	form := widget.NewForm()

	root := ui.NewRootView(func() string {
		r, err := c.GetRoot()
		if err != nil {
			f.ShowError(err)
			return ""
		}
		return r
	})
	form.Append("Root", container.NewBorder(nil, nil, nil, widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		go func() {
			c.SetRoot(root.Text)
			root.Refresh()
		}()
	}), root))

	peerList := &widget.List{
		Length: func() int {
			return len(c.Peers)
		},
		CreateItem: func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), nil), widget.NewLabel(""))
		},
	}
	peerList.UpdateItem = func(index widget.ListItemID, item fyne.CanvasObject) {
		if index < 0 || index >= len(c.Peers) {
			return
		}
		p := c.Peers[index]
		os := item.(*fyne.Container).Objects
		os[0].(*widget.Label).SetText(p)
		os[1].(*widget.Button).OnTapped = func() {
			c.SetPeers(append(c.Peers[:index], c.Peers[index+1:]...)...)
			form.Refresh()
		}
	}
	form.Append("Peers", container.NewVBox(
		peerList,
		container.NewGridWithColumns(2,
			widget.NewButton("Add", func() {
				entry := widget.NewEntry()
				entry.SetPlaceHolder("Domain")
				dialog.ShowForm("Add Peer", "Add", "Cancel", []*widget.FormItem{
					widget.NewFormItem("Peer", entry),
				}, func(ok bool) {
					if ok {
						c.SetPeers(append(c.Peers, entry.Text)...)
						form.Refresh()
					}
				}, f.Window)
			}),
			widget.NewButton("Reset", func() {
				dialog.ShowConfirm("Reset Peers", "Reset peers to default?", func(reset bool) {
					if !reset {
						return
					}
					c.SetPeers(bcgo.GetBCHost())
					form.Refresh()
				}, f.Window)
			}),
		),
	))

	form.Append("Cache", container.NewVBox(
		ui.NewCacheView(func() bcgo.Cache {
			h, err := c.GetCache()
			if err != nil {
				f.ShowError(err)
				return nil
			}
			return h
		}),
		widget.NewButton("Purge", func() {
			dialog.ShowConfirm("Purge Cache", "Remove all data from cache?", func(reset bool) {
				if !reset {
					return
				}
				go func() {
					if err := c.Purge(); err != nil {
						f.ShowError(err)
						return
					}
				}()
				form.Refresh()
			}, f.Window)
		}),
	))

	form.Append("Network", container.NewVBox(
		ui.NewNetworkView(func() bcgo.Network {
			n, err := c.GetNetwork()
			if err != nil {
				f.ShowError(err)
				return nil
			}
			return n
		}),
	))
	dialog.ShowCustom("Settings", "OK", form, f.Window)
}
