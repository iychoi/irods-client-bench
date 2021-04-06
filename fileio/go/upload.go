package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cyverse/go-irodsclient/irods/common"
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
		fmt.Fprintf(os.Stderr, "Give a local source path and an iRODS destination path!\n")
		os.Exit(1)
	}

	srcPath := args[0]
	destPath := args[1]

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
	appName := "upload"
	timeout := time.Second * 200 // 200 sec

	conn := connection.NewIRODSConnection(account, timeout, appName)
	err = conn.Connect()
	if err != nil {
		util.LogErrorf("err - %v", err)
		panic(err)
	}
	defer conn.Disconnect()

	// convert src path into absolute path
	srcPath, err = filepath.Abs(srcPath)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buffersize := 1024 * 1024 * 8 // 8MB

	handle, err := fs.OpenDataObjectWithOperation(conn, destPath, "", "w", common.OPER_TYPE_PUT_DATA_OBJ)
	if err != nil {
		panic(err)
	}

	// copy
	buffer := make([]byte, buffersize)
	for {
		bytesRead, err := f.Read(buffer)
		if err != nil {
			fs.CloseDataObject(conn, handle)
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		err = fs.WriteDataObject(conn, handle, buffer[:bytesRead])
		if err != nil {
			fs.CloseDataObject(conn, handle)
			panic(err)
		}
	}
}
