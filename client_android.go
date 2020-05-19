// +build android

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

package bcfynego

import (
	"log"
	"os"
)

func (c *BCFyneClient) GetRoot() (string, error) {
	if c.BCClient.Root == "" {
		log.Println("bcfynego.BCFyneClient_android.GetRoot()")
		log.Println("*******************init*******************")
		log.Println(os.Environ())
		if _, ok := os.LookupEnv("ROOT_DIRECTORY"); !ok {
			os.Setenv("ROOT_DIRECTORY", os.Getenv("FILESDIR"))
		}
		if _, ok := os.LookupEnv("CACHE_DIRECTORY"); !ok {
			os.Setenv("CACHE_DIRECTORY", os.Getenv("TMPDIR"))
		}
	}
	return c.BCClient.GetRoot()
}
