#!/bin/bash

export GO111MODULE=on
(trap 'kill 0' SIGINT; go run github.com/tikz/bcov -web & cd web && npm start dev)