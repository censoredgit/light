package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"net/url"
)

func init() {
	gob.Register(ErrorBag{})
}

func valuesToBase64(data url.Values) (string, error) {
	strBuf := bytes.NewBuffer([]byte{})

	benc := base64.NewEncoder(base64.RawStdEncoding, strBuf)
	genc := gob.NewEncoder(benc)
	err := genc.Encode(data)
	if err != nil {
		return "", err
	}

	err = benc.Close()
	if err != nil {
		return "", err
	}

	return strBuf.String(), nil
}

func base64ToValues(data string) (url.Values, error) {
	strBuf := bytes.NewBuffer([]byte(data))
	bdec := base64.NewDecoder(base64.RawStdEncoding, strBuf)
	var values url.Values
	err := gob.NewDecoder(bdec).Decode(&values)
	return values, err
}

func errorBagToBase64(data *ErrorBag) (string, error) {
	strBuf := bytes.NewBuffer([]byte{})

	benc := base64.NewEncoder(base64.RawStdEncoding, strBuf)
	genc := gob.NewEncoder(benc)
	err := genc.Encode(data.err)
	if err != nil {
		return "", err
	}

	err = benc.Close()
	if err != nil {
		return "", err
	}

	return strBuf.String(), nil
}

func base64ToErrorBag(data string) (*ErrorBag, error) {
	strBuf := bytes.NewBuffer([]byte(data))
	bdec := base64.NewDecoder(base64.RawStdEncoding, strBuf)
	values := newErrorBag()
	err := gob.NewDecoder(bdec).Decode(&values.err)
	return values, err
}
