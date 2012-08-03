#!/usr/bin/env bash

set -e
set -u

rsync -havP lib/pre-commit .git/hooks/
