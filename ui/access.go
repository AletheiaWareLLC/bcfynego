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
	"aletheiaware.com/bcgo"
	"encoding/base64"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type AccessView struct {
	widget.Form
	alias               *widget.Label
	secretKey           *widget.Label
	encryptionAlgorithm *widget.Label
}

func NewAccessView() *AccessView {
	v := &AccessView{
		alias: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		secretKey: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		encryptionAlgorithm: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
	}
	v.ExtendBaseWidget(v)
	v.alias.ExtendBaseWidget(v.alias)
	v.secretKey.ExtendBaseWidget(v.secretKey)
	v.encryptionAlgorithm.ExtendBaseWidget(v.encryptionAlgorithm)
	v.Append("Alias", v.alias)
	v.Append("Key", v.secretKey)
	v.Append("Encryption", v.encryptionAlgorithm)
	v.Hide()
	return v
}

func (v *AccessView) SetAccess(access *bcgo.Record_Access) {
	if access == nil {
		v.Hide()
		return
	}
	v.alias.SetText(access.Alias)
	v.secretKey.SetText(base64.RawURLEncoding.EncodeToString(access.SecretKey))
	v.encryptionAlgorithm.SetText(access.EncryptionAlgorithm.String())
	if v.Visible() {
		v.Refresh()
	} else {
		v.Show()
	}
}
