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

type Authentication struct {
	Alias              *widget.Label
	Password           *widget.Entry
	AuthenticateButton *widget.Button
}

func NewAuthentication(alias string) *Authentication {
	s := &Authentication{
		Alias:    widget.NewLabel(alias),
		Password: widget.NewPasswordEntry(),
		AuthenticateButton: &widget.Button{
			Style: widget.PrimaryButton,
			Text:  "Authenticate",
		},
	}
	s.Password.SetPlaceHolder("Password")
	// TODO Password is single line, handle enter key by moving to button/auto click
	return s
}

func (s *Authentication) CanvasObject() fyne.CanvasObject {
	return fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		s.Alias,
		s.Password,
		layout.NewSpacer(),
		s.AuthenticateButton,
	)
}
