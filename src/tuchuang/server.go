package tuchuang

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	HostName        string `json:"host_name"`
}

func HandleError(err error) {
	log.Printf("%#v", err)
}

var config Config

const LOG = `0.0.1 2021-07-06 实现基本功能
0.0.2 2021年7月10日 试图修复 CROS 问题`

var VERSION = func() string {
	lines := strings.Split(LOG, "\n")
	var afLines = []string{}
	for _, line := range lines {
		if line != "" {
			afLines = append(afLines, line)
		}
	}
	lastLine := afLines[len(afLines)-1]
	return strings.Split(lastLine, " ")[0]
}()

func init() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		HandleError(err)
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		HandleError(err)
	}
}

func GetBucket() (*oss.Bucket, error) {
	client, err := oss.New(
		config.Endpoint,
		config.AccessKeyId,
		config.AccessKeySecret)
	if err != nil {
		HandleError(err)
	}
	return client.Bucket(config.BucketName)
}

func PutSimple(data []byte, fileName, contentType string) (url string, err error) {
	bucket, err := GetBucket()
	if err != nil {
		HandleError(err)
	}
	url = ""
	uuidHex := utils.GetUUID()
	if len(uuidHex) > 5 {
		uuidHex = uuidHex[:4]
	}
	timeStr := time.Now().Format("20060102")
	fileDirName := timeStr + "/" + uuidHex + "_" + fileName
	option := oss.ContentType(contentType)
	err = bucket.PutObject(fileDirName, bytes.NewReader(data), option)
	if err != nil {
		HandleError(err)
	}
	url = config.HostName + "/" + fileDirName
	return
}

func HandlePut(w http.ResponseWriter, r *http.Request) {
	handleCROS(w)
	result := map[string]string{}
	setMsg := func(content string) {
		result["message"] = content
		marshal, _ := json.Marshal(result)
		_, _ = w.Write(marshal)
	}
	if r.Method != "POST" {
		setMsg("Unsupported HTTP Method, Just support PUT")
		return
	}
	_ = r.ParseForm()
	file, m, err := r.FormFile("file")
	if err != nil {
		setMsg(fmt.Sprintf("%#v", err))
		return
	}
	filename := m.Filename
	defer func(file multipart.File) {
		if file.Close() != nil {
			setMsg(fmt.Sprintf("%#v", err))
		}
	}(file)
	if m.Size > 1024*1024*5 {
		setMsg(fmt.Sprintf("File too big, size %d", m.Size))
		return
	}
	temp, err := ioutil.ReadAll(file)
	if err != nil {
		setMsg(fmt.Sprintf("%#v", err))
		return
	}
	contentType := http.DetectContentType(temp)
	url, err := PutSimple(temp, filename, contentType)
	if err != nil {
		setMsg(fmt.Sprintf("%#v", err))
	} else {
		result["message"] = "OK"
		result["url"] = url
		marshal, _ := json.Marshal(result)
		_, _ = w.Write(marshal)
	}
}

func handleCROS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://static2.mazhangjing.com")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
}

var port = flag.Int("-port", 8089, "HTTP Server Port")

func Serve() {
	flag.Parse()
	log.Printf("Server served at port %d", *port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal(map[string]string{
			"usage":   "Post file to /api/upload with multipart/form-data and named file: file",
			"log":     LOG,
			"version": VERSION,
		})
		_, _ = w.Write(data)
	})
	http.HandleFunc("/api/upload", HandlePut)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
