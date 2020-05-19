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
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcfynego"
)

func main() {
	// Create application
	a := app.New()

	// Create window
	w := a.NewWindow("BC")
	w.SetMaster()

	// Create BC Fyne client
	c := &bcfynego.BCFyneClient{
		App:    a,
		Window: w,
	}

	logo := c.GetLogo()

	nodeButton := widget.NewButton("Node", func() {
		go c.ShowNode()
	})

	w.SetContent(fyne.NewContainerWithLayout(layout.NewBorderLayout(logo, nil, nil, nil), logo, widget.NewAccordionContainer(
		widget.NewAccordionItem("Node", nodeButton))))
	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()
	w.ShowAndRun()
}
