#!/bin/sh

path="/${PWD}"
# make gitlab-ci path=${path}
docker volume create gitlab-runner-config
docker run -d \
  --name gitlab-runner \
  --restart always \
  -v "${path}"://app/ \
  -v //var/run/docker.sock://var/run/docker.sock \
  -v gitlab-runner-config:/etc/gitlab-runner \
  gitlab/gitlab-runner:alpine
docker exec -it -w //app/ gitlab-runner git config --global --add safe.directory "*"
docker exec -it -w //app/ gitlab-runner gitlab-runner exec docker lint_code
docker exec -it -w //app/ gitlab-runner gitlab-runner exec docker unit_tests
docker exec -it -w //app/ gitlab-runner gitlab-runner exec docker build
