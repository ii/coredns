#!/usr/bin/env bash

set -e +o pipefail

if [ "$TRAVIS_PULL_REQUEST" != "false" ] ; then
    curl -H "Authorization: token ${GITHUB_TOKEN}" -X POST \
        --data-urlencode "{\"body\": \"$(cat .benchmark.log | sed "s/\"/'/g")\"}" \
        "https://api.github.com/repos/${TRAVIS_REPO_SLUG}/issues/${TRAVIS_PULL_REQUEST}/comments"
fi
