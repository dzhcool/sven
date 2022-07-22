package picture

import (
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

func Thumb(filename string, maxWidth, maxHeight int) error {
	ext := filepath.Ext(filename)

	var err error
	switch strings.ToLower(ext) {
	case ".jpg":
		err = thumbJpg(filename, uint(maxWidth), uint(maxHeight))
	case ".jpeg":
		err = thumbJpg(filename, uint(maxWidth), uint(maxHeight))
	case ".png":
		err = thumbPng(filename, uint(maxWidth), uint(maxHeight))
	}

	return err
}

func thumbJpg(filename string, maxWidth, maxHeight uint) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}
	file.Close()

	m := resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos2)

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	jpeg.Encode(out, m, nil)

	return nil
}

func thumbPng(filename string, maxWidth, maxHeight uint) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	img, err := png.Decode(file)
	if err != nil {
		return err
	}
	file.Close()

	m := resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos2)

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	png.Encode(out, m)

	return nil
}
