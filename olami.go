package golami

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

func (client *Client) PostText(ctx context.Context, service, text string) (*Result, error) {
	// get timestamp
	timeStamp := time.Now().Local().UnixNano() / int64(time.Millisecond)
	log.Printf("timestamp: %d\n", timeStamp)

	//  Prepare message to generate an MD5 digest.
	signMsg := fmt.Sprintf("%sapi=%sappkey=%stimestamp=%d%s",
		client.appSecret, service, client.appKey, timeStamp, client.appSecret,
	)

	// Generate MD5 digest.
	sign := fmt.Sprintf("%x", md5.Sum([]byte(signMsg)))

	// Prepare rq JSON data
	var rq string
	var apiName string
	switch service {
	case APIServiceSEG:
		rq = text
		apiName = APIServiceSEG
	case APIServiceNLI:
		rq = fmt.Sprintf(`{"data_type":"stt","data":{"input_type":"1","text":"%s"}}`, text)
		apiName = APIServiceNLI
	}

	// Assemble all the HTTP parameters you want to send
	body := bytes.NewBufferString(fmt.Sprintf("api=%s&appkey=%s&timestamp=%d&sign=%s&rq=%s",
		apiName, client.appKey, timeStamp, sign, rq,
	))

	req, err := http.NewRequest("POST", client.url(), body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return nil, err
	}
	resp, err := client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r := &Result{}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}

func (client *Client) PostASR(ctx context.Context, filePath string) (*Result, error) {
	// get timestamp
	timeStamp := time.Now().Local().UnixNano() / int64(time.Millisecond)
	log.Printf("timestamp: %d\n", timeStamp)

	//  Prepare message to generate an MD5 digest.
	signMsg := fmt.Sprintf("%sapi=%sappkey=%stimestamp=%d%s",
		client.appSecret, APIServiceASR, client.appKey, timeStamp, client.appSecret,
	)

	// Generate MD5 digest.
	sign := fmt.Sprintf("%x", md5.Sum([]byte(signMsg)))

	// file stream
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	client.httpClient.Jar, _ = cookiejar.New(nil)
	req, err := http.NewRequest("POST", client.url(), file)
	req.Header.Set("Content-Type", "application/octet-stream")
	// Generate QueryString
	q := req.URL.Query()
	q.Add("api", APIServiceASR)
	q.Add("appkey", client.appKey)
	q.Add("timestamp", fmt.Sprintf("%d", timeStamp))
	q.Add("sign", sign)
	q.Add("seq", APIServiceSEG+","+APIServiceNLI)
	q.Add("stop", "1")
	// TODO option
	q.Add("compress", "0")

	req.URL.RawQuery = q.Encode()

	if err != nil {
		return nil, err
	}
	_, err = client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	q.Del("compress")
	for {
		resp, err := client.get(ctx, q)
		if err != nil {
			return nil, err
		}
		r := &Result{}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&r)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if r.Data.ASR.Final {
			return r, nil
		}
	}
}
