# ############################################################################### 
# 문의처 : innodreamer@gmail.com
# 참고 사이트) https://github.com/angristan/wireguard-install
# ###############################################################################
#!/bin/bash
echo "Start the script!!"

source ./init.env

sudo rm ./barista0*

if [ "$(systemd-detect-virt)" == "openvz" ]; then
    echo "OpenVZ is not supported"
    exit
fi

if [ "$(systemd-detect-virt)" == "lxc" ]; then
    echo "LXC is not supported (yet)."
    echo "VPN can technically run in an LXC container,"
    echo "but the kernel module has to be installed on the host,"
    echo "the container has to be run with some specific parameters"
    echo "and only the tools need to be installed in the container."
    exit
fi

# Check OS version
if [[ -e /etc/debian_version ]]; then
    source /etc/os-release
    OS=$ID # debian or ubuntu
elif [[ -e /etc/fedora-release ]]; then
    OS=fedora
elif [[ -e /etc/centos-release ]]; then
    OS=centos
elif [[ -e /etc/arch-release ]]; then
    OS=arch
else
    echo "Looks like you aren't running this installer on a Debian, Ubuntu, Fedora, CentOS or Arch Linux system"
    exit 1
fi

SERVER_PUB_IPV4=$VPN_SERVER_PUB_IP
echo "SERVER_PUB_IPV4 : " "$SERVER_PUB_IPV4"

SERVER_PUB_NIC=$VPN_SERVER_PUB_NIC
echo "Public interface : " "$SERVER_PUB_NIC" 

SERVER_WG_NIC="barista0"
echo "VPN interface name : " "$SERVER_WG_NIC" 

SERVER_WG_IPV4="10.10.10.1"
echo "Server's VPN IPv4 : " "$SERVER_WG_IPV4"

SERVER_WG_IPV6="fd42:42:42::1"
echo "Server's VPN IPv6 : " "$SERVER_WG_IPV6"

SERVER_PORT=51820
echo "Server's VPN port : " "$SERVER_PORT"

CLIENT1_WG_IPV4="10.10.10.2"
echo "Client1's VPN IPv4 : " "$CLIENT1_WG_IPV4"

CLIENT1_WG_IPV6="fd42:42:42::2"
echo "Client1's VPN IPv6 : " "$CLIENT1_WG_IPV6"

CLIENT2_WG_IPV4="10.10.10.3"
echo "Client2's VPN IPv4 : " "$CLIENT2_WG_IPV4"

CLIENT2_WG_IPV6="fd42:42:42::3"
echo "Client2's VPN IPv6 : " "$CLIENT2_WG_IPV6"

CLIENT3_WG_IPV4="10.10.10.4"
echo "Client3's VPN IPv4 : " "$CLIENT3_WG_IPV4"

CLIENT3_WG_IPV6="fd42:42:42::4"
echo "Client3's VPN IPv6 : " "$CLIENT3_WG_IPV6"

CLIENT4_WG_IPV4="10.10.10.5"
#echo "Client4's VPN IPv4 : " "$CLIENT4_WG_IPV4"

CLIENT4_WG_IPV6="fd42:42:42::5"
#echo "Client4's VPN IPv6 : " "$CLIENT4_WG_IPV6"

CLIENT5_WG_IPV4="10.10.10.6"
#echo "Client5's VPN IPv4 : " "$CLIENT5_WG_IPV4"

CLIENT5_WG_IPV6="fd42:42:42::6"
#echo "Client5's VPN IPv6 : " "$CLIENT5_WG_IPV6"

ENDPOINT="$SERVER_PUB_IPV4:$SERVER_PORT"
echo "ENDPOINT : " "$ENDPOINT"


### 주의) 처음 실행시는 Cloud-Barista가 설치된 machine에서 아래 주석으로 처리된 부분이 실행되어야함.
# Install VPN tools and module
# if [[ "$OS" = 'ubuntu' ]]; then
#     echo "# apt update"
#     sudo apt update

#     echo "# apt -y install software-properties-common dirmngr apt-transport-https lsb-release ca-certificates"
#     sudo apt -y install software-properties-common dirmngr apt-transport-https lsb-release ca-certificates
    
#     echo "# add-apt-repository -y ppa:VPN/VPN"
#     sudo add-apt-repository -y ppa:VPN/VPN
    
#     echo "# apt-get update"
#     sudo apt-get update
    
#     echo "# apt-get install ~~~ "
#     sudo apt-get install "linux-headers-$(uname -r)"
    
#     echo "#apt-get install -y VPN iptables resolvconf"
#     sudo apt-get install -y VPN iptables resolvconf

# elif [[ "$OS" = 'debian' ]]; then
#     echo "deb http://deb.debian.org/debian/ unstable main" > /etc/apt/sources.list.d/unstable.list
#     printf 'Package: *\nPin: release a=unstable\nPin-Priority: 90\n' > /etc/apt/preferences.d/limit-unstable
#     sudo apt update
#     sudo apt-get install "linux-headers-$(uname -r)"
#     sudo apt install VPN iptables resolvconf
# elif [[ "$OS" = 'fedora' ]]; then
#     dnf copr enable jdoss/VPN
#     dnf install VPN-dkms VPN-tools iptables
# elif [[ "$OS" = 'centos' ]]; then
#     curl -Lo /etc/yum.repos.d/VPN.repo https://copr.fedorainfracloud.org/coprs/jdoss/VPN/repo/epel-7/jdoss-VPN-epel-7.repo
#     yum install epel-release
#     yum install VPN-dkms VPN-tools iptables
# elif [[ "$OS" = 'arch' ]]; then
#     pacman -S linux-headers
#     pacman -S VPN-tools iptables VPN-arch
# fi



# Make sure the directory exists (this does not seem the be the case on fedora)
sudo mkdir /etc/VPN > /dev/null 2>&1

echo "### Start to generate VPN keys!!"
umask 077

# Generate key pair for the server
SERVER_PRIV_KEY=$(wg genkey)
SERVER_PUB_KEY=$(echo "$SERVER_PRIV_KEY" | wg pubkey)

# Generate key pair for the client1
CLIENT1_PRIV_KEY=$(wg genkey)
CLIENT1_PUB_KEY=$(echo "$CLIENT1_PRIV_KEY" | wg pubkey)

# Generate key pair for the client2
CLIENT2_PRIV_KEY=$(wg genkey)
CLIENT2_PUB_KEY=$(echo "$CLIENT2_PRIV_KEY" | wg pubkey)

# Generate key pair for the client3
CLIENT3_PRIV_KEY=$(wg genkey)
CLIENT3_PUB_KEY=$(echo "$CLIENT3_PRIV_KEY" | wg pubkey)

# Generate key pair for the client4
CLIENT4_PRIV_KEY=$(wg genkey)
CLIENT4_PUB_KEY=$(echo "$CLIENT4_PRIV_KEY" | wg pubkey)

# Generate key pair for the client5
CLIENT5_PRIV_KEY=$(wg genkey)
CLIENT5_PUB_KEY=$(echo "$CLIENT5_PRIV_KEY" | wg pubkey)


umask 002

echo "### Start to create installation script files!!"
# Add server interface
echo "[Interface]
Address = $SERVER_WG_IPV4/24,$SERVER_WG_IPV6/64
ListenPort = $SERVER_PORT
PrivateKey = $SERVER_PRIV_KEY
PostUp = iptables -A FORWARD -i $SERVER_WG_NIC -j ACCEPT; iptables -t nat -A POSTROUTING -o $SERVER_PUB_NIC -j MASQUERADE; ip6tables -A FORWARD -i $SERVER_WG_NIC -j ACCEPT; ip6tables -t nat -A POSTROUTING -o $SERVER_PUB_NIC -j MASQUERADE
PostDown = iptables -D FORWARD -i $SERVER_WG_NIC -j ACCEPT; iptables -t nat -D POSTROUTING -o $SERVER_PUB_NIC -j MASQUERADE; ip6tables -D FORWARD -i $SERVER_WG_NIC -j ACCEPT; ip6tables -t nat -D POSTROUTING -o $SERVER_PUB_NIC -j MASQUERADE" >> "$SERVER_WG_NIC.conf"

# Add the client1 as a peer to the server
echo "[Peer]
PublicKey = $CLIENT1_PUB_KEY
AllowedIPs = $CLIENT1_WG_IPV4/32,$CLIENT1_WG_IPV6/128" >> "$SERVER_WG_NIC.conf"

# Add the client2 as a peer to the server
echo "[Peer]
PublicKey = $CLIENT2_PUB_KEY
AllowedIPs = $CLIENT2_WG_IPV4/32,$CLIENT2_WG_IPV6/128" >> "$SERVER_WG_NIC.conf"

# Add the client3 as a peer to the server
echo "[Peer]
PublicKey = $CLIENT3_PUB_KEY
AllowedIPs = $CLIENT3_WG_IPV4/32,$CLIENT3_WG_IPV6/128" >> "$SERVER_WG_NIC.conf"

# Add the client4 as a peer to the server
echo "[Peer]
PublicKey = $CLIENT4_PUB_KEY
AllowedIPs = $CLIENT4_WG_IPV4/32,$CLIENT4_WG_IPV6/128" >> "$SERVER_WG_NIC.conf"

# Add the client5 as a peer to the server
echo "[Peer]
PublicKey = $CLIENT5_PUB_KEY
AllowedIPs = $CLIENT5_WG_IPV4/32,$CLIENT5_WG_IPV6/128" >> "$SERVER_WG_NIC.conf"



# Create client1 file with interface
echo "[Interface]
PrivateKey = $CLIENT1_PRIV_KEY
Address = $CLIENT1_WG_IPV4/24,$CLIENT1_WG_IPV6/64" > "$SERVER_WG_NIC-client1.conf"

# Add the server as a peer to the client
echo "[Peer]
PublicKey = $SERVER_PUB_KEY
Endpoint = $ENDPOINT
AllowedIPs = $SERVER_WG_IPV4/32, $SERVER_WG_IPV6/128, $CLIENT1_WG_IPV4/32,$CLIENT1_WG_IPV6/128, $CLIENT2_WG_IPV4/32,$CLIENT2_WG_IPV6/128, $CLIENT3_WG_IPV4/32,$CLIENT3_WG_IPV6/128, $CLIENT4_WG_IPV4/32,$CLIENT4_WG_IPV6/128, $CLIENT5_WG_IPV4/32,$CLIENT5_WG_IPV6/128
PersistentKeepalive = 25" >> "$SERVER_WG_NIC-client1.conf"



# Create client2 file with interface
echo "[Interface]
PrivateKey = $CLIENT2_PRIV_KEY
Address = $CLIENT2_WG_IPV4/24,$CLIENT2_WG_IPV6/64" > "$SERVER_WG_NIC-client2.conf"

# Add the server as a peer to the client
echo "[Peer]
PublicKey = $SERVER_PUB_KEY
Endpoint = $ENDPOINT
AllowedIPs = $SERVER_WG_IPV4/32, $SERVER_WG_IPV6/128, $CLIENT1_WG_IPV4/32,$CLIENT1_WG_IPV6/128, $CLIENT2_WG_IPV4/32,$CLIENT2_WG_IPV6/128, $CLIENT3_WG_IPV4/32,$CLIENT3_WG_IPV6/128, $CLIENT4_WG_IPV4/32,$CLIENT4_WG_IPV6/128, $CLIENT5_WG_IPV4/32,$CLIENT5_WG_IPV6/128
PersistentKeepalive = 25" >> "$SERVER_WG_NIC-client2.conf"


# Create client3 file with interface
echo "[Interface]
PrivateKey = $CLIENT3_PRIV_KEY
Address = $CLIENT3_WG_IPV4/24,$CLIENT3_WG_IPV6/64" > "$SERVER_WG_NIC-client3.conf"

# Add the server as a peer to the client
echo "[Peer]
PublicKey = $SERVER_PUB_KEY
Endpoint = $ENDPOINT
AllowedIPs = $SERVER_WG_IPV4/32, $SERVER_WG_IPV6/128, $CLIENT1_WG_IPV4/32,$CLIENT1_WG_IPV6/128, $CLIENT2_WG_IPV4/32,$CLIENT2_WG_IPV6/128, $CLIENT3_WG_IPV4/32,$CLIENT3_WG_IPV6/128, $CLIENT4_WG_IPV4/32,$CLIENT4_WG_IPV6/128, $CLIENT5_WG_IPV4/32,$CLIENT5_WG_IPV6/128
PersistentKeepalive = 25" >> "$SERVER_WG_NIC-client3.conf"


# Create client4 file with interface
echo "[Interface]
PrivateKey = $CLIENT4_PRIV_KEY
Address = $CLIENT4_WG_IPV4/24,$CLIENT4_WG_IPV6/64" > "$SERVER_WG_NIC-client4.conf"

# Add the server as a peer to the client
echo "[Peer]
PublicKey = $SERVER_PUB_KEY
Endpoint = $ENDPOINT
AllowedIPs = $SERVER_WG_IPV4/32, $SERVER_WG_IPV6/128, $CLIENT1_WG_IPV4/32,$CLIENT1_WG_IPV6/128, $CLIENT2_WG_IPV4/32,$CLIENT2_WG_IPV6/128, $CLIENT3_WG_IPV4/32,$CLIENT3_WG_IPV6/128, $CLIENT4_WG_IPV4/32,$CLIENT4_WG_IPV6/128, $CLIENT5_WG_IPV4/32,$CLIENT5_WG_IPV6/128
PersistentKeepalive = 25" >> "$SERVER_WG_NIC-client4.conf"


# Create client5 file with interface
echo "[Interface]
PrivateKey = $CLIENT5_PRIV_KEY
Address = $CLIENT5_WG_IPV4/24,$CLIENT5_WG_IPV6/64" > "$SERVER_WG_NIC-client5.conf"

# Add the server as a peer to the client
echo "[Peer]
PublicKey = $SERVER_PUB_KEY
Endpoint = $ENDPOINT
AllowedIPs = $SERVER_WG_IPV4/32, $SERVER_WG_IPV6/128, $CLIENT1_WG_IPV4/32,$CLIENT1_WG_IPV6/128, $CLIENT2_WG_IPV4/32,$CLIENT2_WG_IPV6/128, $CLIENT3_WG_IPV4/32,$CLIENT3_WG_IPV6/128, $CLIENT4_WG_IPV4/32,$CLIENT4_WG_IPV6/128, $CLIENT5_WG_IPV4/32,$CLIENT5_WG_IPV6/128
PersistentKeepalive = 25" >> "$SERVER_WG_NIC-client5.conf"


echo "### Succeed to create installation script files!!"