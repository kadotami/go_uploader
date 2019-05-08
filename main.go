package main

import (
    "html/template"
    "io"
    "net/http"
    "os"
    "fmt"
)

const root_path = "/var/www/go_uploader/"
const upload_path = "/var/file/upload/"

func saveHandler(w http.ResponseWriter, r *http.Request) {
    reader, err := r.MultipartReader()

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    for {
        part, err := reader.NextPart()
        if err == io.EOF {
            break
        }

        if part.FileName() == "" {
            continue
        }

        uploadedFile, err := os.Create(upload_path + part.FileName())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            uploadedFile.Close()
            redirectToErrorPage(w,r)
            return
        }

        _, err = io.Copy(uploadedFile, part)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            uploadedFile.Close()
            redirectToErrorPage(w,r)
            return
        }
    }
    http.Redirect(w,r,"/complete",http.StatusFound)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    var templatefile = template.Must(template.ParseFiles(root_path + "upload.html"))
    templatefile.Execute(w, "upload.html")
}

func completeHandler(w http.ResponseWriter, r *http.Request) {
    var templatefile = template.Must(template.ParseFiles(root_path + "complete.html"))
    templatefile.Execute(w, "complete.html")
}

func errorPageHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w,"%s","<p>Internal Server Error</p>")
}

func redirectToErrorPage(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w,r,"/errorPage",http.StatusFound)
}

func main() {
    //ハンドラの登録
    http.HandleFunc("/", uploadHandler)
    http.HandleFunc("/upload", uploadHandler)
    http.HandleFunc("/save",saveHandler)
    http.HandleFunc("/complete",completeHandler)
    http.HandleFunc("/errorPage",errorPageHandler)
    //サーバーの開始
    http.ListenAndServe(":80", nil)
}
