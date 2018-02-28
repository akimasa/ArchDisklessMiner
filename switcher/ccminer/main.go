package ccminer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

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

func init() {
	// miner_name -> nicehash
	dict = map[string]string{
		"scrypt":     "scrypt",
		"sha256":     "sha256",
		"scryptnf":   "scryptnf",
		"x11":        "x11",
		"x13":        "x13",
		"keccak":     "keccak",
		"x15":        "x15",
		"nist5":      "nist5",
		"neoscrypt":  "neoscrypt",
		"lyra2":      "lyra2le",
		"whirlpoolx": "whirlpoolx",
		"qubit":      "qubit",
		"quark":      "quark",
		// "axiom":           "axiom",
		"lyra2v2":     "lyra2rev2",
		"scrypt-jane": "scryptjanenf16",
		"blakecoin":   "blake256r8",
		// "blake256r14":     "blake256r14",
		// #"vanilla": "blake256r8vnl",
		// "hodl":            "hodl",
		// "daggerhashimoto": "daggerhashimoto",
		"decred":      "decred",
		"cryptonight": "cryptonight",
		"lbry":        "lbry",
		// #"equihash":    "equihash",
		"pascal": "pascal",
		"sib":    "x11gost",
		// #"sia":         "sia",
		// #"blake2s":     "blake2s",
		"skunk": "skunk",
	}
}

func (arr *ccHashrateArr) scanCcminer(scanner *bufio.Scanner) {
	startflag := false
	re := regexp.MustCompile(`] +(?P<name>.*?) : *(?P<rate>\d*.*?) kH/s,`)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "Benchmark results for GPU") {
			startflag = true
		}
		if startflag {
			match := re.FindStringSubmatch(text)
			if match != nil {
				// fmt.Println(match[1], match[2])
				hr, err := strconv.ParseFloat(match[2], 64)
				if err != nil {
					continue
				}
				a := make(ccHashrateArr, 1)
				a[0] = ccHashrate{MinerName: match[1], Hashrate: hr}
				*arr = append(*arr, a[0])

			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
func (m CcHashrateMap) ScanCcminer(scanner *bufio.Scanner) {
	startflag := false
	re := regexp.MustCompile(`] +(?P<name>.*?) : *(?P<rate>\d*.*?) kH/s,`)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)
		if strings.Contains(text, "Benchmark results for GPU") {
			startflag = true
		}
		if startflag {
			match := re.FindStringSubmatch(text)
			if match != nil {
				// fmt.Println(match[1], match[2])
				name := match[1]
				hr, err := strconv.ParseFloat(match[2], 64)
				if err != nil {
					continue
				}
				if _, ok := m[name]; !ok {
					m[name] = hr
				} else {
					m[name] += hr
				}

			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func BenchString() (string, error) {
	hrArr := BenchHashrate()
	yamlbyte, err := yaml.Marshal(hrArr)
	if err != nil {
		return "", err
	}
	return string(yamlbyte), nil
}
func BenchHashrate() []Hashrate {
	cmd := exec.Command("bash", "-c", "LD_LIBRARY_PATH=./ ./ccminer --benchmark")
	cmdStdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdStdout)
	cmd.Start()

	m := CcHashrateMap{}
	m.ScanCcminer(scanner)
	cmd.Wait()

	var hrArr []Hashrate
	for k, v := range m {
		hr := Hashrate{NicehashName: k, MinerName: k, Hashrate: v, Miner: "ccminer"}

		if v, ok := dict[k]; ok {
			hr.NicehashName = v
			hrArr = append(hrArr, hr)
		}

	}
	return hrArr

}
