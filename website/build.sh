#!/usr/bin/env bash
set -o errexit #abort if any command fails

if [[ "$(docker images -q slate_app:latest 2> /dev/null)" == "" ]]; then
  echo "Building docker image..."
  docker build -t slate_app:latest .
fi

echo "Building documentation..."

docker run --rm -v $PWD:/usr/src/app/source -w /usr/src/app/source slate_app bundle exec middleman build --clean

echo "Documentation built successfully."