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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Authentication struct {
	Alias              *widget.Label
	Password           *widget.Entry
	AuthenticateButton *widget.Button
}

func NewAuthentication(alias string) *Authentication {
	s := &Authentication{
		Alias:              widget.NewLabel(alias),
		Password:           widget.NewPasswordEntry(),
		AuthenticateButton: widget.NewButton("Authenticate", nil),
	}
	s.Password.PlaceHolder = "Password"
	s.Password.Wrapping = fyne.TextWrapOff
	s.AuthenticateButton.Importance = widget.HighImportance
	return s
}

func (s *Authentication) CanvasObject() fyne.CanvasObject {
	return container.NewGridWithColumns(1,
		s.Alias,
		s.Password,
		layout.NewSpacer(),
		s.AuthenticateButton,
	)
}
