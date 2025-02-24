package encoder

import "encoding/base64"

type Content interface {
	Get() string
}

type content struct {
	data string
}

func (s content) Get() string {
	return s.data
}

func StringContent(data string) Content {
	return content{data: data}
}

func BytesContent(data []byte) Content {
	return content{data: base64.StdEncoding.EncodeToString(data)}
}
