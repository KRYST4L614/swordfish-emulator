#!/bin/sh

path="/${PWD}"
# make gitlab-ci path=${path}
docker run -d \
  --name gitlab-runner \
  --restart always \
  -v "${path}":"${path}" \
  -v //var/run/docker.sock://var/run/docker.sock \
  gitlab/gitlab-runner:latest
docker exec -it -w "${path}" gitlab-runner git config --global --add safe.directory "*"
docker exec -it -w "${path}" gitlab-runner gitlab-runner exec docker lint_code
docker exec -it -w "${path}" gitlab-runner gitlab-runner exec docker unit_tests
docker exec -it -w "${path}" gitlab-runner gitlab-runner exec docker build
