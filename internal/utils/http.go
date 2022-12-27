package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

func HttpGet(url string, timeout int) ([]byte, error) {
	client := http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		// fmt.Printf("http.send.get.method.fail:%s\n", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// fmt.Printf("Get url:%s, ReadAll err:%s", url, err)
		return nil, err
	}

	return body, nil
}

func HttpPost(url string, data interface{}, contentType string, timeout int32) ([]byte, error) {
	return HttpRequest("POST", url, data, contentType, timeout)
}

func HttpPut(url string, data interface{}, contentType string, timeout int32) ([]byte, error) {
	return HttpRequest("PUT", url, data, contentType, timeout)
}

func HttpRequest(method string, url string, data interface{}, contentType string, timeout int32) ([]byte, error) {
	if len(contentType) == 0 {
		contentType = "application/json"
	}

	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		// fmt.Printf(">>>Post url:%s, NewRequest err:%s\n", url, err)
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Add("content-type", contentType)
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// fmt.Printf(">>>Post url:%s, do err:%s\n", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	return result, err
}

func HttpUpload(url string, values map[string]io.Reader, token string, filename string, timeout int32) ([]byte, error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}

		fw, err := w.CreateFormFile(key, filename)
		if _, err = io.Copy(fw, r); err != nil {
			// fmt.Printf(">>>Upload url:%s, Copy err:%s\n", url, err)
			return nil, err
		}
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		// fmt.Printf(">>>Upload url:%s, NewRequest err:%s\n", url, err)
		return nil, err
	}

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Add("file_upload_token", token)

	//忽略证书
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Submit the request
	var client = &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		// fmt.Printf(">>>Upload url:%s, do err:%s\n", url, err)
		return nil, err
	}

	all, err := ioutil.ReadAll(res.Body)
	return all, err
}
