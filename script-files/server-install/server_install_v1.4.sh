#!/bin/bash
echo "Start to install VPN server!!"

if [ "$(systemd-detect-virt)" == "openvz" ]; then
    echo "OpenVZ is not supported"
    exit
fi

if [ "$(systemd-detect-virt)" == "lxc" ]; then
    echo "LXC is not supported (yet)."
    echo "WireGuard can technically run in an LXC container,"
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

# Install WireGuard tools and module
if [[ "$OS" = 'ubuntu' ]]; then

    sudo apt -y install ufw
    sudo ufw default deny incoming
    sudo ufw default allow outgoing
    sudo ufw allow ssh
    sudo ufw allow 2379/tcp 
    sudo ufw allow 2380/tcp
    sudo ufw allow 51820/udp   
    echo "y" | sudo ufw enable
    sudo ufw status
    # VM에서 기존 열린 port 확인해서 열어야함.

    echo "# apt update"
    sudo apt update

    echo "# apt -y install software-properties-common dirmngr apt-transport-https lsb-release ca-certificates"
    sudo apt -y -q install software-properties-common dirmngr apt-transport-https lsb-release ca-certificates
    
    echo "# add-apt-repository -y ppa:wireguard/wireguard"
    sudo add-apt-repository -y ppa:wireguard/wireguard
    
    echo "# apt-get update"
    sudo apt-get update
    
    echo "# apt-get install ~~~ "
    sudo apt-get install -y -q "linux-headers-$(uname -r)"
    
    echo "#apt-get install -y wireguard iptables resolvconf"
    sudo apt-get install -y -q wireguard iptables resolvconf

elif [[ "$OS" = 'debian' ]]; then
    echo "deb http://deb.debian.org/debian/ unstable main" > /etc/apt/sources.list.d/unstable.list
    printf 'Package: *\nPin: release a=unstable\nPin-Priority: 90\n' > /etc/apt/preferences.d/limit-unstable
    sudo apt update
    sudo apt-get install "linux-headers-$(uname -r)"
    sudo apt install wireguard iptables resolvconf

elif [[ "$OS" = 'fedora' ]]; then
    dnf copr enable jdoss/wireguard
    dnf install wireguard-dkms wireguard-tools iptables
elif [[ "$OS" = 'centos' ]]; then
    curl -Lo /etc/yum.repos.d/wireguard.repo https://copr.fedorainfracloud.org/coprs/jdoss/wireguard/repo/epel-7/jdoss-wireguard-epel-7.repo
    yum install epel-release
    yum install wireguard-dkms wireguard-tools iptables
elif [[ "$OS" = 'arch' ]]; then
    pacman -S linux-headers
    pacman -S wireguard-tools iptables wireguard-arch
fi

# Enable routing on the server
sudo echo "net.ipv4.ip_forward = 1
net.ipv6.conf.all.forwarding = 1" > /etc/sysctl.d/wg.conf

sudo sysctl --system

# Make sure the directory exists (this does not seem the be the case on fedora)
sudo mkdir /etc/wireguard > /dev/null 2>&1

sudo cp -f ./barista0.conf /etc/wireguard/

chmod 600 -R /etc/wireguard/

SERVER_WG_NIC="barista0"
echo "### VPN interface name : " "$SERVER_WG_NIC" 

# Start the Wireguard service
sudo wg-quick up $SERVER_WG_NIC

# Enable the Wireguard service to automatically restart on boot
sudo systemctl enable "wg-quick@$SERVER_WG_NIC"

echo "### Installation Finished!!(VPN Server)"

#sudo wg show

echo "### Running VPN IP : " 
ifconfig | grep 10.10.10

echo "######################################################################"
