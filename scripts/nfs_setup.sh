#!/bin/bash

#Check if script running by root
if [ "$EUID" -ne 0 ]; then
    echo "This script must be run by root privilege"
    exit 1
fi

echo "Installing nfs-kernel-server..."
apt update
apt install nfs-kernel-server

if systemctl is-active --quiet nfs-server; then
    echo "nfs-kernel-server successfully installed"
else
    echo "Failed to install nfs-kernel-server"
    exit 1
fi

echo "Root permission for root directory..."
chmod 777 /

echo "Add current user to sudo group..."
usermod -aG sudo "$USER"

if groups "$USER" | grep -q '\bsudo\b'; then
    echo "User $USER successfully added to sudo group"
else 
    echo "Error during try add $USER to sudo group"
    exit 1
fi

exit 0