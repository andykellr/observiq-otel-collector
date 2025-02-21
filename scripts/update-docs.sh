#!/bin/sh
# Copyright  observIQ, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


TARGET_VERSION=$1
if [ -z "$TARGET_VERSION" ]; then
    echo "Must specify a target version"
    exit 1
fi

read -r -d '' DOC_FILES << EOF
docs/processors.md
docs/extensions.md
docs/exporters.md
docs/receivers.md
exporter/googlecloudexporter/README.md
EOF

for doc in $DOC_FILES
do
    echo "$doc"
    # Point contrib links to new version
    sed -i '' -Ee \
        "s|https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/v[^/]*|https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/$TARGET_VERSION|" \
        "$doc"
    # Point core links to new version
    sed -i '' -Ee \
        "s|https://github.com/open-telemetry/opentelemetry-collector/blob/v[^/]*|https://github.com/open-telemetry/opentelemetry-collector/blob/$TARGET_VERSION|" \
        "$doc"
done
