package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	path      string
	reader    bool
	chunk_len int
	files     map[string]*bytes.Buffer
	directory string
)

const (
	host string = "localhost"
	port string = "8080"
)

func init() {
	flag.StringVar(&path, "path", "", "Path")
	flag.BoolVar(&reader, "r", false, "Reader switch")
	flag.Parse()
	log.Printf("Path: %s\n", path)
	log.Printf("Reader mode: %s\n", If(reader, "on", "off").(string))
	directory = path + "files/"
	walkthrough(directory)
}

func main() {
	http.HandleFunc("/", downloadFile)

	log.Printf("Listening on port: %s...\n", port)
	log.Printf("Directory: %s\n", directory)

	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func downloadFile(w http.ResponseWriter, r *http.Request) {

	filename := strings.TrimPrefix(r.URL.Path, "/")

	fileByteArr := files[filename].Bytes()
	mimeType := http.DetectContentType(fileByteArr)
	fileSize := len(string(fileByteArr))

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Length", strconv.Itoa(fileSize))

	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(fileByteArr))
}

func readFileWithoutReader(path string) (buffer *bytes.Buffer) {
	start := time.Now()
	log.Printf("Start to read file")

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	buffer = bytes.NewBuffer(make([]byte, 0))
	chunk := make([]byte, 1024)

	for {
		if chunk_len, err = file.Read(chunk); err != nil {
			break
		}
		buffer.Write(chunk[:chunk_len])
	}
	if err != io.EOF {
		log.Fatal("Error: ", err)
	} else {
		err = nil
	}
	elapsed := time.Since(start)
	log.Printf("Read file without bufio reader took %s", elapsed)
	return
}

func readFileWithReader(path string) (buffer *bytes.Buffer) {
	start := time.Now()
	log.Printf("Start to read file")

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	reader := bufio.NewReader(file)
	buffer = bytes.NewBuffer(make([]byte, 0))
	chunk := make([]byte, 1024)

	for {
		if chunk_len, err = reader.Read(chunk); err != nil {
			break
		}
		buffer.Write(chunk[:chunk_len])
	}
	if err != io.EOF {
		log.Fatal("Error: ", err)
	} else {
		err = nil
	}
	elapsed := time.Since(start)
	log.Printf("Read file with bufio reader took %s", elapsed)
	return
}

func walkthrough(path string) {
	r := isFlagPassed("reader")
	files = make(map[string]*bytes.Buffer)
	items, _ := ioutil.ReadDir(path)
	for _, item := range items {
		if !item.IsDir() {
			if r {
				files[item.Name()] = readFileWithReader(path + item.Name())
			} else {
				files[item.Name()] = readFileWithoutReader(path + item.Name())
			}
		} else {
			subitems, _ := ioutil.ReadDir(item.Name())
			for _, subitem := range subitems {
				if !subitem.IsDir() {
					log.Println(item.Name() + "/" + subitem.Name())
				}
			}
		}
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
