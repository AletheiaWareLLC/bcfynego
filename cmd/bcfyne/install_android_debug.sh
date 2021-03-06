#!/bin/bash
#
# Copyright 2020-2021 Aletheia Ware LLC
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

(cd $GOPATH/src/aletheiaware.com/bcfynego/ui/data/ && ./gen.sh)
go fmt $GOPATH/src/aletheiaware.com/bcfynego/...
go vet $GOPATH/src/aletheiaware.com/bcfynego/...
go test $GOPATH/src/aletheiaware.com/bcfynego/...
(cd $GOPATH/src/aletheiaware.com/bcfynego/cmd/bcfyne && fyne package -os android -appID com.aletheiaware.bc -name BC_unaligned)
(cd $GOPATH/src/aletheiaware.com/bcfynego/cmd/bcfyne && ${ANDROID_HOME}/build-tools/28.0.3/zipalign -f 4 BC_unaligned.apk BC.apk)
(cd $GOPATH/src/aletheiaware.com/bcfynego/cmd/bcfyne && adb install -r -g BC.apk)
(cd $GOPATH/src/aletheiaware.com/bcfynego/cmd/bcfyne && adb logcat -c && adb logcat | tee android.log)
