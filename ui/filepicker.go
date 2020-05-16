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
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
	"log"
)

var _ desktop.Cursorable = (*filePicker)(nil)
var _ fyne.Tappable = (*filePicker)(nil)
var _ fyne.Widget = (*filePicker)(nil)

type filePicker struct {
	widget.BaseWidget
	window fyne.Window
	icon   *canvas.Image
	entry  *widget.Entry
}

func newFilePicker(w fyne.Window, e *widget.Entry) *filePicker {
	pr := &filePicker{
		window: w,
		icon:   canvas.NewImageFromResource(theme.FolderIcon()),
		entry:  e,
	}
	pr.ExtendBaseWidget(pr)
	return pr
}

func (r *filePicker) CreateRenderer() fyne.WidgetRenderer {
	return &filePickerRenderer{
		icon:  r.icon,
		entry: r.entry,
	}
}

func (r *filePicker) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (r *filePicker) Tapped(*fyne.PointEvent) {
	log.Println("filePicker.Tapped")
	// Show open file dialog
	dialog.ShowFileOpen(func(reader fyne.FileReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, r.window)
			return
		}
		if reader == nil {
			return
		}
		// Set entry text to file uri
		r.entry.SetText(reader.URI().String())
		reader.Close()
	}, r.window)
}

var _ fyne.WidgetRenderer = (*filePickerRenderer)(nil)

type filePickerRenderer struct {
	entry *widget.Entry
	icon  *canvas.Image
}

func (r *filePickerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *filePickerRenderer) Destroy() {
}

func (r *filePickerRenderer) Layout(size fyne.Size) {
	r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	r.icon.Move(fyne.NewPos((size.Width-theme.IconInlineSize())/2, (size.Height-theme.IconInlineSize())/2))
}

func (r *filePickerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
}

func (r *filePickerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.icon}
}

func (r *filePickerRenderer) Refresh() {
	canvas.Refresh(r.icon)
}
