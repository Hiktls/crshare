package main

import (
	"fmt"
	"bytes"
	"log"
	"io"
	"io/ioutil"
	"crypto/md5"
	"strings"
	"encoding/hex"
	"net/http"
)

type Upload struct {
	file string
}



func welcomePage(w http.ResponseWriter, r *http.Request){
	
	if r.Method != "GET"{
		http.Error(w,"GET not supported.",http.StatusNotFound)
	}

	fmt.Fprint(w,"What are you doing here lol")
}

func uploadPage(w http.ResponseWriter, r *http.Request){
	r.ParseMultipartForm(32 << 20) 
    var buf bytes.Buffer
    file, header, err := r.FormFile("file")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    name := strings.Split(header.Filename, ".")
    fmt.Printf("File name %s\n", name[0])
    io.Copy(&buf, file)
    contents := buf.String()
		fileNameHash := md5.Sum([]byte(name[0]+name[1]))
		fileName := "./files/"+hex.EncodeToString(fileNameHash[:])[:6]+"."+name[1]
		errr := ioutil.WriteFile(fileName,[]byte(contents),0644)
		if errr != nil {
			panic(errr)
		}
    buf.Reset()
		fmt.Fprintf(w,fileName)
    return
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "pp2.jpg")
}

func main(){
	 fileServer := http.FileServer(http.Dir("./static"))
	 http.Handle("/",fileServer)
	 http.HandleFunc("/where",welcomePage)
	 http.HandleFunc("/upload",uploadPage)
	 fs := http.FileServer(http.Dir("C:\\Users\\kaano\\go\\crshare\\files"))
	 http.Handle("/files/",http.StripPrefix("/files",fs))
	 fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080",nil); err != nil {
		log.Fatal(err)
	}

}
