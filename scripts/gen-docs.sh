#!/bin/bash

prefix="github.com/vanilla-os/sdk"
versions=$(ls pkg)

# Generate references for each module
for pkg_version in $versions; do
  modules=$(ls pkg/$pkg_version)

  mkdir -p docs/references/$pkg_version
  echo "# $pkg_version\n" > docs/references/$pkg_version/README.md

  for module in $modules; do
    if [ $module = "VERSION" ]; then
      continue
    fi

    echo "Generating refs for $prefix/pkg/$pkg_version/$module"
    godoc2md $prefix/pkg/$pkg_version/$module > docs/references/$pkg_version/$module.md

    echo "- [$module](/references/$pkg_version/$module.md)" >> docs/references/$pkg_version/README.md
  done
done