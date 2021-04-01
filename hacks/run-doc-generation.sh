#!/bin/bash

cd ./cli
go build \
  -o ./bin/kubectl-shovel \
  .

HOME="/home/user" ./bin/kubectl-shovel doc
