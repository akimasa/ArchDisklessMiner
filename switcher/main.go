package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"./ccminer"
	"./excavator"
	"gopkg.in/yaml.v2"
)

var userpass string
var gpunum int
var benchmode bool
var location string

type hashratearr []ccminer.Hashrate

func (hr *hashratearr) getHashrate(NicehashName string) (Hashrate float64) {
	for _, element := range *hr {
		if strings.Compare(element.NicehashName, NicehashName) == 0 {
			return element.Hashrate
		}
	}
	return 0
}
func (hr *hashratearr) getMiner(NicehashName string) (Hashrate string) {
	for _, element := range *hr {
		if strings.Compare(element.NicehashName, NicehashName) == 0 {
			return element.Miner
		}
	}
	return ""
}
func (hr *hashratearr) getMinerName(NicehashName string) (Hashrate string) {
	for _, element := range *hr {
		if strings.Compare(element.NicehashName, NicehashName) == 0 {
			return element.MinerName
		}
	}
	return ""
}

type maxearn struct {
	Earning      float64
	NicehashName string
	Name         string
	Miner        string
	Port         int
}

func getMaxEarn(hr hashratearr) (me maxearn, err error) {
	me.Earning = 0
	resp, err := http.Get("https://api.nicehash.com/api?method=simplemultialgo.info")
	if err != nil {
		return me, err
	}
	dec := json.NewDecoder(resp.Body)
	type Message struct {
		Name, Text string
	}
	type Simplemultialgo struct {
		Paying string `json:"paying"`
		Name   string `json:"name"`
		Port   int    `json:"port"`
		Algo   int    `json:"algo"`
	}
	type Result struct {
		Simplemultialgo []Simplemultialgo `json:"simplemultialgo"`
	}
	type Body struct {
		Method string `json:"method"`
		Result Result `json:"result"`
	}
	for {
		var m Body
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		for _, s := range m.Result.Simplemultialgo {
			// fmt.Println(i, s)
			hashrate := hr.getHashrate(s.Name)
			paying, err := strconv.ParseFloat(s.Paying, 64)
			if err != nil {
				continue
			}
			earning := hashrate * paying
			if earning == 0 {
				continue
			}
			if earning > me.Earning {
				me.Earning = earning
				me.NicehashName = s.Name
				me.Name = hr.getMinerName(s.Name)
				me.Miner = hr.getMiner(s.Name)
				me.Port = s.Port
			}
		}

	}
	return me, nil
}
func bench() []ccminer.Hashrate {
	ccminerHr := ccminer.BenchHashrate()

	eth := excavator.BenchmarkAlgo("daggerhashimoto", 30)
	ethHr := ccminer.Hashrate{NicehashName: "daggerhashimoto", MinerName: "daggerhashimoto", Hashrate: float64(eth / 1000), Miner: "excavator"}
	ccminerHr = append(ccminerHr, ethHr)

	equihash := excavator.BenchmarkAlgo("equihash", 10)
	equihashHr := ccminer.Hashrate{NicehashName: "equihash", MinerName: "equihash", Hashrate: float64(equihash / 1000), Miner: "excavator"}

	ccminerHr = append(ccminerHr, equihashHr)

	yamlbyte, err := yaml.Marshal(ccminerHr)
	if err != nil {
		panic(err)
	}
	fmt.Println(`saving hashrate.yaml`)
	ymlfile, err := os.OpenFile("hashrate.yaml", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(`error opening hashrate.yaml!`)
		return ccminerHr
	}
	ymlfile.Write(yamlbyte)
	ymlfile.Close()

	return ccminerHr
}
func init() {
	flag.IntVar(&gpunum, "gpunum", 1, "Number of your gpus")
	flag.StringVar(&userpass, "wallet", "", "Your BTC wallet")
	flag.StringVar(&location, "location", "jp", "Nicehash server locaiton")
	flag.BoolVar(&benchmode, "benchmark", false, "Benchmark mode")
}
func main() {
	flag.Parse()
	if benchmode {
		hrs := bench()
		yamlbyte, err := yaml.Marshal(hrs)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(yamlbyte))
		os.Exit(0)
	}
	if strings.Compare("", userpass) == 0 {
		panic("no wallet!")
	}
	userpass += ":x"

	var hr hashratearr

	yamlfile, err := ioutil.ReadFile("hashrate.yaml")
	if os.IsNotExist(err) {
		fmt.Println(`hashrate.yml does't exist.`)
		fmt.Println(`starting benchmark...`)
		hr = bench()
	} else {
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(yamlfile, &hr)
		if err != nil {
			panic(err)
		}
	}
	firsttime := true
	lastcmd := ""
	var execCmd *exec.Cmd
	for {
		me, err := getMaxEarn(hr)
		if err != nil {
			continue
		}
		fmt.Printf("Select name:%s, earning:%f\n", me.NicehashName, me.Earning)
		cmd := ""
		if strings.Compare("ccminer", me.Miner) == 0 {
			url := fmt.Sprintf("stratum+tcp://%s.%s.nicehash.com:%d", me.NicehashName, location, me.Port)
			cmd = fmt.Sprintf("LD_LIBRARY_PATH=./ ./ccminer -o %s -a %s -O %s", url, me.Name, userpass)
		} else if strings.Compare("excavator", me.Miner) == 0 {
			excavator.MakeJSON(me.NicehashName, me.Port, userpass, location)
			cmd = fmt.Sprintf("ALGO=%s ./excavator -c excavator-run.json", me.NicehashName)
		} else {
			tmp := me.Miner
			hostport := fmt.Sprintf("%s.%s.nicehash.com:%d", me.NicehashName, location, me.Port)
			tmp = strings.Replace(tmp, "%HOSTPORT%", hostport, -1)
			tmp = strings.Replace(tmp, "%USERPASS%", userpass, -1)
			cmd = tmp
		}
		if strings.Compare(cmd, lastcmd) != 0 {
			if firsttime {
				firsttime = false
			} else {
				execCmd.Process.Kill()
				execCmd.Process.Wait()
			}
			execCmd = exec.Command("bash", "-c", cmd)
			execCmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGTERM}
			execCmd.Stderr = os.Stderr
			execCmd.Stdout = os.Stdout
			execCmd.Start()
			fmt.Println(execCmd.Process.Pid)
			fmt.Println("switch cmd", cmd)
		}
		lastcmd = cmd
		time.Sleep(60 * time.Second)
	}
}
