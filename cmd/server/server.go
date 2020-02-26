package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"servers/pkg"
	"strings"
)

func main() {
	const addr = "0.0.0.0:9999"
	log.Printf("server try starting: %s", addr)
	err, _ := serverStart(addr)
	if err != nil {
		log.Fatalf("can't start server on: %s, error: %v", addr, err)
	}
}

func serverStart(addr string) (err error, conn net.Conn) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("can't listen from: %s, it error: %v", addr, err)
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			log.Fatalf("can't close connect: %v", err)
		}
	}()
	log.Printf("server start success!: %s", addr)
	for {
		conn, err := listener.Accept()
		log.Printf("try accept connect on: %s", addr)
		if err != nil {
			log.Printf("can't accept connect %v", err)
			continue
		}
		log.Printf("Server find connect on: %s", addr)
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) (err error) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("can't server close connect: %e", err)
			return
		}
	}()
	log.Printf("connect is stable!")
	reader := bufio.NewReader(conn)
	cmdClient, err := pkg.ReadLine(reader)
	if err != nil {
		log.Printf("error while reading: %v", err)
		return err
	}
	index := strings.IndexByte(cmdClient, ':')
	cmd, arg := cmdClient[:index], cmdClient[index+1:]
	log.Printf("command from client: %s", cmd)
	log.Printf("commandArgument from client: %s", arg)
	switch cmd {
	case "download":
		{
			err := downloadFromServer(conn, arg)
			if err != nil {
				log.Printf("can't downlod file from server to client: %v", err)
			}
		}
	case "upload":
		{
			err := uploadToServer(conn, arg)
			if err != nil {
				log.Printf("can't upload file from client to server: %v", err)
			}
		}
	case "list":
		{
			err := listFilesFromServer(conn)
			if err != nil {
				log.Printf("can't send list files to client: %v", err)
			}
		}
	default:
		fmt.Println("error request from client")
	}
	return nil
}

func downloadFromServer(conn net.Conn, arg string) (err error) {
	writer := bufio.NewWriter(conn)
	log.Printf("client try download file: %s", arg)
	arg = strings.TrimSuffix(arg, "\n")
	rootFiles := "./cmd/server/files/"
	log.Printf("try find file: %s, from directory: %s", arg, rootFiles)
	downloadFile := rootFiles + arg
	file, err := os.Open(downloadFile)
	if err != nil {
		log.Printf("can't find file: %s, it error: %v", arg, err)
		_ = pkg.WriteLine("error file not found", writer)
		return err
	}
	log.Printf("file found: %s, from directory: %s", arg, rootFiles)
	err = pkg.WriteLine("download:ok", writer)
	if err != nil {
		log.Printf("error while writing: %v", err)
		return err
	}
	log.Printf("try copy file: %s, from directory: %v", arg, rootFiles)
	byteFile, err := io.Copy(writer, file)
	if err != nil {
		log.Printf("can't copy file: %s, it error: %v", arg, err)
		return err
	}
	log.Printf("file copied: %s, size by bytes: %d", arg, byteFile)
	log.Printf("file copy success")
	log.Printf("file sent to client success")
	return nil
}

func uploadToServer(conn net.Conn, arg string) (err error) {
	arg = strings.TrimSuffix(arg, "\n")
	log.Printf("client try upload file to server: %s", arg)
	log.Printf("try upload file: %s from client", arg)
	reader := bufio.NewReader(conn)
	_, err = pkg.ReadLine(reader)
	if err != nil {
		log.Printf("can't read: %v", err)
		return err
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		if err != io.EOF {
			log.Printf("can't read data: %v", err)
			return err
		}
	}
	check := string(bytes)
	log.Printf("check file: %s exist in client", arg)
	if check == "error\n" {
		log.Printf("file: %s is not exist in client", arg)
		return err
	}
	downloadDir := "./cmd/server/files/"
	log.Printf("download file: %s, size by bytes: %d, to: %s", arg, len(bytes), downloadDir)
	downloadFile := downloadDir + arg
	log.Printf("try create file: %s", downloadFile)
	file, err := os.OpenFile(downloadFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("can't create file: %s, error: %v", arg, err)
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("can't close file: %v", err)
		}
	}()
	log.Printf("create file: %s, success", downloadFile)
	log.Printf("try write to file: %s", arg)
	_, err = file.Write(bytes)
	if err != nil {
		log.Printf("can't write to file: %s, error: %v", arg, err)
		return err
	}
	log.Printf("file: %s, download success!", arg)
	log.Println("upload to server from client success!")
	return nil
}

func listFilesFromServer (conn net.Conn) (err error) {
	const serverFiles = "./cmd/server/files"
	writeFileLists := bufio.NewWriter(conn)
	log.Printf("try get list from: %s", serverFiles)
	getListFile, err := pkg.ListFiles(serverFiles)
	if err != nil {
		log.Printf("can't get server Files list: %v", err)
		return err
	}
	log.Printf("server Files list get success.")
	log.Printf("try send server Files list to client.")
	err = pkg.WriteLine(getListFile, writeFileLists)
	if err != err {
		log.Printf("can't send server Files list to client: %v", err)
		return err
	}
	log.Printf("server Files list send success to client")
	return nil
}