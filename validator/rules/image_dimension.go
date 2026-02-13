package rules

import (
	"errors"
	"fmt"
	"github.com/censoredgit/light/validator/input"
	"github.com/censoredgit/light/validator/support"
	"image"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
)

type ImageDimensionRule struct {
	Rule
	maxSize image.Point
	minSize image.Point
}

func ImageDimension(maxSize image.Point) *ImageDimensionRule {
	return &ImageDimensionRule{
		Rule: Rule{
			alias: "ImageDimension",
			sType: support.Files,
		},
		maxSize: maxSize,
	}
}

func (r *ImageDimensionRule) SetMin(minSize image.Point) *ImageDimensionRule {
	r.minSize = minSize
	return r
}

func (r *ImageDimensionRule) Process(inputData *input.Data, fieldName string) error {
	exists, src := inputData.Has(fieldName)
	if !exists {
		return nil
	}

	if src == input.SourceValues {
		return r.Err("The :field field is invalid.", fieldName)
	}

	for _, imageData := range inputData.AllFiles(fieldName) {
		fileRule := File([]string{"image/jpeg", "image/png", "image/webp", "image/bmp"})
		err := fileRule.Process(inputData, fieldName)
		if err != nil {
			return err
		}

		var imgBounds image.Rectangle

		var imgDecoder = jpeg.Decode
		switch fileRule.ParsedMime() {
		case "image/png":
			imgDecoder = png.Decode
		case "image/webp":
			imgDecoder = webp.Decode
		case "image/bmp":
			imgDecoder = bmp.Decode
		}

		file, err := imageData.Open()
		if err != nil {
			return errors.New(fmt.Sprint(r.Alias(), ": internal error 0x1"))
		}

		img, err := imgDecoder(file)
		if err != nil {
			_ = file.Close()
			return r.Err("Image type not supported", fieldName)
		}

		imgBounds = img.Bounds()
		err = file.Close()
		if err != nil {
			return errors.New(fmt.Sprint(r.Alias(), ": internal error 0x2"))
		}

		if imgBounds.Max.X > r.maxSize.X || imgBounds.Max.Y > r.maxSize.Y {
			return r.Err(fmt.Sprintf("Max size for :field is %dx%d", r.maxSize.X, r.maxSize.Y), fieldName)
		}

		if imgBounds.Max.X < r.minSize.X || imgBounds.Max.Y < r.minSize.Y {
			return r.Err(fmt.Sprintf("Min size for :field is %dx%d", r.minSize.X, r.minSize.Y), fieldName)
		}
	}

	return nil
}
