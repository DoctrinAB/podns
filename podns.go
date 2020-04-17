package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"text/tabwriter"
)

type pod struct {
	Name, NodeIP, ContainerID, Pid string
}

type workerFunc func(*pod, *sync.WaitGroup, chan error)

var remoteUser = flag.String("u", "centos", "remote user")

func main() {
	flag.Parse()
	if len(os.Args) < 3 {
		flag.Usage()
		return
	}
	log.SetFlags(0) // stderr
	tabw := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
	remoteCmd := os.Args[1]
	pods := createPods(os.Args[2:])
	log.Println("* get node_ip and container_id")
	errList := runWorkers(ipAndCid, pods)
	if errList != nil {
		log.Fatal(errList)
	}
	printIpAndCid(tabw, pods)
	log.Println("* get pid")
	errList = runWorkers(getPid(*remoteUser), pods)
	if errList != nil {
		log.Fatal(errList)
	}
	printPid(tabw, pods)
	log.Println(fmt.Sprintf("* running %s in pid namespace", remoteCmd))
	errList = runWorkers(remoteCommand(*remoteUser, remoteCmd), pods)
	if errList != nil {
		log.Fatal(errList)
	}
	log.Println("* done")
}

func createPods(podnames []string) (pods []*pod) {
	for _, podname := range podnames {
		pods = append(pods, &pod{
			Name: podname,
		})
	}
	return
}

func printIpAndCid(w *tabwriter.Writer, pods []*pod) {
	for _, p := range pods {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s", p.Name, p.NodeIP, p.ContainerID))
	}
	w.Flush()
}

func printPid(w *tabwriter.Writer, pods []*pod) {
	for _, p := range pods {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s", p.Name, p.Pid))
	}
	w.Flush()
}

func runWorkers(fn workerFunc, pods []*pod) (errList []error) {
	errc := make(chan error)
	go func(fn workerFunc, pods []*pod) {
		var wg sync.WaitGroup
		wg.Add(len(pods))
		for _, p := range pods {
			go fn(p, &wg, errc)
		}
		wg.Wait()
		close(errc)
	}(fn, pods)
	for err := range errc {
		errList = append(errList, err)
	}
	return
}

func ipAndCid(p *pod, wg *sync.WaitGroup, errc chan error) {
	defer wg.Done()
	jsonpath := "-o=jsonpath='{.status.hostIP},{.status.containerStatuses[0].containerID}'"
	out, err := exec.Command("oc", "get", "pod", p.Name, jsonpath).Output()
	if err != nil {
		errc <- err
		return
	}
	parts := strings.Split(string(out), ",")
	p.NodeIP = strings.ReplaceAll(parts[0], "'", "")
	p.ContainerID = strings.ReplaceAll(strings.TrimPrefix(parts[1], "docker://"), "'", "")
}

func getPid(remoteUser string) workerFunc {
	return func(p *pod, wg *sync.WaitGroup, errc chan error) {
		defer wg.Done()
		host := fmt.Sprintf("%s@%s", remoteUser, p.NodeIP)
		out, err := exec.Command("ssh", host, "sudo", "docker", "inspect", "-f", `'{{.State.Pid}}'`, p.ContainerID).Output()
		if err != nil {
			errc <- err
			return
		}
		p.Pid = strings.TrimSuffix(string(out), "\n")
	}
}

func remoteCommand(remoteUser, remoteCmd string) workerFunc {
	return func(p *pod, wg *sync.WaitGroup, errc chan error) {
		defer wg.Done()
		host := fmt.Sprintf("%s@%s", remoteUser, p.NodeIP)
		cmd := exec.Command("ssh", host, "sudo", "nsenter", "-t", p.Pid, "-n", remoteCmd)
		cmd.Stdout = os.Stdout
		err := cmd.Start()
		if err != nil {
			errc <- err
			return
		}
		err = cmd.Wait()
		if err != nil {
			errc <- err
			return
		}
	}
}
