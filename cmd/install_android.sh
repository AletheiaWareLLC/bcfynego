#!/bin/bash
#
# Copyright 2020 Aletheia Ware LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e
set -x

(cd $GOPATH/src/github.com/AletheiaWareLLC/bcfynego/ui/data/ && ./gen.sh)
go fmt $GOPATH/src/github.com/AletheiaWareLLC/{bcfynego,bcfynego/...}
go vet $GOPATH/src/github.com/AletheiaWareLLC/{bcfynego,bcfynego/...}
go test $GOPATH/src/github.com/AletheiaWareLLC/{bcfynego,bcfynego/...}
ANDROID_NDK_HOME=${ANDROID_HOME}/ndk-bundle/
(cd $GOPATH/src/github.com/AletheiaWareLLC/bcfynego/cmd && fyne package -os android -appID com.aletheiaware.bc -icon $GOPATH/src/github.com/AletheiaWareLLC/bcfynego/ui/data/logo.png -name BC_unaligned)
(cd $GOPATH/src/github.com/AletheiaWareLLC/bcfynego/cmd && ${ANDROID_HOME}/build-tools/28.0.3/zipalign -f 4 BC_unaligned.apk BC.apk)
(cd $GOPATH/src/github.com/AletheiaWareLLC/bcfynego/cmd && adb install -r -g BC.apk)
#(cd $GOPATH/src/github.com/AletheiaWareLLC/bcfynego/cmd && adb logcat com.aletheiaware.bc:V org.golang.app:V *:S | tee android.log)
(cd $GOPATH/src/github.com/AletheiaWareLLC/bcfynego/cmd && adb logcat -c && adb logcat | tee android.log)
