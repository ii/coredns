#!/usr/bin/env bash

set -e +o pipefail

if [ "$TRAVIS_PULL_REQUEST" != "false" ] ; then
   cat .benchmark.log
fi
