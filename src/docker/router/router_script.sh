# Enable package forwarding in the container.
echo "net.ipv4.ip_forward=1" | tee -a /etc/sysctl.conf

LOG_FILE=/usr/src/router/router.log
touch $LOG_FILE

# Continuous monitoring of interfaces
while true; do
    # Get the list of interfaces that are not loopback
    INTERFACES=$(ip -o link show | awk -F': ' '{print $2}' | grep -v 'lo')

    for iface in $INTERFACES; do
        # Check if we have already applied iptables for this interface
        if ! iptables -t nat -C POSTROUTING -o "$iface" -j MASQUERADE 2>/dev/null; then
            # Apply NAT if not already applied
            iptables -t nat -A POSTROUTING -o "$iface" -j MASQUERADE
            iptables -A FORWARD -o "$iface" -j LOG --log-prefix "FW_FORWARD: " --log-level 4 2>&1 | tee -a "$LOG_FILE"
            echo "NAT and logging set for interface: $iface" | tee -a "$LOG_FILE"
        fi
    done
    # Sleep before checking again
    sleep 5
done