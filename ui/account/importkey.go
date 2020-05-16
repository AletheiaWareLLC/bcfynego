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

package account

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type ImportKey struct {
	Alias           *widget.Entry
	Access          *widget.Entry
	ImportKeyButton *widget.Button
}

func NewImportKey() *ImportKey {
	i := &ImportKey{
		Alias:  widget.NewEntry(),
		Access: widget.NewEntry(),
		ImportKeyButton: &widget.Button{
			Style: widget.PrimaryButton,
			Text:  "Import Key",
		},
	}
	i.Alias.SetPlaceHolder("Alias")
	i.Access.SetPlaceHolder("Access Code")
	// TODO Alias is single line, handle enter key by moving to password
	// TODO Access is single line, handle enter key by moving to button/auto click
	return i
}

func (i *ImportKey) CanvasObject() fyne.CanvasObject {
	return fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		i.Alias,
		i.Access,
		layout.NewSpacer(),
		i.ImportKeyButton,
	)
}
