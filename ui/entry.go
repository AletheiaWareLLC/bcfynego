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
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"sort"
)

type EntryView struct {
	widget.Form
	hash                 *widget.Label
	timestamp            *widget.Label
	creator              *widget.Label
	access               *widget.Box
	payload              *widget.Label
	compressionAlgorithm *widget.Label
	encryptionAlgorithm  *widget.Label
	signature            *widget.Label
	signatureAlgorithm   *widget.Label
	reference            *widget.Box
	meta                 *widget.Box
}

func NewEntryView() *EntryView {
	v := &EntryView{
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
		creator: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		access: widget.NewVBox(),
		payload: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		compressionAlgorithm: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		encryptionAlgorithm: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		signature: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		signatureAlgorithm: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		},
		reference: widget.NewVBox(),
		meta:      widget.NewVBox(),
	}
	v.ExtendBaseWidget(v)
	v.hash.ExtendBaseWidget(v.hash)
	v.timestamp.ExtendBaseWidget(v.timestamp)
	v.creator.ExtendBaseWidget(v.creator)
	v.access.ExtendBaseWidget(v.access)
	v.payload.ExtendBaseWidget(v.payload)
	v.compressionAlgorithm.ExtendBaseWidget(v.compressionAlgorithm)
	v.encryptionAlgorithm.ExtendBaseWidget(v.encryptionAlgorithm)
	v.signature.ExtendBaseWidget(v.signature)
	v.signatureAlgorithm.ExtendBaseWidget(v.signatureAlgorithm)
	v.reference.ExtendBaseWidget(v.reference)
	v.meta.ExtendBaseWidget(v.meta)
	v.Append("Hash", v.hash)
	v.Append("Timestamp", v.timestamp)
	v.Append("Creator", v.creator)
	v.Append("Access", v.access)
	v.Append("Payload", v.payload)
	v.Append("Compression", v.compressionAlgorithm)
	v.Append("Encryption", v.encryptionAlgorithm)
	v.Append("Signature", v.signature)
	v.Append("Signature", v.signatureAlgorithm)
	v.Append("References", v.reference)
	v.Append("Metadata", v.meta)
	v.Hide()
	return v
}

func (v *EntryView) SetHash(hash []byte) {
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

func (v *EntryView) SetRecord(entry *bcgo.Record) {
	if entry == nil {
		v.Hide()
		return
	}
	v.timestamp.SetText(bcgo.TimestampToString(entry.Timestamp))
	v.creator.SetText(entry.Creator)
	var accesses []fyne.CanvasObject
	for _, a := range entry.Access {
		v := NewAccessView()
		v.SetAccess(a)
		accesses = append(accesses, v)
	}
	v.access.Children = accesses
	v.access.Refresh()
	v.payload.SetText(base64.RawURLEncoding.EncodeToString(entry.Payload))
	v.compressionAlgorithm.SetText(entry.CompressionAlgorithm.String())
	v.encryptionAlgorithm.SetText(entry.EncryptionAlgorithm.String())
	v.signature.SetText(base64.RawURLEncoding.EncodeToString(entry.Signature))
	v.signatureAlgorithm.SetText(entry.SignatureAlgorithm.String())
	var references []fyne.CanvasObject
	for _, r := range entry.Reference {
		v := NewReferenceView()
		v.SetReference(r)
		references = append(references, v)
	}
	v.reference.Children = references
	v.reference.Refresh()
	keys := make([]string, len(entry.Meta))
	i := 0
	for k := range entry.Meta {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	var metas []fyne.CanvasObject
	for _, k := range keys {
		metas = append(metas, container.NewGridWithColumns(2,
			&widget.Label{
				Text: k,
				TextStyle: fyne.TextStyle{
					Monospace: true,
				},
			},
			&widget.Label{
				Text: entry.Meta[k],
				TextStyle: fyne.TextStyle{
					Monospace: true,
				},
			},
		))
	}
	v.meta.Children = metas
	v.meta.Refresh()
	if v.Visible() {
		v.Refresh()
	} else {
		v.Show()
	}
}
