package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cyverse/go-irodsclient/irods/connection"
	"github.com/cyverse/go-irodsclient/irods/fs"
	"github.com/cyverse/go-irodsclient/irods/types"
	"github.com/cyverse/go-irodsclient/irods/util"
)

func main() {
	util.SetLogLevel(9)

	// Parse cli parameters
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Give an iRODS source path and a local destination path!\n")
		os.Exit(1)
	}

	srcPath := args[0]
	destPath := args[1]
	destFilePath := destPath

	stat, err := os.Stat(destPath)
	if err != nil {
		if os.IsNotExist(err) {
			// file not exists, it's a file
			destFilePath = destPath
		} else {
			panic(err)
		}
	} else {
		if stat.IsDir() {
			irodsFileName := filepath.Base(srcPath)
			destFilePath = filepath.Join(destPath, irodsFileName)
		} else {
			panic(fmt.Errorf("File %s already exists", destPath))
		}
	}

	// Read account configuration from YAML file
	yaml, err := ioutil.ReadFile("account.yml")
	if err != nil {
		panic(err)
	}

	account, err := types.CreateIRODSAccountFromYAML(yaml)
	if err != nil {
		panic(err)
	}

	// Create a file system
	appName := "download"
	timeout := time.Second * 200 // 200 sec

	conn := connection.NewIRODSConnection(account, timeout, appName)
	err = conn.Connect()
	if err != nil {
		util.LogErrorf("err - %v", err)
		panic(err)
	}
	defer conn.Disconnect()

	// convert dest path into absolute path
	destFilePath, err = filepath.Abs(destFilePath)
	if err != nil {
		util.LogErrorf("err - %v", err)
		panic(err)
	}

	handle, _, err := fs.OpenDataObject(conn, srcPath, "", "r")
	if err != nil {
		panic(err)
	}

	f, err := os.Create(destFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buffersize := 1024 * 1024 * 8 // 8MB

	// copy
	for {
		buffer, err := fs.ReadDataObject(conn, handle, buffersize)
		if err != nil {
			fs.CloseDataObject(conn, handle)
			panic(err)
		}

		if buffer == nil || len(buffer) == 0 {
			// EOF
			fs.CloseDataObject(conn, handle)
			return
		} else {
			_, err = f.Write(buffer)
			if err != nil {
				fs.CloseDataObject(conn, handle)
				panic(err)
			}
		}
	}
}
