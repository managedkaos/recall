docker run --rm --env GITHUB_TOKEN="$GITHUB_TOKEN" --volume "$PWD:/work" ghcr.io/managedkaos/get-action-version-number:main --update-in-place --workflow go.yml
docker run --rm --env GITHUB_TOKEN="$GITHUB_TOKEN" --volume "$PWD:/work" ghcr.io/managedkaos/get-action-version-number:main --update-in-place --workflow release-build.yml
docker run --rm --env GITHUB_TOKEN="$GITHUB_TOKEN" --volume "$PWD:/work" ghcr.io/managedkaos/get-action-version-number:main --update-in-place --workflow test-binaries.yml

