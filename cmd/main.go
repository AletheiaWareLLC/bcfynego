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
	"flag"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcclientgo"
	"github.com/AletheiaWareLLC/bcfynego"
	"github.com/AletheiaWareLLC/bcgo"
	"log"
)

var peer = flag.String("peer", "", "BC peer")

func main() {
	// Create Application
	a := app.New()

	// Create Window
	w := a.NewWindow("BC")
	w.SetMaster()

	peers := bcgo.SplitRemoveEmpty(*peer, ",")
	if len(peers) == 0 {
		peers = append(peers,
			bcgo.GetBCHost(), // Add BC host as peer
		)
	}

	// Create BC Client
	c := &bcclientgo.BCClient{
		Peers: peers,
	}

	// Create BC Fyne
	f := &bcfynego.BCFyne{
		App:    a,
		Window: w,
	}

	logo := f.GetLogo()

	w.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(logo, nil, nil, nil), logo, widget.NewScrollContainer(widget.NewAccordionContainer(
		widget.NewAccordionItem("Node", widget.NewVBox(
			widget.NewButton("Node", func() {
				go func() {
					n, err := f.GetNode(c)
					if err != nil {
						f.ShowError(err)
					} else {
						f.ShowNode(n)
					}
				}()
			}),
			widget.NewButton("NewNode", func() {
				log.Println("// TODO go c.NewNode()")
			}),
			widget.NewButton("ExistingNode", func() {
				log.Println("// TODO go c.ExistingNode()")
			}),
			widget.NewButton("SetNode", func() {
				log.Println("// TODO go c.SetNode()")
			}),
		)),
		widget.NewAccordionItem("Alias", widget.NewVBox(
			widget.NewButton("Register", func() {
				log.Println("// TODO go c.Register()")
			}),
			widget.NewButton("Alias", func() {
				log.Println("// TODO go c.Alias()")
			}),
		)),
		widget.NewAccordionItem("Account", widget.NewVBox(
			widget.NewButton("ShowAccount", func() {
				go f.ShowAccount(c)
			}),
			widget.NewButton("ShowAccessDialog", func() {
				log.Println("// TODO go c.ShowAccessDialog()")
			}),
			widget.NewButton("ImportKeys", func() {
				log.Println("// TODO go c.ImportKeys()")
			}),
			widget.NewButton("ExportKeys", func() {
				log.Println("// TODO go c.ExportKeys()")
			}),
		)),
		widget.NewAccordionItem("ShowError", widget.NewButton("ShowError", func() {
			go f.ShowError(fmt.Errorf("Sample Error"))
		})),
		widget.NewAccordionItem("Root", widget.NewVBox(
			widget.NewButton("GetRoot", func() {
				go log.Println(c.GetRoot())
			}),
			widget.NewButton("SetRoot", func() {
				log.Println("// TODO go log.Println(c.SetRoot())")
			}),
		)),
		widget.NewAccordionItem("Peers", widget.NewVBox(
			widget.NewButton("GetDefaultPeers", func() {
				go func() {
					ps, err := c.GetDefaultPeers()
					if err != nil {
						f.ShowError(err)
					} else {
						log.Println(ps)
					}
				}()
			}),
			widget.NewButton("GetPeers", func() {
				go func() {
					ps, err := c.GetPeers()
					if err != nil {
						f.ShowError(err)
					} else {
						log.Println(ps)
					}
				}()
			}),
			widget.NewButton("SetPeers", func() {
				log.Println("// TODO go log.Println(c.SetPeers())")
			}),
		)),
		widget.NewAccordionItem("Cache", widget.NewVBox(
			widget.NewButton("GetCache", func() {
				go func() {
					ch, err := c.GetCache()
					if err != nil {
						f.ShowError(err)
					} else {
						log.Println(ch)
					}
				}()
			}),
			widget.NewButton("SetCache", func() {
				log.Println("// TODO go go log.Println(c.SetCache())")
			}),
			widget.NewButton("Purge", func() {
				log.Println("// TODO go c.Purge()")
			}),
		)),
		widget.NewAccordionItem("Network", widget.NewVBox(
			widget.NewButton("GetNetwork", func() {
				go func() {
					n, err := c.GetNetwork()
					if err != nil {
						f.ShowError(err)
					} else {
						log.Println(n)
					}
				}()
			}),
			widget.NewButton("SetNetwork", func() {
				log.Println("// TODO go go log.Println(c.SetNetwork())")
			}),
			widget.NewButton("Pull", func() {
				log.Println("// TODO go c.Pull()")
			}),
			widget.NewButton("Push", func() {
				log.Println("// TODO go c.Push()")
			}),
		)),
		widget.NewAccordionItem("Channel", widget.NewVBox(
			widget.NewButton("Head", func() {
				log.Println("// TODO go c.Head()")
			}), widget.NewButton("Block", func() {
				log.Println("// TODO go c.Block()")
			}), widget.NewButton("Record", func() {
				log.Println("// TODO go c.Record()")
			}), widget.NewButton("Read", func() {
				log.Println("// TODO go c.Read()")
			}), widget.NewButton("ReadKey", func() {
				log.Println("// TODO go c.ReadKey()")
			}), widget.NewButton("ReadPayload", func() {
				log.Println("// TODO go c.ReadPayload()")
			}), widget.NewButton("Write", func() {
				log.Println("// TODO go c.Write()")
			}), widget.NewButton("Mine", func() {
				log.Println("// TODO go c.Mine()")
			}),
		)),
	))))
	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()
	w.ShowAndRun()
}
