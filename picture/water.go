package picture

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func Water(filename, watername string) error {
	ext := filepath.Ext(filename)

	var err error
	var img image.Image

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	switch strings.ToLower(ext) {
	case ".jpg":
		img, err = jpeg.Decode(file)
	case ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	}

	if err != nil {
		return err
	}
	file.Close()

	waterfile, err := os.Open(watername)
	if err != nil {
		return err
	}
	waterimg, err := png.Decode(waterfile)
	if err != nil {
		return err
	}
	waterfile.Close()

	// 开始打水印
	b := img.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, img, image.ZP, draw.Src)

	wb := waterimg.Bounds()
	ow := b.Dx()  // 原图宽度
	oh := b.Dy()  // 原图高度
	ww := wb.Dx() // 水印宽度
	wh := wb.Dy() // 水印高度

	margin := 20
	offsetTop := image.Pt(margin, margin)
	offsetMiddle := image.Pt(ow/2-ww/2, oh/2-wh/2)
	offsetBottom := image.Pt(ow-ww-margin, oh-wh-margin)

	// 头部水印
	draw.Draw(m, wb.Add(offsetTop), waterimg, image.ZP, draw.Over)
	draw.Draw(m, wb.Add(offsetMiddle), waterimg, image.ZP, draw.Over)
	draw.Draw(m, wb.Add(offsetBottom), waterimg, image.ZP, draw.Over)

	// 保存
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	switch strings.ToLower(ext) {
	case ".jpg":
		jpeg.Encode(out, m, &jpeg.Options{100})
	case ".jpeg":
		jpeg.Encode(out, m, &jpeg.Options{100})
	case ".png":
		png.Encode(out, m)
	}

	return err
}
