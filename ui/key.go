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
	"encoding/base64"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type KeyView struct {
	widget.Label
	updater func() []byte
}

func NewKeyView(updater func() []byte) *KeyView {
	v := &KeyView{
		Label: widget.Label{
			Alignment: fyne.TextAlignLeading,
			TextStyle: fyne.TextStyle{Monospace: true},
			Wrapping:  fyne.TextTruncate,
		},
		updater: updater,
	}
	v.ExtendBaseWidget(v)
	v.update()
	return v
}

func (v *KeyView) Refresh() {
	v.update()
	v.Label.Refresh()
}

func (v *KeyView) update() {
	if u := v.updater; u != nil {
		if bytes := u(); len(bytes) != 0 {
			base := base64.RawURLEncoding.EncodeToString(bytes)
			var publicKeyRunes []rune
			for i, v := range []rune(base) {
				if i > 0 && i%64 == 0 {
					publicKeyRunes = append(publicKeyRunes, '\n')
				}
				publicKeyRunes = append(publicKeyRunes, v)
			}
			v.Text = string(publicKeyRunes)
		}
	}
}
