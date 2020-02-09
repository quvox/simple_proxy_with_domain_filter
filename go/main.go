/**
 *  [2020] Takeshi Kubo
 *  All Rights Reserved.
 */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/go-httpproxy/httpproxy"
	"github.com/pkg/profile"
	"log"
	"main/env"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	VERSION = "0.2.5"
)


var (
	port  = flag.Int("port", 8080, "proxy port")
	blockFile = flag.String("block", "", "filename of block list")
	showVersion = flag.Bool("v", false, "show version")

	blockList = map[string]bool{}
)


func readBlockList(filePath string) {
	blockList = map[string]bool{}

	fp, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil { return }
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(line) == 0 { continue }
		//log.Println(line)
		blockList[line] = true
	}
}


func inBlockList(host string) bool {
	fqdn := strings.Split(host, ":")[0]
	_, ok := blockList[fqdn]
	if ok { return true }
	lvs := strings.Split(fqdn, ".")
	for i := 1; i < len(lvs)-1; i++ {
		h := strings.Join(lvs[i:], ".")
		_, ok = blockList[h]
		if ok { return true }
	}
	return false
}


func hijackConnection(wrt http.ResponseWriter) net.Conn {
	hj, ok := wrt.(http.Hijacker)
	if !ok {
		panic("XXXX httpserver does not support hijacking")
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		panic("XXXX Cannot hijack connection " + err.Error())
	}

	return conn
}


func OnAccept(ctx *httpproxy.Context, w http.ResponseWriter, r *http.Request) bool {
	//log.Printf("[%s] receive HTTP request\n", r.Host)
	if inBlockList(r.Host) {
		log.Printf("[%s] BLOCK!!\n", r.Host)
		hijackConnection(w).Close()
		return true
	}
	return false
}

func OnConnect(ctx *httpproxy.Context, host string) (
	ConnectAction httpproxy.ConnectAction, newHost string) {
	// Apply "Man in the Middle" to all ssl connections. Never change host.
	return httpproxy.ConnectProxy, host
}

func OnRequest(ctx *httpproxy.Context, req *http.Request) (
	resp *http.Response) {
	// Log proxying requests.
	log.Printf("INFO: Proxy: %s %s", req.Method, req.URL.String())
	return
}

func OnResponse(ctx *httpproxy.Context, req *http.Request,
	resp *http.Response) {
	// Add header "Via: go-httpproxy".
	resp.Header.Add("Via", "go-httpproxy")
}


func main() {
	if env.DEBUG { defer profile.Start(profile.ProfilePath("./")).Stop() }

	//defer profile.Start(profile.MemProfile).Stop()
	flag.Parse()
	if *showVersion {
		fmt.Printf("version: %s\n", VERSION)
		return
	}

	if *blockFile != "" { readBlockList(*blockFile) }

	prx, _ := httpproxy.NewProxy()
	// Set handlers.
	prx.OnAccept = OnAccept
	prx.OnConnect = OnConnect
	prx.OnRequest = OnRequest
	prx.OnResponse = OnResponse

	// Listen...
	http.ListenAndServe(fmt.Sprintf(":%d", *port), prx)
}
