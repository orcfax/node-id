#! /usr/bin/bash

# Verify the signatures of the executable files created by the goreleaser
# package.

echo "verifying signatures, to run manually: 'gpg --verify <path to file>'"
echo ""

for i in `find dist/ -name "*.sig" -type f`; do
    gpg --verify $i &> /dev/null
    if [ $? -eq 0 ]; then
      echo "signature verify successful: '$i'"
    else
      echo "signature verify failed: '$i'"
    fi
done
