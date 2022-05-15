#!/usr/bin/env bash

echo "Installing $1 from $2"
go install $2@latest
