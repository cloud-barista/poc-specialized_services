# poc-specialized_services
Cloud-Barista 'Specialized Services' based on Multi-Cloud Infra Service(MCIS)

* 현재는 멀티 클라우드 기반 동일 subnet을 제공하는 MC-VPN 기능 코드만 push된 상태임.

[ # 구동 전 참고 사항 # ]
1. CB-Store(https://github.com/cloud-barista/cb-store), CB-Log(https://github.com/cloud-barista/cb-log), CB-Spider(https://github.com/cloud-barista/cb-spider), CB-Tumblebug(https://github.com/cloud-barista/cb-tumblebug)가 설치 및 구동된 상태에서 실행해야함.

2. 처음 실행시는 Cloud-Barista가 설치된 machine에 wireguard가 설치되어야함.
   ./script-files/create-client-server_scripts_v1.5.sh 93라인 참고해서 그부분 주석을 없애면됨.

[ 구동 순서)
1. source ./setup.evn 실행
2. cd /rest-runtime/
3. ./build.sh를 실행하거나 go run *.go 실행

* 문의처 : innodreamer@gmail.com