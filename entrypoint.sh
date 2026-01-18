#!/bin/sh
set -e

# =============================================================================
# WebSSH Entrypoint with Tailscale Support
# =============================================================================
# 如果提供了 TS_AUTHKEY 环境变量，则启动 Tailscale 以获得 IPv4/IPv6 双栈出站能力
# 这对于只有 IPv6 出站的 PaaS 平台（如 Koyeb）特别有用
# =============================================================================

start_tailscale() {
    echo "[Tailscale] Starting tailscaled in userspace mode..."
    tailscaled --tun=userspace-networking --statedir=/var/lib/tailscale &
    
    # 等待 tailscaled 启动
    sleep 2
    
    echo "[Tailscale] Authenticating with authkey..."
    tailscale up --authkey="${TS_AUTHKEY}" --hostname="${TS_HOSTNAME:-webssh}"
    
    echo "[Tailscale] Connected! Tailscale IP:"
    tailscale ip -4 2>/dev/null || echo "(no IPv4)"
    tailscale ip -6 2>/dev/null || echo "(no IPv6)"
}

# 如果提供了 Tailscale Auth Key，则启动 Tailscale
if [ -n "${TS_AUTHKEY}" ]; then
    start_tailscale
else
    echo "[Tailscale] TS_AUTHKEY not set, skipping Tailscale initialization."
    echo "[Tailscale] To enable IPv4 outbound on IPv6-only platforms, set TS_AUTHKEY."
fi

# 启动 WebSSH
echo "[WebSSH] Starting WebSSH server..."
exec ./webssh
