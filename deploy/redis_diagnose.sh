#!/bin/bash

echo "==== Redis 本地连接诊断脚本 ===="

# 检查 Redis 进程是否存在
echo -e "\n[1] 检查 Redis 进程状态..."
pgrep -x redis-server > /dev/null
if [ $? -eq 0 ]; then
    echo "✅ Redis 正在运行"
else
    echo "❌ Redis 未运行，请执行：sudo systemctl start redis"
    exit 1
fi

# 显示监听端口
echo -e "\n[2] 检查 Redis 监听端口..."
ss -ltnp | grep 6379 || echo "❌ Redis 未监听在 6379 端口"

# 获取配置文件路径（默认）
CONF="/etc/redis/redis.conf"
if [ ! -f "$CONF" ]; then
    echo "❌ Redis 配置文件未找到：$CONF"
    exit 1
fi

# 检查配置项
echo -e "\n[3] 检查 redis.conf 配置项..."
echo -n "bind: "; grep -E '^bind ' $CONF || echo "(未设置，可能监听所有地址)"
echo -n "port: "; grep -E '^port ' $CONF
echo -n "protected-mode: "; grep -E '^protected-mode ' $CONF
echo -n "requirepass: "; grep -E '^requirepass ' $CONF || echo "(未设置密码)"

# 本地连接测试
echo -e "\n[4] 尝试本地连接 Redis..."
redis-cli ping > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ 本地连接成功 (PING -> PONG)"
else
    echo "⚠️ 连接失败，尝试获取错误信息..."
    redis-cli ping
fi

# 推荐动作
echo -e "\n[5] 建议动作："
echo "- 确保 bind 包含 127.0.0.1 或 0.0.0.0"
echo "- 确保 protected-mode 为 yes，避免 0.0.0.0 + 无密码造成安全风险"
echo "- 如果 requirepass 设置了密码，使用 redis-cli -a <password> 连接"
echo "- 如果仍有问题，查看日志：sudo journalctl -u redis 或 /var/log/redis/redis-server.log"

echo -e "\n=== 诊断完毕 ==="