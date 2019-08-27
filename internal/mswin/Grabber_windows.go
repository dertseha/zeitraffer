package mswin

import (
	"image"
	"unsafe"
)

// Grabber is the type that provides images from the desktop
type Grabber struct {
	winDC  HDC
	memDC  HDC
	bmp    HBITMAP
	oldBmp HGDIOBJ
	xVirt  int
	yVirt  int

	buffer []byte
	img    *image.RGBA
}

// NewGrabber creates a new instance of a grabber.
func NewGrabber() (*Grabber, error) {
	_ = SetThreadDpiAwarenessContext(DPI_AWARENESS_CONTEXT_SYSTEM_AWARE)
	winDC := GetDC(0)
	memDC, err := CreateCompatibleDC(winDC)
	if err != nil {
		ReleaseDC(0, winDC)
		return nil, err
	}

	xVirt := GetSystemMetrics(SM_XVIRTUALSCREEN)
	yVirt := GetSystemMetrics(SM_YVIRTUALSCREEN)
	width := GetSystemMetrics(SM_CXVIRTUALSCREEN)
	height := GetSystemMetrics(SM_CYVIRTUALSCREEN)
	// dpi := GetDpiForSystem()
	// fmt.Printf("bounds: %vx%v - %vx%v - dpi %v\n", xVirt, yVirt, width, height, dpi)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	bmp, err := CreateCompatibleBitmap(winDC, width, height)
	if err != nil {
		DeleteDC(memDC)
		ReleaseDC(0, winDC)
		return nil, err
	}
	oldBmp, err := SelectObject(memDC, HGDIOBJ(bmp))
	if err != nil {
		DeleteDC(memDC)
		ReleaseDC(0, winDC)
		return nil, err
	}

	return &Grabber{
		winDC:  winDC,
		memDC:  memDC,
		bmp:    bmp,
		oldBmp: oldBmp,
		xVirt:  xVirt,
		yVirt:  yVirt,

		buffer: make([]byte, len(img.Pix)+3),
		img:    img,
	}, nil
}

// Dispose releases any allocated resources.
func (grabber *Grabber) Dispose() {
	DeleteDC(grabber.memDC)
	DeleteObject(HGDIOBJ(grabber.bmp))
	ReleaseDC(0, grabber.winDC)
}

// Grab makes a snapshot and returns the image. The returned image is valid
// until the grabber is used again.
func (grabber *Grabber) Grab() image.Image {
	width := grabber.img.Rect.Max.X
	height := grabber.img.Rect.Max.Y

	var header BITMAPINFOHEADER
	header.BiSize = uint32(unsafe.Sizeof(header))
	header.BiPlanes = 1
	header.BiBitCount = 32
	header.BiWidth = int32(width)
	header.BiHeight = int32(-height)
	header.BiCompression = BI_RGB
	header.BiSizeImage = 0

	BitBlt(grabber.memDC, 0, 0, width, height, grabber.winDC, grabber.xVirt, grabber.yVirt, SRCCOPY)

	bufferAddr := uintptr(unsafe.Pointer(&(grabber.buffer[0])))
	alignedAddr := ((bufferAddr + 3) / 4) * 4
	_, _ = SelectObject(grabber.memDC, grabber.oldBmp)
	GetDIBits(grabber.memDC, grabber.bmp, 0, uint32(height), unsafe.Pointer(alignedAddr), &header, DIB_RGB_COLORS)
	_, _ = SelectObject(grabber.memDC, HGDIOBJ(grabber.bmp))
	i := 0
	srcBuffer := grabber.buffer[alignedAddr-bufferAddr:]
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// BGRA => RGBA, and set A to 255
			grabber.img.Pix[i], grabber.img.Pix[i+1], grabber.img.Pix[i+2], grabber.img.Pix[i+3] =
				srcBuffer[i+2], srcBuffer[i+1], srcBuffer[i], 255
			i += 4
		}
	}

	return grabber.img
}
