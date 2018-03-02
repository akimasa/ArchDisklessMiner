package ethminer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Hashrate struct {
	NicehashName string  `yaml:"nicehash_name"`
	MinerName    string  `yaml:"miner_name"`
	Hashrate     float64 `yaml:"hashrate"`
	Miner        string  `yaml:"miner"`
}

type ccHashrate struct {
	MinerName string
	Hashrate  float64
}
type ccHashrateArr []ccHashrate

type CcHashrateMap map[string]float64

var dict map[string]string

func scanEthminer(scanner *bufio.Scanner) chan string {

	re := regexp.MustCompile(`speed: (?P<rate>\d*.\d*?) H/s`)
	receiver := make(chan string)
	go func() {
		hr := ""
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Println(text)
			match := re.FindStringSubmatch(text)
			if match != nil {
				fmt.Println(match[1])
				hr = match[1]
				receiver <- hr
				return
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

	}()
	return receiver
}

func BenchString() (string, error) {
	hrArr := BenchHashrate()
	yamlbyte, err := yaml.Marshal(hrArr)
	if err != nil {
		return "", err
	}
	return string(yamlbyte), nil
}
func BenchHashrate() float64 {
	cmd := exec.Command("bash", "-c", "./bin/ethminer -Z- U")
	cmdStdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdStdout)
	cmd.Start()

	receiver := scanEthminer(scanner)
	for {
		select {
		case recv := <-receiver:
			// hogehoge
			fmt.Println(recv)
		case <-time.After(time.Second * 6):
			// timeout!
			cmd.Process.Kill()
			cmd.Process.Wait()
			return 0
		}
	}
	cmd.Wait()

	return 1

}
