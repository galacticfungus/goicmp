package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"syscall"
)

type restapi struct {
    // System file descriptor.
    icmpChanCount uint64
	fs http.FileSystem
}
// Payload is the JSON Type used for reporting the number of ICMP messages recieved
type Payload struct {
	Count uint64 `json:"count"`
}

// InitServerEndPoints Why is this function exported
func initServer() *restapi {
	api := new(restapi)
	api.icmpChanCount = 0
	api.fs = http.Dir("./build")
	println("Starting server on localhost:8080")
	http.HandleFunc("/data", api.dataRequest)
	http.HandleFunc("/health_check", api.healthCheck)
	// Basically we intercept file requests and decide what can actually be served
	http.Handle("/", http.FileServer(api))
	return api
}

func (restApi *restapi) healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Ready</h1>")
}

func (restApi *restapi) dataRequest(w http.ResponseWriter, r *http.Request) {
	p := Payload {
		Count: atomic.LoadUint64(&restApi.icmpChanCount),
	}
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error encoding JSON: ", err.Error())
		http.Error(w, err.Error(), 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		fmt.Println("Failed to write JSON: ", err.Error())
		http.Error(w, err.Error(), 500)
	}
	
}

func main() {
	api := initServer()
	go api.interceptIcmp()

	
	
	http.ListenAndServe(":8080", nil)

	
}

func (restApi *restapi) interceptIcmp() {
	// This is currently counting both sending and recieving
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	//socket := netSocket{fd: fd}
	if err != nil {
        fmt.Println("Failed to create the raw socket: Error was", err)
		return
    }
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	for {
		buf := make([]byte, 1024) // Blocks until it recieves data
		_, err := f.Read(buf) // numRead
		if err != nil {
			fmt.Println("Error occured while reading the buffer: ",err)
		}
		// fmt.Printf("Read %d bytes", numRead)
		// fmt.Printf("Buffer is: %s", buf)
		header := buf[0]
		ipSize := 0x0f & header // We are interested in the size of the ip header so we can skip it
		// ipSize is the length of the header is DWORDS or the number of 4 bytes
		// We don't need to check IP protocol type since we only see ICMP packets
		// in this socket
		ipSize *= 4
		icmpHeader := buf[ipSize:ipSize+8]
		if icmpHeader[0] == 8 { // 0 is reply 8 is request, we don't care about our responses to ICMP requests
			atomic.AddUint64(&restApi.icmpChanCount, 1)
		}
	}
}

type checkedFileSystem struct {
    fs http.FileSystem
}

// This is used to return 404 for any directory entries on the server or any files we dont want served directly
func (restApi restapi) Open(path string) (http.File, error) {
    f, err := restApi.fs.Open(path)
    if err != nil {
        return nil, err
    }
	fmt.Println("Request for ", path)
    s, err := f.Stat()
    if s.IsDir() {
		// If the directory has an index.html then we can serve that file instead
		// All paths will be relative to /build since that is where we server static files from
		// ie paths not relative to /build will not trigger the file system
        index := filepath.Join(path, "index.html")
        if _, err := restApi.fs.Open(index); err != nil {
            closeErr := f.Close()
            if closeErr != nil {
                return nil, closeErr
            }
			fmt.Println("Denied")
            return nil, err
        }
    }

    return f, nil
} 
