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
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type BlockView struct {
	widget.Form
	hash      *widget.Label
	timestamp *widget.Label
	channel   *widget.Label
	length    *widget.Label
	previous  *widget.Label
	miner     *widget.Label
	nonce     *widget.Label
	entry     *widget.Box
}

func NewBlockView() *BlockView {
	v := &BlockView{
		hash: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
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
		length: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		previous: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		miner: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		nonce: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		entry: widget.NewVBox(),
	}
	v.ExtendBaseWidget(v)
	v.hash.ExtendBaseWidget(v.hash)
	v.timestamp.ExtendBaseWidget(v.timestamp)
	v.channel.ExtendBaseWidget(v.channel)
	v.length.ExtendBaseWidget(v.length)
	v.previous.ExtendBaseWidget(v.previous)
	v.miner.ExtendBaseWidget(v.miner)
	v.nonce.ExtendBaseWidget(v.nonce)
	v.entry.ExtendBaseWidget(v.entry)
	v.Append("Hash", v.hash)
	v.Append("Timestamp", v.timestamp)
	v.Append("Channel", v.channel)
	v.Append("Length", v.length)
	v.Append("Previous", v.previous)
	v.Append("Miner", v.miner)
	v.Append("Nonce", v.nonce)
	v.Append("Entries", v.entry)
	v.Hide()
	return v
}

func (v *BlockView) SetHash(hash []byte) {
	if hash == nil {
		v.Hide()
		return
	}
	v.hash.SetText(base64.RawURLEncoding.EncodeToString(hash))
	if v.Visible() {
		v.Refresh()
	} else {
		v.Show()
	}
}

func (v *BlockView) SetBlock(block *bcgo.Block) {
	if block == nil {
		v.Hide()
		return
	}
	v.timestamp.SetText(bcgo.TimestampToString(block.Timestamp))
	v.channel.SetText(block.ChannelName)
	v.length.SetText(fmt.Sprintf("%d", block.Length))
	v.previous.SetText(base64.RawURLEncoding.EncodeToString(block.Previous))
	v.miner.SetText(block.Miner)
	v.nonce.SetText(fmt.Sprintf("%d", block.Nonce))
	var entries []fyne.CanvasObject
	for _, e := range block.Entry {
		v := NewEntryView()
		v.SetHash(e.RecordHash)
		v.SetRecord(e.Record)
		entries = append(entries, v)
	}
	v.entry.Children = entries
	v.entry.Refresh()
	if v.Visible() {
		v.Refresh()
	} else {
		v.Show()
	}
}
