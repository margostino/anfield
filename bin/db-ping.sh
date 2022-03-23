#!/bin/bash

mongo <<EOF
db.runCommand("ping").ok
EOF