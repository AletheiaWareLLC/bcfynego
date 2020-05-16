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

type SignUp struct {
	Alias        *widget.Entry
	Password     *widget.Entry
	Confirm      *widget.Entry
	SignUpButton *widget.Button
}

func NewSignUp() *SignUp {
	s := &SignUp{
		Alias:    widget.NewEntry(),
		Password: widget.NewPasswordEntry(),
		Confirm:  widget.NewPasswordEntry(),
		SignUpButton: &widget.Button{
			Style: widget.PrimaryButton,
			Text:  "Sign Up",
		},
	}
	s.Alias.SetPlaceHolder("Alias")
	s.Password.SetPlaceHolder("Password")
	s.Confirm.SetPlaceHolder("Confirm Password")
	return s
}

func (s *SignUp) CanvasObject() fyne.CanvasObject {
	return fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		s.Alias,
		s.Password,
		s.Confirm,
		layout.NewSpacer(),
		s.SignUpButton,
	)
}
