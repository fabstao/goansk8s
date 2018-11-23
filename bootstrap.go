/*
* (C) 2018 Intel EOS OTC
* fabianx.salamanca.dominguez@intel.com
 */

package main

import (
	"fmt"
	"os"
	"text/template"
	"strconv"
)

var templ *template.Template

type VM struct {
	Name    string
	Flavor  string
	Network string
	Key     string
	SG      string
	Image   string
	State   string
	Debug   bool
}

// Just a pretty (?) way to print CLI messages
func banner(m1 string, m2 string) {
	fmt.Println("=============================================")
	fmt.Println("")
	fmt.Println(" ", m1)
	fmt.Println(" ", m2)
	fmt.Println("")
	fmt.Println("=============================================")
}

func init() {
	templ = template.Must(template.ParseFiles("deploy.goyaml"))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	args := os.Args
	if args == nil || len(args) < 5 {
		fmt.Printf("ERROR: Usage %v <present|absent> <image> <cluster-name> <qty of nodes>", args[0])
		panic("")
	}
	state := args[1]
	image := args[2]
	debug := false
	if state == "present" {
		debug = true
	}
	cluster := args[3]
	nodes, err := strconv.Atoi(args[4])
	checkError(err)

	banner("Output ansible Playbook for automatic ", "K8s cluster creation on Openstack")
	file, err := os.Create("goadeploy.yaml")
	checkError(err)
	var k8scluster []VM

	//Create VMs for k8s cluster
	for iter := 0; iter < nodes; iter++ {
		k8scluster = append(k8scluster, VM{
			Name:    cluster + strconv.Itoa(iter),
			Flavor:  "10",
			Network: "provider",
			Key:     "EOSKEY",
			SG:      "Linux-Generic",
			Image:   image,
			State:   state,
			Debug:   debug,
		})
	}
	err = templ.ExecuteTemplate(file, "deploy.goyaml", k8scluster)
	checkError(err)
	banner("Playbook created successfully you now may run:", "ansible-playbook ./goadeploy.yaml")
}
