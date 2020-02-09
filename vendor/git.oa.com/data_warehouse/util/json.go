package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func LoadJson(r io.Reader, v interface{}) error {
	d := json.NewDecoder(r)
	if err := d.Decode(v); err != nil {
		return err
	}
	return nil
}

func LoadJsonFromFile(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return LoadJson(f, v)
}

func SaveJsonToFile(path string, v interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	e := json.NewEncoder(f)
	if err := e.Encode(v); err != nil {
		return err
	}
	return nil
}

func LoadJsonFromURL(url string, v interface{}) error {
	cli := &http.Client{}
	resp, err := cli.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Get %s , status code=%d", url, resp.StatusCode)
	}
	defer resp.Body.Close()
	return LoadJson(resp.Body, v)
}

func JsonReader(v interface{}) (io.Reader, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}

func JsonString(o interface{}) string {
	buffer, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(buffer)
}

func WriteJson(w io.Writer, o interface{}) error {
	ec := json.NewEncoder(w)
	return ec.Encode(o)
}

func PostJson(url string, o interface{}) (*http.Response, error) {
	body, err := JsonReader(o)
	if err != nil {
		return nil, err
	}
	cli := &http.Client{}
	resp, err := cli.Post(url, "application/json", body)
	return resp, err
}

type Time time.Time

const (
	timeFormart = "2006-01-02 15:04:05"
)

func (t Time) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(timeFormart))
	return []byte(stamp), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	if err != nil {
		log.Printf("Parse time error: %v", err)
		return nil
	}
	*t = Time(now)
	return nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormart)
}
