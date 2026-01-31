#!/bin/bash

# VibeSync: Host-Level Security Hardening Script (v1.0)
# This script applies non-intrusive security hardening using free/lightweight OS tools.

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi

echo "🛡️ VibeSync: Initiating Host Hardening..."

# 1. Kernel Hardening (sysctl)
echo "Configuring sysctl for network hardening..."
cat <<EOF > /etc/sysctl.d/99-vibesync-hardening.conf
# Ignore ICMP redirects (prevent MITM)
net.ipv4.conf.all.accept_redirects = 0
net.ipv6.conf.all.accept_redirects = 0
net.ipv4.conf.all.send_redirects = 0

# Ignore source-routed packets
net.ipv4.conf.all.accept_source_route = 0
net.ipv6.conf.all.accept_source_route = 0

# Enable reverse path filtering (prevent spoofing)
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1

# TCP SYN Cookies (DoS protection)
net.ipv4.tcp_syncookies = 1

# Restrict dmesg to root
kernel.dmesg_restrict = 1
EOF
sysctl -p /etc/sysctl.d/99-vibesync-hardening.conf

# 2. Firewall Lockdown (UFW)
if command -v ufw > /dev/null; then
    echo "Configuring UFW for localhost-only bridge traffic..."
    ufw default deny incoming
    ufw default allow outgoing
    # Allow bridge communication on loopback
    ufw allow in on lo
    # Explicitly allow the bridge ports on localhost
    ufw allow from 127.0.0.1 to any port 8085 proto tcp   # Unity
    ufw allow from 127.0.0.1 to any port 22000 proto tcp  # Blender
    ufw allow from 127.0.0.1 to any port 8000:30000 proto tcp # Multi-Engine Multiplexer Range
    echo "UFW configured for multi-engine sync. Range: 8000-30000."
    echo "UFW configured. Reminder: 'ufw enable' to activate."
else
    echo "UFW not found. Skipping firewall configuration."
fi

# 3. Suggest Lightweight Packages
echo "Checking for recommended lightweight security packages..."
PACKAGES=("apparmor" "fail2ban" "tini" "libcap")
for pkg in "${PACKAGES[@]}"; do
    if ! command -v "$pkg" > /dev/null && ! dpkg -l "$pkg" > /dev/null 2>&1; then
        echo "Tip: Consider installing '$pkg' for enhanced isolation (apt install $pkg)."
    fi
done

# 4. AppArmor Profile (Placeholder logic)
if [ -d /etc/apparmor.d/ ]; then
    echo "AppArmor detected. System is ready for process isolation."
fi

echo "✅ VibeSync: Host Hardening Applied. Operations are NOT blocked."
