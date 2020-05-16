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

package ui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcclientgo"
	"log"
)

type Client struct {
	bcclientgo.Client
	App    fyne.App
	Window fyne.Window
}

func (c *Client) GetLogo() fyne.CanvasObject {
	log.Println("GetLogo")
	return widget.NewLabel("BC")
}

func (c *Client) ShowNode() {
	log.Println("ShowNode")
}
