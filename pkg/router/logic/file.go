package router

import (
	"CloudPlatform/pkg/handle"
	"CloudPlatform/pkg/router"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func init() {
	router.Register(func(router gin.IRoutes) {
		router.POST("/file/upload", Upload)
	}, router.V1)
}

const (
	LIMIT_SIZE = 5 << 20 // 限制上传文件的大小, 这里是5MB

	COS_BUCKET_NAME = "abcp"
	COS_APP_ID      = "1306179590"
	COS_REGION      = "ap-nanjing"
	COS_SECRET_ID   = "AKID61h8KHyQdbKn10jDSSChy0nKbdR0KPEo"
	COS_SECRET_KEY  = "NUMyLSKwJQKzOYS7WgaFfpjOozua0B7a"
	COS_URL_FORMAT  = "http://%s-%s.cos.%s.myqcloud.com" // 此项固定
)

func Upload(c *gin.Context) {
	r := handle.NewResponse(c)
	f, err := c.FormFile("file")
	if err != nil {
		glog.Errorf("file upload error! msg: %s", err.Error())
		r.Error(handle.INTERNAL_ERROR)
		return
	}
	// 大小限制
	if f.Size > LIMIT_SIZE {
		r.Error(handle.PARAM_FILE_SIZE_TOO_BIG)
		return
	}
	file, _ := f.Open()
	name := strings.Split(f.Filename, `.`)

	fileName := base64.StdEncoding.EncodeToString([]byte(name[0])) + name[1]
	url, err := upload(fileName, file)
	if err != nil {
		glog.Errorf("file upload error! msg:", err.Error())
		r.Error(handle.INTERNAL_FILE_UPLOAD_ERROR)
		return
	}
	m := make(map[string]interface{})
	m["file_url"] = url
	r.Success(m)
}

func upload(fileName string, file io.Reader) (string, error) {
	URL := fmt.Sprintf(COS_URL_FORMAT, COS_BUCKET_NAME, COS_APP_ID, COS_REGION)
	u, _ := url.Parse(URL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Timeout: 30 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  COS_SECRET_ID,
			SecretKey: COS_SECRET_KEY,
		},
	})
	key := "/abcp/" + fileName

	_, err := client.Object.Put(
		context.Background(), key, file, nil,
	)
	return URL + key, err
}
