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
	"fmt"
	"fyne.io/fyne/widget"
)

type CacheView struct {
	widget.Label
	updater func() bcgo.Cache
}

func NewCacheView(updater func() bcgo.Cache) *CacheView {
	v := &CacheView{
		updater: updater,
	}
	v.ExtendBaseWidget(v)
	v.update()
	return v
}

func (v *CacheView) Refresh() {
	v.update()
	v.Label.Refresh()
}

func (v *CacheView) update() {
	if u := v.updater; u != nil {
		if c := u(); c != nil {
			v.Text = fmt.Sprintf("%T %v", c, c)
		}
	}
}
