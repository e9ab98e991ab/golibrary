/**
 * @author : godfeer@aliyun.com
 * @date : 2018/6/11/011
 **/
package base

import (
	"code.google.com/p/graphics-go/graphics"
	//"code.google.com/p/graphics-go/graphics/interp"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
)

type imgDataW struct {
	ImgData []byte
}

func (this *imgDataW) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 {
		this.ImgData = append(this.ImgData, p...)
	}
	return
}

type Img struct {
	SrcImgFile string
}

func (this *Img) getType() string {
	slice := strings.Split(this.SrcImgFile, ".")
	return slice[len(slice)-1]
}

//生成缩略图
func (this *Img) CreateSmall(x, y int) (b []byte, err error) {
	var srcF *os.File
	if srcF, err = os.Open(this.SrcImgFile); err != nil {
		return
	} else {
		defer srcF.Close()
	}

	if x <= 0 && y <= 0 {
		if b, err = ioutil.ReadAll(srcF); err != nil {
			return
		} else {
			return b, nil
		}
	}

	switch imgType := this.getType(); imgType {
	case "png":
		return this.png(srcF, x, y)
	case "jpg":
		return this.jpg(srcF, x, y)
	case "gif":
		return this.gif(srcF, x, y)
	}
	return
}

func (this *Img) png(srcF *os.File, x, y int) (b []byte, err error) {
	var oldImg image.Image
	if oldImg, err = png.Decode(srcF); err != nil {
		return
	}
	oldY := oldImg.Bounds().Dy()
	oldX := oldImg.Bounds().Dx()
	tX, tY := 0, 0
	if x > 0 && y > 0 {
		if (x * oldY / oldX) > y {
			tX = y * oldX / oldY
			tY = y
		} else {
			tX = x
			tY = x * oldY / oldX
		}
	} else {
		if x <= 0 {
			tX = y * oldX / oldY
			tY = y
			x = tX
		} else {
			tX = x
			tY = x * oldY / oldX
			y = tY
		}
	}
	newImg := image.NewRGBA(image.Rect(0, 0, tX, tY)) //
	graphics.Scale(newImg, oldImg)                    //缩略图
	outImg := image.NewRGBA(image.Rect(0, 0, x, y))   //
	white := color.RGBA{255, 255, 255, 255}           //缩略图底色
	draw.Draw(outImg, outImg.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)
	floatX, floatY := 0, 0
	if x-tX > 0 {
		floatX = (x - tX) / 2
	}
	if y-tY > 0 {
		floatY = (y - tY) / 2
	}
	draw.Draw(outImg, image.Rect(floatX, floatY, x, y), newImg, image.ZP, draw.Over)
	imgDataW := &imgDataW{
		ImgData: make([]byte, 0),
	}
	if err = png.Encode(imgDataW, outImg); err != nil {
		return
	}
	return imgDataW.ImgData, nil
}

func (this *Img) jpg(srcF *os.File, x, y int) (b []byte, err error) {
	var oldImg image.Image
	if oldImg, err = jpeg.Decode(srcF); err != nil {
		return
	}
	oldY := oldImg.Bounds().Dy()
	oldX := oldImg.Bounds().Dx()
	tX, tY := 0, 0
	if x > 0 && y > 0 {
		if (x * oldY / oldX) > y {
			tX = y * oldX / oldY
			tY = y
		} else {
			tX = x
			tY = x * oldY / oldX
		}
	} else {
		if x <= 0 {
			tX = y * oldX / oldY
			tY = y
			x = tX
		} else {
			tX = x
			tY = x * oldY / oldX
			y = tY
		}
	}
	newImg := image.NewRGBA(image.Rect(0, 0, tX, tY)) //
	graphics.Scale(newImg, oldImg)                    //缩略图
	outImg := image.NewRGBA(image.Rect(0, 0, x, y))   //
	white := color.RGBA{255, 255, 255, 255}           //缩略图底色
	draw.Draw(outImg, outImg.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)
	floatX, floatY := 0, 0
	if x-tX > 0 {
		floatX = (x - tX) / 2
	}
	if y-tY > 0 {
		floatY = (y - tY) / 2
	}
	draw.Draw(outImg, image.Rect(floatX, floatY, x, y), newImg, image.ZP, draw.Over)
	imgDataW := &imgDataW{
		ImgData: make([]byte, 0),
	}
	if err = jpeg.Encode(imgDataW, outImg, &jpeg.Options{jpeg.DefaultQuality}); err != nil {
		return
	}
	return imgDataW.ImgData, nil
}

func (this *Img) gif(srcF *os.File, x, y int) (b []byte, err error) {
	var oldImg image.Image
	if oldImg, err = gif.Decode(srcF); err != nil {
		return
	}
	oldY := oldImg.Bounds().Dy()
	oldX := oldImg.Bounds().Dx()
	tX, tY := 0, 0
	if x > 0 && y > 0 {
		if (x * oldY / oldX) > y {
			tX = y * oldX / oldY
			tY = y
		} else {
			tX = x
			tY = x * oldY / oldX
		}
	} else {
		if x <= 0 {
			tX = y * oldX / oldY
			tY = y
			x = tX
		} else {
			tX = x
			tY = x * oldY / oldX
			y = tY
		}
	}

	floatX, floatY := 0, 0 //游离边距
	if x-tX > 0 {
		floatX = (x - tX) / 2
	}
	if y-tY > 0 {
		floatY = (y - tY) / 2
	}

	newImg := image.NewRGBA(image.Rect(0, 0, tX, tY)) //
	graphics.Scale(newImg, oldImg)                    //缩略图

	outImg := image.NewRGBA(image.Rect(0, 0, x, y))
	white := new(color.RGBA) //缩略图底色
	white.A = 255
	white.B = 255
	white.G = 255
	white.R = 255
	draw.Draw(outImg, outImg.Bounds(), &image.Uniform{white}, image.ZP, draw.Src) //添加底色底色
	draw.Draw(outImg, image.Rect(floatX, floatY, x, y), newImg, image.ZP, draw.Over)

	imgDataW := &imgDataW{
		ImgData: make([]byte, 0),
	}
	var opt gif.Options
	opt.NumColors = 256
	if err = gif.Encode(imgDataW, outImg, nil); err != nil {
		return
	}
	return imgDataW.ImgData, nil
}
