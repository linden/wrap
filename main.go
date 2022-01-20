package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var ID string
var path string
var bundle string
var origin string

func init() {
	flag.StringVar(&ID, "ID", "", "Apple ID, used for signing")
	flag.StringVar(&path, "path", "", "path to the binary")
	flag.StringVar(&bundle, "bundle", "", "bundle name")
	flag.StringVar(&origin, "origin", "example", "your name, for the signing (no spaces)")
	flag.Parse()

	if ID == "" {
		log.Fatal("ID is required")
	}

	if path == "" {
		log.Fatal("path is required")
	}

	if bundle == "" {
		log.Fatal("bundle is required")
	}
}

func main() {
	err := os.RemoveAll("./" + bundle + ".app")

	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir("./"+bundle+".app", 0755)

	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir("./"+bundle+".app"+"/Contents", 0755)

	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir("./"+bundle+".app"+"/Contents/MacOS", 0755)

	if err != nil {
		log.Fatal(err)
	}

	executable, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalf("failed to find executable at %s", path)
	}

	err = ioutil.WriteFile("./"+bundle+".app"+"/Contents/MacOS/"+bundle, executable, 0755)

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./"+bundle+".app"+"/Contents/Info.plst",
		[]byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
		<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
		<plist version="1.0">
		<dict>
		    <key>CFBundleExecutable</key>
		    <string>%s</string>
		    <key>CFBundleIdentifier</key>
		    <string>com.%s.%s</string>
		    <key>CFBundleInfoDictionaryVersion</key>
		    <string>6.0</string>
		    <key>CFBundleName</key>
		    <string>%s</string>
		    <key>CFBundlePackageType</key>
		    <string>APPL</string>
		    <key>CFBundleShortVersionString</key>
		    <string>1.0.0</string>
		    <key>CFBundleVersion</key>
		    <string>20</string>
		    <key>LSMinimumSystemVersion</key>
		    <string>10.6</string>
		    <key>LSUIElement</key>
		    <true/>
		    <key>LSBackgroundOnly</key>
		    <true/>
		</dict>
		</plist>`, bundle, origin, bundle, bundle)), 0755)

	if err != nil {
		log.Fatal(err)
	}

	var stderr bytes.Buffer

	command := exec.Command("/usr/bin/codesign", "--force", "--verify", "--verbose", "--deep", "--sign", ID, bundle+".app")

	command.Stderr = &stderr

	err = command.Run()

	if err != nil {
		log.Fatalf("failed to sign bundle: %s", stderr.String())
	}
}
