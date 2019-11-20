// PoC-Specialized_Services based on a MCIS
// Create a MCIS-VPN environment on a MCIS consisting of multiple VMs
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by B.T. Oh, innodreamer@gmail.com, 2019.11.

package mcisvpn

import (
	"fmt"
	"bytes"
	"log"
	"os/exec"
	"strings"
	"os"
	"io/ioutil"
	"strconv"
	"encoding/json"
	"net/http"


	// CB-Store
	cbstore "github.com/cloud-barista/cb-store"
	icbs "github.com/cloud-barista/cb-store/interfaces"
	"github.com/cloud-barista/cb-store/config"
    "github.com/sirupsen/logrus"

	"github.com/cloud-barista/poc-specialized_services/vm-ssh-util"
	"github.com/cloud-barista/cb-tumblebug/src/common"	
	// "github.com/cloud-barista/cb-tumblebug/src/mcir"	

	// "github.com/hashicorp/go-multierror"

	)

type KeyValue struct {
	Key   string
	Value string
}

// CB-Store
var store icbs.Store

func init() {
	cblog = config.Cblogger
	store = cbstore.GetStore()
}

type vpnReq struct {
	Name string `json:"MCIS-VPN_Name"`
}

type mcisStatusInfo struct {
	Id     string         `json:"id"`
	Name   string         `json:"name"`
	//Vm_num string         `json:"vm_num"`
	Status string         `json:"status"`
	Vm     []vmStatusInfo `json:"vm"`
}

type vmStatusInfo struct {
	Id        string `json:"id"`
	Csp_vm_id string `json:"csp_vm_id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Public_ip string `json:"public_ip"`
	VPN_priavate_ip string
	VPN_status string
}

var cblog *logrus.Logger

func init() {
        cblog = config.Cblogger
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}


// func CreatMcisVpn(nsId string, mcisId string, vpnreq *vpnReq) (mcisStatusInfo, int, []byte, error) {
func CreatMcisVpn(nsId string, mcisId string) (mcisStatusInfo, int, []byte, error) {

	// Specialized Services Framwork Path
	SPECIAL_SERVICE_ROOT := os.Getenv("SPECIAL_SERVICE_ROOT")

	AWS_EC2_KEYPATH1 := os.Getenv("AWS_EC2_KEYPATH1")
	AWS_EC2_KEYPATH2 := os.Getenv("AWS_EC2_KEYPATH2")
	// MS_AZURE_KEYPATH := os.Getenv("MS_AZURE_KEYPATH")
	// GCE_KEYPATH := os.Getenv("GCE_KEYPATH")

	TUMBLEBUG_URL := os.Getenv("TUMBLEBUG_URL")
	
	// fmt.Println("vpnreq.Name : " + vpnreq.Name)


	// {{ip}}:{{t-port}}/ns/{{ns_id}}/mcis/{{MCIS-01}}?action=status
	url := TUMBLEBUG_URL + "/ns/" + nsId + "/mcis/" + mcisId + "?action=status"
	method := "GET"

	fmt.Println("URL : " + url)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	fmt.Println("Called Get MCIS status info API!!")
	if err != nil {
		cblog.Error(err)
		content := mcisStatusInfo{}
		// return content, res.StatusCode, nil, err
		fmt.Println("Content :", content)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	if err != nil {
		cblog.Error(err)
		content := mcisStatusInfo{}
		// return content, res.StatusCode, body, err
		fmt.Println("Content :", content)
	}

	fmt.Println("HTTP Status code " + strconv.Itoa(res.StatusCode))
	switch {
	case res.StatusCode >= 400 || res.StatusCode < 200:
		err := fmt.Errorf("HTTP Status code " + strconv.Itoa(res.StatusCode))
		fmt.Println("body: ", string(body))
		cblog.Error(err)
		content := mcisStatusInfo{}
		// return content, res.StatusCode, body, err
		fmt.Println("Content :", content)
	}

	// type mcisStatusInfo struct {
	// 	Id     string         `json:"id"`
	// 	Name   string         `json:"name"`
	// 	//Vm_num string         `json:"vm_num"`
	// 	Status string         `json:"status"`
	// 	Vm     []vmStatusInfo `json:"vm"`
	// }

	// type vmStatusInfo struct {
	// 	Id        string `json:"id"`
	// 	Csp_vm_id string `json:"csp_vm_id"`
	// 	Name      string `json:"name"`
	// 	Status    string
	// 	Public_ip string `json:"public_ip"`
	//  VPN_priavate_ip string
	//  VPN_status string
	// }

	type VPNInfo struct {
		Id      string
		Name    string
		Status  string // available, unavailable
		Vm     []vmStatusInfo
	}

	vpn := VPNInfo{}
	err2 := json.Unmarshal(body, &vpn)
	if err2 != nil {
		fmt.Println("Marshalling Error :", err2)
	}


	VM_NUM := len(vpn.Vm)

	// //VM1
	// VM_IP[0] = "15.164.49.113"
	// //VM2
	// VM_IP[1] = "52.141.19.181"
	// //VM3
	// VM_IP[2] = "34.97.45.247"


	VM_UserName := make(map[int]string)
	// VM_UserName[0] = "ubuntu"
	// VM_UserName[1] = "sean"
	// VM_UserName[2] = "sean"

	// VM_UserName := "ubuntu"


	VM_KeyPath := make(map[int]string)
	// VM_KeyPath[0] = AWS_EC2_KEYPATH
	// VM_KeyPath[1] = MS_AZURE_KEYPATH
	// VM_KeyPath[2] = GCE_KEYPATH

	// VM_KeyPath := AWS_EC2_KEYPATH


	VM_IP := make(map[int]string)
	VPN_priavate_ip := make(map[int]string)

	i := 0
	for _, value := range vpn.Vm {
		fmt.Println("VM"+ strconv.Itoa(i+1) + " Info :", value)
		VM_IP[i] = vpn.Vm[i].Public_ip	
		fmt.Println("Public_ip : ", VM_IP[i])

		// VM_IP[i] = vpn.Vm[i].userName	
		VM_UserName[i] = "ubuntu"

		
		if strings.Contains(vpn.Vm[i].Name, "aws-ap-east-1") { 
			VM_KeyPath[i] = AWS_EC2_KEYPATH1	
			fmt.Println("VM_KeyPath :", VM_KeyPath[i])
		} else if strings.Contains(vpn.Vm[i].Name, "eu-central-1") { 
			VM_KeyPath[i] = AWS_EC2_KEYPATH2	
			fmt.Println("VM_KeyPath :", VM_KeyPath[i])
		} else {
			fmt.Println("No Key files!!")
		}


		VPN_priavate_ip[i] = "10.10.10." + strconv.Itoa(i+1)
		fmt.Println("Allocated VPN_priavate_ip : ", VPN_priavate_ip[i])
		vpn.Vm[i].VPN_priavate_ip = VPN_priavate_ip[i]

		i++
	}


	// if k == 1 {
	// 	println("One")
	// } else if k == 2 {  //같은 라인
	// 	println("Two")
	// } else {   //같은 라인
	// 	println("Other")
	// }

	
	// AWS EC2
	// VM1_IP := "15.164.49.113"
	// VM1_UserName := "ubuntu"
    // VM1_KeyPath  := "/home/sean/.aws/ETRI-sean_oh_2-key.pem"

	// Azure
	// VM2_IP := "52.141.19.181"
	// VM2_UserName := "sean"
    // VM2_KeyPath  := "/home/sean/.azure/azure-vm-key"

	// GCE
	// VM3_IP := "34.97.45.247"
	// VM3_UserName := "sean"
    // VM3_KeyPath  := "/home/sean/.gcp/gce-vm-key"


	cmd001 := exec.Command("bash", "../script-files/dos2unix.sh")
	cmd001.Stdin = strings.NewReader("");
	var out1 bytes.Buffer;
	cmd001.Stdout = &out1;
	err001 := cmd001.Run();
	if err001 != nil {
			log.Fatal(err001);
	}
	fmt.Printf("Output \n",out1.String());


	//Copy the script file to the VM1 to get the nic information of the VM1
	SourceFile := SPECIAL_SERVICE_ROOT + "/script-files/server-install/get-vpn-server-nic.sh"
    TargetFile := "/home/" + VM_UserName[0] + "/get-vpn-server-nic.sh"

	err0 := sshutil.SshCopyWithKeyPath(VM_IP[0], VM_UserName[0], VM_KeyPath[0], SourceFile, TargetFile)
	if err0 != nil {
	os.Stderr.WriteString(err0.Error())
	}

	
	//Run the script to get the nic information from the VM1
	Command01 := "./get-vpn-server-nic.sh"

	result01, err01 := sshutil.SshRunWithKeyPath(VM_IP[0], VM_UserName[0], VM_KeyPath[0], Command01)
	if err01 != nil {
	os.Stderr.WriteString(err01.Error())
	}
	fmt.Println(result01)

	VPN_SERVER_PUB_NIC := result01
	fmt.Println("VPN_SERVER_PUB_NIC : " + VPN_SERVER_PUB_NIC)	


	//Write an environment file about VPN_SERVER_PUB_IP and VPN_SERVER_PUB_NIC for bash script 
    f, err := os.Create("init.env")
    checkError(err)
    
    n, err := f.WriteString("export VPN_SERVER_PUB_IP=" + VM_IP[0] + "\n" + "export VPN_SERVER_PUB_NIC=" + VPN_SERVER_PUB_NIC + "\n")
    checkError(err)
    fmt.Printf("File Wrote %d bytes\n", n)
    
    f.Sync()


	//Create the VPN ~.conf files to install VPN on each VM
	cmd002 := exec.Command("bash", "../script-files/create-client-server_scripts_v1.5.sh")
	cmd002.Stdin = strings.NewReader("");
	var out2 bytes.Buffer;
	cmd002.Stdout = &out2;
	err002 := cmd002.Run();
	if err002 != nil {
		log.Fatal(err002);
	}
	fmt.Printf("Output \n",out2.String());


	//Copyt the VPN ~.conf file to each VM
	Source := make(map[int]string)
	Target := make(map[int]string)
	copy_err := make(map[int]error)

	Source[0] = SPECIAL_SERVICE_ROOT + "/rest-runtime/barista0.conf"
	// Do not use the path like $HOME/~~~ 

	for i := 0; i < VM_NUM ; i++ {
		if i != 0 {
		Source[i] = SPECIAL_SERVICE_ROOT + "/rest-runtime/barista0-client" + strconv.Itoa(i) + ".conf"
		// fmt.Println("Source Num : " + strconv.Itoa(i));
		// fmt.Println("Source : " + Source[i]);
		}
		Target[i] = "/home/" + VM_UserName[i] + "/barista0.conf"

		copy_err[i] = sshutil.SshCopyWithKeyPath(VM_IP[i], VM_UserName[i], VM_KeyPath[i], Source[i], Target[i])
		if copy_err[i] != nil {
		os.Stderr.WriteString(copy_err[i].Error())
		}
	}


	//Copyt each script file to each VM to install VPN application according to the ~.conf file
	ServerInstallScript := SPECIAL_SERVICE_ROOT + "/script-files/server-install/server_install_v1.4.sh"
	ClientInstallScript := SPECIAL_SERVICE_ROOT + "/script-files/client-install/client_install_v1.4.sh"
	// Do not use the path like $HOME/~~/~ 

	VPN_install := make(map[int]string)
	install_err := make(map[int]error)
	
	VPN_install[0] = "/home/" + VM_UserName[0] + "/server_install_v1.4.sh"

	install_err[0] = sshutil.SshCopyWithKeyPath(VM_IP[0], VM_UserName[0], VM_KeyPath[0], ServerInstallScript, VPN_install[0])
	if install_err[0] != nil {
	os.Stderr.WriteString(install_err[0].Error())
	}


	for i := 0; i < VM_NUM ; i++ {
		if i != 0 {
		VPN_install[i] = "/home/" + VM_UserName[i] + "/client_install_v1.4.sh"

		install_err[i] = sshutil.SshCopyWithKeyPath(VM_IP[i], VM_UserName[i], VM_KeyPath[i], ClientInstallScript, VPN_install[i])
		if install_err[i] != nil {
		os.Stderr.WriteString(install_err[i].Error())
		}
		}
	}
	


	// ### Run each of script file to install VPN application according to the ~.conf file on each VM ###	
	// var err_result error

	run_result := make(map[int]string)
	run_command := make(map[int]string)
	run_err := make(map[int]error)

	run_command[0] = "sudo ./server_install_v1.4.sh"

	run_result[0], run_err[0] = sshutil.SshRunWithKeyPath(VM_IP[0], VM_UserName[0], VM_KeyPath[0], run_command[0])
	if run_err[0] != nil {
	os.Stderr.WriteString(run_err[0].Error())
	}
	fmt.Println(run_result[0])

	// err_result = multierror.Append(err_result, run_err[0])


	for i := 0; i < VM_NUM ; i++ {
		if i != 0 {
		run_command[i] = "sudo ./client_install_v1.4.sh"

		run_result[i], run_err[i] = sshutil.SshRunWithKeyPath(VM_IP[i], VM_UserName[i], VM_KeyPath[i], run_command[i])
		if run_err[i] != nil {
		os.Stderr.WriteString(run_err[i].Error())

		// err_result = multierror.Append(err_result, run_err[i])
			}
		}
	}





	// Command05 := "ifconfig | grep destination"

	// result05, err05 := sshutil.SshRunWithKeyPath(VM1_IP, VM1_UserName, VM1_KeyPath, Command05)
	// if err05 != nil {
	// os.Stderr.WriteString(err05.Error())
	// }
	// fmt.Println(result05)


	// type vpnInfo struct {
	// 	Id      string
	// 	Name    string
	// 	GuestOS string // Windows7, Ubuntu etc.
	// 	Status  string // available, unavailable

	// 	KeyValueList []common.KeyValue
	// }


	//Run the script to get the VPN status information from the VM1
	StatusCommand := "ifconfig | grep destination"

	VPN_Status := make(map[int]string)
	run_err2 := make(map[int]error)
	VPN_RuningIP := make(map[int]string)


    for k := 0; k < len(vpn.Vm); k++ {
		VPN_Status[k], run_err2[k] = sshutil.SshRunWithKeyPath(VM_IP[k], VM_UserName[k], VM_KeyPath[k], StatusCommand)
		if run_err2[k] != nil {
		os.Stderr.WriteString(run_err2[k].Error())
		}
		fmt.Println("VPN_Status :"+ VPN_Status[k])

		VPN_RuningIP[k] = "10.10.10." + strconv.Itoa(k+1)

		if strings.Contains(VPN_Status[k], VPN_RuningIP[k]) { 
		fmt.Println(vpn.Vm[k].VPN_status) 
		vpn.Vm[k].VPN_status = "VPN is Runing"
		} else {
		vpn.Vm[k].VPN_status = "VPN is Not Runing"
		}
	}



	vpncontent := mcisStatusInfo{}
	vpncontent.Id = common.GenUuid()
	vpncontent.Vm = vpn.Vm
	// content.Name = vpnreq.Name
	// vpncontent.Status = vpn.Status
	vpncontent.Status = "Ref. VPN info of each VM!!"


	// Store the metadata to CB-Store.
	// fmt.Println("============= PUT VPNInfo =============")
	// Key := mcir.genResourceKey(nsId, "VPN", vpncontent.Id)
	// Val, _ := json.Marshal(vpncontent)
	// cbStorePutErr := store.Put(string(Key), string(Val))
	// if cbStorePutErr != nil {
	// 	cblog.Error(cbStorePutErr)
	// 	return vpncontent, res.StatusCode, body, cbStorePutErr
	// }
	// keyValue, _ := store.Get(string(Key))
	// fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
	// fmt.Println("=======================================")
	// return vpncontent, res.StatusCode, body, nil



	// fmt.Println("Content :", vpncontent)
	fmt.Println("===========================")
	return vpncontent, res.StatusCode, body, install_err[0]

	// $$$$$
	// return "Success!!", err_result
	// return &result1, err01
}