#!/bin/bash

prefix="github.com/vanilla-os/sdk"

for pkg_version in $(ls pkg); do
  for module in $(ls pkg/$pkg_version); do
    if [ $module = "VERSION" ]; then
      continue
    fi

    echo "Generating refs for $prefix/pkg/$pkg_version/$module"
    mkdir -p docs/references/$pkg_version
    godoc2md $prefix/pkg/$pkg_version/$module > docs/references/$pkg_version/$module.md
    done
done
