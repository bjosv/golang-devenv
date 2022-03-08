package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
)

const cpu_limit int = 10

func getInt(str string) int {
	parts := strings.Split(str, ",")
	if len(parts) > 1 {
		i, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatal(err)
		}
		return i
	}
	return 0
}

func main() {

	var pid string
	var port string
	flag.StringVar(&pid, "pid", "0", "pid to stat")
	flag.StringVar(&port, "port", "6060", "pprof port")
	url := fmt.Sprintf("http://localhost:%s/debug/pprof/profile", port)

	cmd := exec.Command("pidstat", "-p", pid, "1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(stdout)

	cpupeak := false
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		s := strings.Fields(str)
		//fmt.Printf("> %s\n", str)
		if len(s) > 10 {
			fmt.Printf("> usr:%d system:%d cpu:%d\n", getInt(s[4]), getInt(s[5]), getInt(s[8]))
			if getInt(s[8]) >= cpu_limit {
				if cpupeak == false {
					cpupeak = true
					resultfile := fmt.Sprintf("profile-%d.out", rand.Int())
					sh := fmt.Sprintf("curl -sK -v %s > %s", url, resultfile)
					fmt.Printf("Peaking, call cmd=%s\n", sh)

					cmd := exec.Command("sh", "-c", sh)
					stdout, err := cmd.Output()
					if err != nil {
						fmt.Println(err.Error())
					}
					fmt.Print(string(stdout))
				}
				fmt.Printf("> GOT cpu:%d\n", getInt(s[8]))
			} else {
				if cpupeak == true {
					fmt.Printf("Not peaking\n")
				}
				cpupeak = false
			}
		}
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
