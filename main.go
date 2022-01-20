package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
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

	err = ioutil.WriteFile("./"+bundle+".app"+"/Contents/Info.plst", []byte(strings.Join([]string{
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>",
		"<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">",
		"<plist version=\"1.0\">",
		"<dict>",
		"\t<key>CFBundleExecutable</key>",
		"\t<string>" + bundle + "</string>",
		"\t<key>CFBundleIdentifier</key>",
		"\t<string>com." + origin + "." + bundle + "</string>",
		"\t<key>CFBundleInfoDictionaryVersion</key>",
		"\t<string>6.0</string>",
		"\t<key>CFBundleName</key>",
		"\t<string>" + bundle + "</string>",
		"\t<key>CFBundlePackageType</key>",
		"\t<string>APPL</string>",
		"\t<key>CFBundleShortVersionString</key>",
		"\t<string>1.0.0</string>",
		"\t<key>CFBundleVersion</key>",
		"\t<string>20</string>",
		"\t<key>LSMinimumSystemVersion</key>",
		"\t<string>10.6</string>",
		"\t<key>LSUIElement</key>",
		"\t<true/>",
		"\t<key>LSBackgroundOnly</key>",
		"\t<true/>",
		"</dict>",
		"</plist>",
	}, "\r\n")), 0755)

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
