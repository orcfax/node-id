package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ross-spencer/node-id/pkg/identity"
)

// appname is used by goreleaser.
const appname = "node-id"

const idLoc = "/tmp/.node-identity.json"
const timeFormat = "2006-01-02 15:04:05"

var (
	version = "dev-0.0.0"
	commit  = "000000000000000000000000000000000baddeed"
	date    = "1970-01-01T00:00:01Z"

	// command line args.
	vers      bool
	inspect   bool
	list      bool
	websocket string
)

func init() {
	flag.BoolVar(&list, "list", false, "list the current identity information")
	flag.BoolVar(&list, "l", false, "list the current identity information")
	flag.StringVar(&websocket, "ws", "", "websocket to communicate with")
	flag.StringVar(&websocket, "w", "", "websocket to communicate with")
	flag.BoolVar(&inspect, "inspect", false, "inspect full IP Info data (read-only)")
	flag.BoolVar(&inspect, "i", false, "inspect full IP Info data (read-only)")
	flag.BoolVar(&vers, "version", false, "return version")
	flag.BoolVar(&vers, "v", false, "return version")
}

var usage string = fmt.Sprintf(`Usage of %s:
  -l, --list list the current identity information
  -w, --ws websocket to communicate with
  -i, --inspect inspect full IP Info data (read-only)
  -v, --version output version
  -h, --help prints help information
`, appname)

// listIdentity lists the current identity information.
func listIdentity() {
	var ident identity.Identity
	var err error
	if !identity.Exists(idLoc) {
		log.Println(
			"identity has not yet being created, please create one and try again",
		)
		return
	}
	// Load the identity to enable it to be updated.
	ident, err = identity.LoadCache(idLoc)
	if err != nil {
		// In future implementations we need the handshake to
		// determine whether or not a new identity can just be
		// created.
		log.Printf("error loading existing identity: '%s' cannot retrieve previous data", err)
		return
	}
	val, _ := json.MarshalIndent(ident, "", "   ")
	fmt.Println(string(val))
}

// outputIdentity outputs a new, or updates an existing identity for the
// given machine.
func outputIdentity(websocket string) error {
	var ident identity.Identity
	var err error
	if !identity.Exists(idLoc) {
		log.Println(
			"identity has not yet being created, please create one and try again",
		)
	}
	// Load the identity to enable it to be updated.
	ident, err = identity.LoadCache(idLoc)
	if err != nil {
		// In future implementations we need the handshake to
		// determine whether or not a new identity can just be
		// created.
		log.Printf("error loading existing identity: '%s' cannot retrieve previous data", err)
	}

	if (ident == identity.Identity{}) {
		log.Printf("creating a new id")
	} else {
		log.Printf("id exists: '%s'", ident.NodeID)
		log.Printf("first initialized: '%s'", ident.InitializationDate)
	}

	loc := identity.GetIdentity(
		ident.NodeID,
		ident.InitializationDate,
		websocket,
	)
	val, _ := json.MarshalIndent(loc, "", "   ")
	log.Println(string(val))
	log.Printf("outputting to: '%s'", idLoc)
	err = os.WriteFile(idLoc, val, 0644)
	return err
}

func main() {
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
	if vers {
		fmt.Fprintf(os.Stderr, "%s commit: %s date: %s\n", version, commit, date)
		os.Exit(0)
	} else if flag.NFlag() <= 0 { // can access args w/ len(os.Args[1:]) too
		fmt.Fprintf(os.Stderr, "Usage:  %s\n", appname)
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-list] ...")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-ws] ... ")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-inspect] ... ")
		fmt.Fprintln(os.Stderr, "        OPTIONAL: [-version] ...")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Output: [STRING] {collector-identity}")
		fmt.Fprintf(os.Stderr, "Output: [FILE]   %s\n", idLoc)
		fmt.Fprintln(os.Stderr, "Output: [STRING] {IP Info Data}")
		fmt.Fprintf(os.Stderr, "Output: [STRING] '%s ...'\n\n", version)
		flag.Usage()
		os.Exit(0)
	}

	if list {
		listIdentity()
		os.Exit(0)
	}

	if inspect {
		ipInfo := identity.IPInfoDefault()
		val, _ := json.MarshalIndent(ipInfo, "", "   ")
		fmt.Println(string(val))
		os.Exit(0)
	}

	if websocket == "" {
		log.Println("please provide a valid websocket value, cannot 'nil'")
		os.Exit(1)
	}
	err := outputIdentity(websocket)
	if err != nil {
		log.Printf("problem writing id file: '%s'", err)
		os.Exit(1)
	}
}
