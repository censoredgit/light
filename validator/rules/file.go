package rules

import (
	"errors"
	"fmt"
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
	"io"
	"net/http"
)

type FileRule struct {
	Rule
	mime       []string
	parsedMime string
	maxCount   uint8
}

func File(mime []string) *FileRule {
	return &FileRule{
		Rule: Rule{
			alias: "File",
			sType: support.Files,
		},
		mime:     mime,
		maxCount: 1,
	}
}

func (r *FileRule) ParsedMime() string {
	return r.parsedMime
}

func (r *FileRule) SetMaxCount(count uint8) *FileRule {
	r.maxCount = count
	return r
}

func (r *FileRule) Process(inputData *input.Data, fieldName string) error {
	exists, src := inputData.Has(fieldName)
	if !exists {
		return nil
	}

	if src == input.SourceValues {
		return r.Err("The :field field is invalid.", fieldName)
	}

	if len(r.mime) > 0 {
		for _, f := range inputData.AllFiles(fieldName) {

			fl, err := f.Open()
			if err != nil {
				return errors.New(fmt.Sprint(r.Alias(), ": internal error 0x1"))
			}

			fileHeader := make([]byte, 512)
			_, err = fl.Read(fileHeader)
			if err != nil {
				return r.Err("The :field field must be a file.", fieldName)
			}

			isCorrectType := false
			r.parsedMime = http.DetectContentType(fileHeader)
			for _, m := range r.mime {
				if m == r.parsedMime {
					isCorrectType = true
					break
				}
			}
			if !isCorrectType {
				_ = fl.Close()
				return r.Err("File type not supported", fieldName)
			}

			_, err = fl.Seek(0, io.SeekStart)
			if err != nil {
				return errors.New(fmt.Sprint(r.Alias(), ": internal error 0x2"))
			}
			err = fl.Close()
			if err != nil {
				return errors.New(fmt.Sprint(r.Alias(), ": internal error 0x3"))
			}
		}
	}

	return nil
}
