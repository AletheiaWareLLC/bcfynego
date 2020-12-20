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
	"encoding/base64"
	"flag"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcclientgo"
	"github.com/AletheiaWareLLC/bcfynego"
	"github.com/AletheiaWareLLC/bcfynego/ui"
	"github.com/AletheiaWareLLC/bcfynego/ui/data"
	"github.com/AletheiaWareLLC/bcgo"
	"log"
	"os"
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

	// Set environment variable
	if a.Settings().BuildType() == fyne.BuildRelease {
		os.Setenv("LIVE", "true")
	}

	// Create BC Client
	c := bcclientgo.NewBCClient(bcgo.SplitRemoveEmpty(*peer, ",")...)

	// Create BC Fyne
	f := bcfynego.NewBCFyne(a, w)

	logo := f.GetLogo()

	address := widget.NewEntry()
	address.SetPlaceHolder("Channel")
	block := ui.NewBlockView()

	w.SetContent(container.NewBorder(logo, nil, nil, nil, container.NewBorder(
		container.NewBorder(nil, nil, nil, widget.NewHBox(
			widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
				go func() {
					parts := strings.Split(address.Text, "/")
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
				}()
			}),

			widget.NewButtonWithIcon("", data.NewPrimaryThemedResource(data.AccountIcon), func() {
				go f.ShowAccount(c)
			}),
			widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
				go settings(f, c)
			}),
		), address),
		nil,
		nil,
		nil,
		widget.NewScrollContainer(block),
	)))
	w.Resize(fyne.NewSize(800, 600))
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
	form.Append("Peers", widget.NewVBox(
		peerList,
		container.NewGridWithColumns(2,
			widget.NewButton("Add", func() {
				dialog.ShowEntryDialog("Add Peer", "Enter Peer Domain", func(peer string) {
					c.SetPeers(append(c.Peers, peer)...)
					form.Refresh()
				}, f.Window)
			}),
			widget.NewButton("Reset", func() {
				dialog.ShowConfirm("Reset Peers", "Reset peers to default", func(reset bool) {
					if !reset {
						return
					}
					c.SetPeers(bcgo.GetBCHost())
					form.Refresh()
				}, f.Window)
			}),
		),
	))

	form.Append("Cache", widget.NewVBox(
		ui.NewCacheView(func() bcgo.Cache {
			h, err := c.GetCache()
			if err != nil {
				f.ShowError(err)
				return nil
			}
			return h
		}),
	))

	form.Append("Network", widget.NewVBox(
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
