#!/usr/bin/env bash

cd web
go build -o web
cd ../auth
go build -o auth
cd ../backend
go build -o backend