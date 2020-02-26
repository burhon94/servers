package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"servers/pkg"
	"strings"
)

var cmd = flag.String("cmd", "download", "command:")

var arg = flag.String("file", "html.pdf", "file name:")

func main() {
	const addr = "localhost:9999"
	log.Printf("client try connecting to: %s", addr)
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("can't dial to: %s, it error: %v", addr, err)
	}
	log.Printf("client connected to %s", addr)
	client(dial)
}

func client(dial net.Conn) {
	defer func() {
		err := dial.Close()
		if err != nil {
			log.Printf("can't close connect , it error: %v", err)
		}
	}()
	flag.Parse()
	*cmd = strings.ToLower(*cmd)
	switch *cmd {
	case "download":
		{
			err := downloadClient(dial, *cmd)
			if err != nil {
				log.Printf("can't download file from server, it error: %v", err)
			}
		}
	case "upload":
		{
			err := uploadClient(dial, *cmd)
			if err != nil {
				log.Printf("can't upload file from client: %v", err)
			}
		}
	case "list":
		{
			err := listClient(dial, *cmd)
			if err != nil {
				log.Printf("can't get files list from server: %v", err)
			}
		}
	default:
		fmt.Printf("incorrect command: %s\n", *cmd)
		return
	}
}

func downloadClient(dial net.Conn, cmd string) (err error) {
	cmd = cmd + ":"
	cmd = cmd + *arg
	cmdSend := bufio.NewWriter(dial)
	log.Printf("try send command to server: %s: ", cmd)
	err = pkg.WriteLine(cmd, cmdSend)
	if err != nil {
		log.Printf("can' t send command: %s, to server, error: %v", cmd, err)
		return err
	}
	log.Println("command send success")
	log.Printf("try download file: %s from server", *arg)
	reader := bufio.NewReader(dial)
	line, err := pkg.ReadLine(reader)
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
	check := line
	log.Printf("check file: %s exist in server", *arg)
	if check == "error file not found\n" {
		log.Printf("file: %s is not exist in server", *arg)
		return err
	}
	log.Printf("file: %s exist in server", *arg)
	downloadDir := "./cmd/client/download/"
	log.Printf("check dir status")
	_, err = os.Stat(downloadDir)
	if err != nil {
		log.Printf("try create directory to save: %s", downloadDir)
		err := pkg.MkDir(downloadDir)
		if err != nil {
			log.Printf("can't create directory: %s, error: %v", downloadDir, err)
			return err
		}
		log.Printf("create directory success: %s", downloadDir)
	}
	log.Printf("download file size by bytes: %d", len(bytes))
	downloadFile := downloadDir + *arg
	log.Printf("try create file: %s", *arg)
	file, err := os.OpenFile(downloadFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("can't create file: %s, error: %v", *arg, err)
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("can't close file: %v", err)
			return
		}
	}()
	log.Printf("create file: %s success", *arg)
	log.Printf("try write to file: %s", *arg)
	_, err = file.Write(bytes)
	if err != nil {
		log.Printf("can't write to file: %s, error: %v", *arg, err)
		return err
	}
	log.Printf("file: %s, download success!", *arg)
	return nil
}

func uploadClient(dial net.Conn, cmd string) (err error) {
	log.Printf("You need have directory \"./cmd/client/upload/\" and put there your files.")
	cmd = cmd + ":"
	cmd = cmd + *arg
	cmdSend := bufio.NewWriter(dial)
	log.Printf("try send command to server: %s: ", cmd)
	err = pkg.WriteLine(cmd, cmdSend)
	if err != nil {
		log.Printf("can' t send command: %s, to server, error: %v", cmd, err)
		return err
	}
	log.Println("command send success")
	log.Printf("try upload file: %s to server", *arg)
	writer := bufio.NewWriter(dial)
	err = pkg.WriteLine("upload: ok", writer)
	if err != nil {
		log.Printf("error while writing: %v", err)
		return err
	}
	*arg = strings.TrimSuffix(*arg, "\n")
	uploadDir := "./cmd/client/upload/"
	uploadFile := uploadDir + *arg
	file, err := os.Open(uploadFile)
	if err != nil {
		log.Printf("can't find file: %s, it error: %v", *arg, err)
		_ = pkg.WriteLine("error", writer)
		return err
	}
	log.Printf("try copy file: %s, from directory: %v", *arg, uploadDir)
	byteFile, err := io.Copy(writer, file)
	if err != nil {
		log.Printf("can't copy file: %s, it error: %v", *arg, err)
		return err
	}
	log.Printf("copy file: %s, size by bytes: %d", *arg, byteFile)
	log.Printf("file copy success")
	log.Printf("file: %s send to server success", *arg)
	return nil
}

func listClient(dial net.Conn, cmd string) (err error) {
	cmd = cmd + ":"
	cmdSend := bufio.NewWriter(dial)
	log.Printf("try send command to server: %s", cmd)
	err = pkg.WriteLine(cmd, cmdSend)
	if err != nil {
		log.Printf("can't  send command: %s to server", cmd)
		return err
	}
	log.Println("command send success")
	log.Println("try get files list from server")
	reader := bufio.NewReader(dial)
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
		return err
	}
	log.Printf("available files:\n%s\n", string(bytes))
	return nil
}
