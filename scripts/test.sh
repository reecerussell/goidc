#!/bin/bash

export CLIENTS_TABLE_NAME=goidc-clients-test
export USERS_TABLE_NAME=goidc-users-test

go test ./...