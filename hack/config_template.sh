#!/bin/bash

if test -f "${HOME}/.cloudagent.conf.yaml"; then
    exit 0
fi

cat << 'EOF' > ~/.cloudagent.conf.yaml
persistentSpec: "@every 60s"
gardenerProjects: []
EOF
