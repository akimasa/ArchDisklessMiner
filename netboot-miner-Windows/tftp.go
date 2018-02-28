package main // readHandler is called when client starts file download from server
import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/pin/tftp"
)

func readHandler(filename string, rf io.ReaderFrom) error {
	// asset, err := Asset(filename)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "%v\n", err)
	// 	return err
	// }
	// r := bytes.NewReader(asset)
	// file := r

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	n, err := rf.ReadFrom(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes sent\n", n)
	return nil
}

// writeHandler is called when client starts file upload to server
func writeHandler(filename string, wt io.WriterTo) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	n, err := wt.WriteTo(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes received\n", n)
	return nil
}
func httpHandler(w http.ResponseWriter, r *http.Request) {
	// remoteaddr := r.RemoteAddr
	intfaddr, _ := net.InterfaceAddrs()
	fmt.Println(intfaddr)
	remoteaddr := r.RemoteAddr
	remoteaddr = strings.Split(remoteaddr, ":")[0]
	myaddr := ""
	for _, v := range intfaddr {

		ipAddr, ipNet, err := net.ParseCIDR(v.String())
		if err != nil {
			fmt.Println(err)
		}

		if ipNet.Contains(net.ParseIP(remoteaddr)) {
			myaddr = ipAddr.String()
			break
		}
	}
	fmt.Println(r.RequestURI)
	if strings.Compare(r.RequestURI, "/default.ipxe") == 0 {
		responseTxt := "#!ipxe\n" +
			"kernel http://%s:80/linux quiet ip=:::::eth0:dhcp serverip=%s\n" +
			"initrd http://%s:80/initrd\n" +
			"boot"
		fmt.Fprintf(w, responseTxt, myaddr, myaddr, myaddr)
	} else if strings.Compare(r.RequestURI, "/linux") == 0 {
		http.ServeFile(w, r, "linux")
	} else if strings.Compare(r.RequestURI, "/initrd") == 0 {
		http.ServeFile(w, r, "initrd")
	} else if strings.Compare(r.RequestURI, "/arch.sfs") == 0 {
		http.ServeFile(w, r, "arch.sfs")
	}
}
func main() {
	go func() {
		cmd := exec.Command("./dhcp/dhcp.exe")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stdout, "dhcp: %v\n", err)
			os.Exit(1)
		}
	}()
	go func() {

		http.HandleFunc("/", httpHandler) // ハンドラを登録してウェブページを表示させる
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "server: %v\n", err)
			os.Exit(1)
		}
	}()
	// use nil in place of handler to disable read or write operations
	s := tftp.NewServer(readHandler, nil)
	// s.SetTimeout(5 * time.Second)  // optional
	err := s.ListenAndServe(":69") // blocks until s.Shutdown() is called
	if err != nil {
		fmt.Fprintf(os.Stdout, "server: %v\n", err)
		os.Exit(1)
	}

}
