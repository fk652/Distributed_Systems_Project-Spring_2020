#!/usr/bin/env bash

./auth.sh & ./backend.sh & ./web.sh && kill $;
