#!/bin/bash
#
# Start script for emergency-auth-code-api
PORT="20188"

exec ./emergency-auth-code-api "-bind-addr=:${PORT}"