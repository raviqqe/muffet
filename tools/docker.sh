#!/bin/sh

set -e

check_buildx(){
    if docker buildx version &>/dev/null; then
        return 0
    fi
    if docker buildx ls &>/dev/null; then
        return 0
    fi
    return 1
}

build_multiarch() {
    if ! check_buildx; then
        echo "Error: Docker buildx is not installed or not working"
        return 1
    fi

    local platforms=${1:-"linux/amd64"}
    
    docker buildx create --name multiarch --use || true
    
    echo "Building for platforms: $platforms"
    
    # Build for each git tag and latest
    for tag in latest $(git tag --points-at | sed s/^v//); do
        docker buildx build \
            --platform "$platforms" \
            --tag raviqqe/muffet:$tag \
            --push \
            .
    done
}


build_multiarch "$@"