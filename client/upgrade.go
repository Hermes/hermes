package client

import (
	"fmt"
	"net/http"
	"encoding/json"
	"os"
	"io"
	"bytes"
	"crypto/md5"
	"encoding/hex"
)

type ClientVersion struct {
	Version 	float32 	// current client version
	Build 		int 		// current client build
	Checksum 	string 		// checksum of client version
}

func checkVersion(o ClientVersion) (ClientVersion, bool) {

	var n ClientVersion

	// Checking for latest version
	resp, err := http.Get("https://raw.github.com/olanmatt/hermes/master/build/version")
	if err != nil {
		fmt.Println("Error: No able to get latest version")
	}
	defer resp.Body.Close()

	// Parsing JSON
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()
	b := []byte(s)
	err = json.Unmarshal(b, &n)
	if err != nil {
		fmt.Println(err)
	}

	// Version comparison
	if o.Version < n.Version {
		return n, true
	} else if o.Version == n.Version && o.Build < n.Build {
		return n, true
	}
	return o, false
}

func Upgrade(v float32, b int) {
	version, result := checkVersion(ClientVersion{v, b, ""})
	if result {

		fmt.Println("New version of client found")

		// Creating file container
		out, err := os.Create("hermes_")
		if err != nil {
			fmt.Println(err)
		}
		defer out.Close()

		// Retrieving latest build
		resp, err := http.Get("https://raw.github.com/olanmatt/hermes/master/build/hermes")
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		io.Copy(out, resp.Body)

    	// Verifying checksum
    	file, _ := os.Open("hermes_")
		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
    	h := md5.New()
    	io.WriteString(h, buf.String())
    	if hex.EncodeToString(h.Sum(nil)) == version.Checksum {

			// Set file permissions
			err = os.Chmod("hermes_", 0777)
			if err != nil {
				fmt.Println(err)
			}
			os.Remove("hermes")
			os.Rename("hermes_", "hermes")
		} else {
			os.Remove("hermes_")
			fmt.Println("Error: Failed to update")
		}

	} else {
		fmt.Println("No new version of client found")
	}

}