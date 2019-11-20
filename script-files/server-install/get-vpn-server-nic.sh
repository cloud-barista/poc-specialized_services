#!/bin/bash
# To Get VPN_SERVER_PUB_NIC
ip -4 route ls | grep default | grep -Po '(?<=dev )(\S+)' | head -1