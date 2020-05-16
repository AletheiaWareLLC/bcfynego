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

package ui_test

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/AletheiaWareLLC/bcfynego/ui/account"
	"github.com/AletheiaWareLLC/bcfynego/ui/data"
	"testing"
)

func Test_UI(t *testing.T) {
	a := test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		builder func(fyne.Window) fyne.CanvasObject
	}{
		"account_import_key": {
			builder: func(w fyne.Window) fyne.CanvasObject {
				return account.NewImportKey().CanvasObject()
			},
		},
		"account_sign_in": {
			builder: func(w fyne.Window) fyne.CanvasObject {
				return account.NewSignIn().CanvasObject()
			},
		},
		"account_sign_up": {
			builder: func(w fyne.Window) fyne.CanvasObject {
				return account.NewSignUp().CanvasObject()
			},
		},
		"logo": {
			builder: func(w fyne.Window) fyne.CanvasObject {
				img := &canvas.Image{
					Resource: data.NewThemedResource(data.Logo),
					FillMode: canvas.ImageFillOriginal,
				}
				img.SetMinSize(fyne.NewSize(480, 240))
				return img
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			window := a.NewWindow(name)
			object := tt.builder(window)
			window.SetContent(fyne.NewContainerWithLayout(layout.NewCenterLayout(), object))
			window.Resize(object.MinSize().Max(fyne.NewSize(400, 300)))
			test.AssertImageMatches(t, name+".png", window.Canvas().Capture())
			window.Close()
		})
	}
}
