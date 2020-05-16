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

package data

import (
	"encoding/hex"
	"fmt"
	"fyne.io/fyne"
	"image/color"
	"strconv"
	"strings"
)

type ThemedResource struct {
	source fyne.Resource
}

func NewThemedResource(source fyne.Resource) *ThemedResource {
	return &ThemedResource{
		source: source,
	}
}

func (res *ThemedResource) Name() string {
	// TODO log.Println("ThemedResource.Name")
	return res.source.Name()
}

func (res *ThemedResource) Content() []byte {
	// TODO log.Println("ThemedResource.Content")
	background := colorToHexString(fyne.CurrentApp().Settings().Theme().BackgroundColor())
	icon := colorToHexString(fyne.CurrentApp().Settings().Theme().IconColor())
	primary := "#0fffff" //colorToHexString(fyne.CurrentApp().Settings().Theme().PrimaryColor())
	textsize := fyne.CurrentApp().Settings().Theme().TextSize()
	svg := string(res.source.Content())
	svg = strings.ReplaceAll(svg, "background", background)
	svg = strings.ReplaceAll(svg, "flasksize", strconv.Itoa(textsize))
	svg = strings.ReplaceAll(svg, "icon", icon)
	svg = strings.ReplaceAll(svg, "primary", primary)
	svg = strings.ReplaceAll(svg, "textsize", strconv.Itoa(textsize*2))
	// TODO cache until theme changes
	return []byte(svg)
}

func colorToHexString(color color.Color) string {
	r, g, b, a := color.RGBA()
	cBytes := []byte{byte(r), byte(g), byte(b), byte(a)}
	return fmt.Sprintf("#%s", hex.EncodeToString(cBytes))
}
