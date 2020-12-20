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
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcgo"
)

type ReferenceView struct {
	widget.Form
	timestamp *widget.Label
	channel   *widget.Label
	block     *widget.Label
	record    *widget.Label
	index     *widget.Label
}

func NewReferenceView() *ReferenceView {
	v := &ReferenceView{
		timestamp: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		channel: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		block: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		record: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		index: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
	}
	v.ExtendBaseWidget(v)
	v.timestamp.ExtendBaseWidget(v.timestamp)
	v.channel.ExtendBaseWidget(v.channel)
	v.block.ExtendBaseWidget(v.block)
	v.record.ExtendBaseWidget(v.record)
	v.index.ExtendBaseWidget(v.index)
	v.Append("Timestamp", v.timestamp)
	v.Append("Channel", v.channel)
	v.Append("Block", v.block)
	v.Append("Record", v.record)
	v.Append("Index", v.index)
	v.Hide()
	return v
}

func (v *ReferenceView) SetReference(reference *bcgo.Reference) {
	if reference == nil {
		v.Hide()
		return
	}
	v.timestamp.SetText(bcgo.TimestampToString(reference.Timestamp))
	v.channel.SetText(reference.ChannelName)
	v.block.SetText(base64.RawURLEncoding.EncodeToString(reference.BlockHash))
	v.record.SetText(base64.RawURLEncoding.EncodeToString(reference.RecordHash))
	v.index.SetText(fmt.Sprintf("%d", reference.Index))
	if v.Visible() {
		v.Refresh()
	} else {
		v.Show()
	}
}
