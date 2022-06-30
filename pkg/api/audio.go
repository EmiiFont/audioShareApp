package api

import (
	"bytes"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Audio struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Size     int64    `json:"size"`
	Location string   `json:"location"`
	Type     string   `json:"type"`
	Tags     []string `json:"tags"`
}

func connectToFirebase() *firebase.App {
	gopath := os.Getenv("GOPATH")
	opt := option.WithCredentialsFile(gopath + "/memeaudio-af80d-firebase-adminsdk-n6lr3-2eaada8830.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		_ = fmt.Errorf("error initializing app: %v", err)
	}
	return app
}

func apiResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"hello world!"}`))
}
func uploadFiles(file *bytes.Buffer) {
	firebaseApp := connectToFirebase()
	fmt.Println(firebaseApp)
	client, err := firebaseApp.Storage(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)

	defer cancel()

	bucket, err := client.Bucket("memeaudio-af80d.appspot.com")
	if err != nil {
		log.Fatalln(err)
	}
	bucketObj := bucket.Object("newguid12")

	//bucketObj = bucketObj.If(storage.Conditions{DoesNotExist: true})
	wc := bucketObj.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		fmt.Errorf("Writer.Close: %v", err)
	}
	//fmt.Fprintf(w, "Blob uploaded.\n")
}

func ListAudios(w http.ResponseWriter, r *http.Request) {
	filesLinks := listFiles()
	audios, _ := json.Marshal(filesLinks)
	w.Write([]byte(audios))
}

func listFiles() []Audio {
	// bucket := "bucket-name"
	ctx := context.Background()
	firebaseApp := connectToFirebase()
	client, err := firebaseApp.Storage(ctx)
	if err != nil {
		fmt.Errorf("storage.NewClient: %v", err)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	it, err := client.Bucket("memeaudio-af80d.appspot.com")

	if err != nil {
		log.Fatalln(err)
	}
	theObjects := it.Objects(ctx, nil)
	var fileSlice []Audio
	for {
		attrs, err := theObjects.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Errorf("Bucket(%q).Objects: %v", "default bucket", err)
		}
		fmt.Println(attrs.Name)
		fileSlice = append(fileSlice, Audio{Name: attrs.Name,
			Type:     attrs.ContentType,
			Size:     attrs.Size,
			Location: attrs.MediaLink,
			ID:       "Id should come from firestore",
		})
	}
	return fileSlice
}

func UploadAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method.", 405)
	}
	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, errFile := r.FormFile("file")
	if errFile != nil {
		http.Error(w, errFile.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(fileHeader)

	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		fmt.Println(err)
	}

	uploadFiles(buf)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}
	name := r.FormValue("name")
	size, _ := strconv.Atoi(r.FormValue("size"))
	fileType := r.FormValue("type")

	theAudio := Audio{
		Name: name,
		Size: int64(size),
		Type: fileType,
	}
	//fileTags, _ := strconv.Atoi(r.FormValue("Reviews"))
	marshal, err := json.Marshal(theAudio)
	w.Write([]byte(marshal))
}
