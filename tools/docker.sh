#!/bin/sh

set -e

check_buildx(){
    docker buildx version &>/dev/null
    return $?
}

build_image() {
    local platforms=${1:-} 
    if [ ! -z "$platforms" ] && check_buildx; then
        echo "Building for platforms: $platforms"
        # Build for each git tag and latest
        for tag in latest $(git tag --points-at | sed s/^v//); do
          docker buildx build \
             --platform "$platforms" \
             --tag raviqqe/muffet:$tag \
             --push \
             .
        done
    else
        for tag in latest $(git tag --points-at | sed s/^v//); do
          docker build -t raviqqe/muffet:$tag .
          docker push raviqqe/muffet:$tag
        done
    fi
}


build_image "$@"