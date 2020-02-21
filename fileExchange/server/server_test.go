package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"servers/fileExchange/pkg/operations"
	"testing"
	"time"
)

func Test_ConnectToServer(t *testing.T) {
	var host = "localhost"
	var port = rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err, _ := serverStart(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
	}()
	time.Sleep(3_000_000_000)
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	defer func() {
		err := dial.Close()
		if err != nil {
			log.Printf("can't close client conn: %v", err)
		}
	}()
}

func Test_DownloadFromServerOK(t *testing.T)  {
	var host = "localhost"
	var port = rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err, conn := serverStart(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
		log.Printf("try download")
		err = downloadFromServer(conn, "krik.txt")
		if err != nil {
			t.Fatalf("can't download from server: %v", err)
		}
		log.Printf("download success")
	}()
	time.Sleep(3_000_000_000)
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	defer func() {
		err := dial.Close()
		if err != nil {
			t.Fatalf("can't close client conn: %v", err)
		}
	}()
	log.Println("try download file from server to client")
	writer := bufio.NewWriter(dial)
	err = operations.WriteLine("download:krik.txt", writer)
	if err != nil {
		t.Fatalf("can't send command to server: %v", err)
	}
	reader := bufio.NewReader(dial)
	line, err := operations.ReadLine(reader)
	if err != nil {
		t.Fatalf("can't read file from server: %v", err)
	}
	fmt.Printf("result file from server: %s", line)
	if line != "download:ok\n" {
		t.Fatalf("can't download file from server: %s", line)
	}
	log.Printf("download file from server success")
}

func Test_DownloadFromServerNotOK(t *testing.T)  {
	var host = "localhost"
	var port = rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err, conn := serverStart(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
		log.Printf("try download")
		err = downloadFromServer(conn, "krik123.txt")
		if err == nil {
			t.Fatalf("just be error: %v", err)
		}
		log.Printf("just be download failed")
	}()
	time.Sleep(3_000_000_000)
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	defer func() {
		err := dial.Close()
		if err != nil {
			t.Fatalf("can't close client conn: %v", err)
		}
	}()
	log.Println("try download file from server to client")
	writer := bufio.NewWriter(dial)
	err = operations.WriteLine("download:krik123.txt", writer)
	if err != nil {
		t.Fatalf("command not sent: %v", err)
	}
	reader := bufio.NewReader(dial)
	line, err := operations.ReadLine(reader)
	if err != nil {
		t.Fatalf("can't read file from server: %v", err)
	}
	fmt.Printf("result file from server: %s", line)
	if line != "error file not found\n" {
		t.Fatalf("just be - error file not found: %s", line)
	}
	log.Printf("just be download failed from server")
}

func Test_UploadToServerOK(t *testing.T)  {
	var host = "localhost"
	var port = rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err, conn := serverStart(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
		log.Printf("try upload")
		err = uploadToServer(conn, "win8.txt")
		if err != nil {
			t.Fatalf("can't upload from server: %v", err)
		}
		log.Printf("upload success")
	}()
	time.Sleep(5_000_000_000)
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	defer func() {
		err := dial.Close()
		if err != nil {
			t.Fatalf("can't close client conn: %v", err)
		}
	}()
	log.Println("try upload file from client to server")
	writer := bufio.NewWriter(dial)
	err = operations.WriteLine("upload:win8.txt", writer)
	if err != nil {
		t.Fatalf("can't send command to server: %v", err)
	}
	log.Println("command sent to server success")
	log.Println("try find uploaded file from server")
	dir, err := ioutil.ReadDir("./files")
	if err != nil {
		t.Fatalf("can't read files_test: ")
	}
	for _, file := range dir{
		if file.Name() == "win8.txt" {
			log.Println("found file uploaded, to server success")
		}
	}
	log.Printf("compare file from client and file uploded to server")
	file, err := os.Open("./upload/win8.txt")
	if err != nil {
		t.Fatalf("can't open file to upload: %v", err)
	}
	byteFile, err := io.Copy(writer, file)
	if err != nil {
		t.Fatalf("can't copy files_test to upload: %v", err)
	}
	log.Printf("copied bytes: %d", byteFile)
	uploadFile, err := os.Open("./files/win8.txt")
	if err != nil {
		t.Fatalf("can't open file to upload: %v", err)
	}
	byteUploadFile, err := io.Copy(writer, uploadFile)
	if err != nil {
		t.Fatalf("can't copy files_test to upload: %v", err)
	}
	if byteFile != byteUploadFile {
		t.Fatalf("file upload incorrect")
	}
	log.Printf("file upload corrected")
}

func Test_UploadToServerNotOK(t *testing.T)  {
	var host = "localhost"
	var port = rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err, conn := serverStart(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
		log.Printf("try upload")
		err = uploadToServer(conn, "command123.txt")
		if err != nil {
			t.Fatalf("can't upload from server: %v", err)
		}
		log.Printf("can't upload file just it")
	}()
	time.Sleep(5_000_000_000)
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	defer func() {
		err := dial.Close()
		if err != nil {
			t.Fatalf("can't close client conn: %v", err)
		}
	}()
	log.Println("try upload file from client to server")
	writer := bufio.NewWriter(dial)
	err = operations.WriteLine("upload:command123.txt", writer)
	if err != nil {
		t.Fatalf("can't send command to server: %v", err)
	}
	log.Println("command sent to server success")
	log.Println("try find uploaded file from server")
	dir, err := ioutil.ReadDir("./files")
	if err != nil {
		t.Fatalf("can't read files_test: ")
	}
	for _, file := range dir{
		if file.Name() == "command123.txt" {
			t.Fatalf("file just not be exist on server")
		}
	}
	log.Printf("file upload incorrect")
}

func ExampleServerFiles()  {
	var host = "localhost"
	var port = rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err, conn := serverStart(addr)
		if err != nil {
			log.Printf("can't start server: %v", err)
			err := listFilesFromServer(conn)
			if err != nil {
				log.Printf("can't getFilesList: %v", err)
			}
		}
	}()
	time.Sleep(3_000_000_000)
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("can't connect to server: %v", err)
	}
	defer func() {
		err := dial.Close()
		if err != nil {
			log.Printf("can't close client conn: %v", err)
		}
	}()
	list, err := operations.ListFiles("./files")
	if err != nil {
		log.Fatalf("can't get server files_test: %v", err)
	}
	fmt.Println(list)
	//Output: krik.txt
	//win8.txt
	//krik.txt
	//win8.txt
}