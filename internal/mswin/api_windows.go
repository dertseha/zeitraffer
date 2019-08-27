package mswin

import (
	"fmt"
	"syscall"
	"unsafe"
)

// HANDLE represents a generic handle.
type HANDLE uintptr

// HWND is a window handle.
type HWND HANDLE

// HGDIOBJ is a handle for a graphical object.
type HGDIOBJ HANDLE

// HDC is a handle for a device context.
type HDC HANDLE

// HBITMAP is a handle for a in-memory bitmap.
type HBITMAP HANDLE

// BITMAPINFOHEADER describes a .bmp file.
type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

// Various MS Windows API constants.
// nolint: golint
const (
	SM_XVIRTUALSCREEN  = 76
	SM_YVIRTUALSCREEN  = 77
	SM_CXVIRTUALSCREEN = 78
	SM_CYVIRTUALSCREEN = 79
	HORZRES            = 8
	VERTRES            = 10
	BI_RGB             = 0
	InvalidParameter   = 2
	DIB_RGB_COLORS     = 0
	SRCCOPY            = 0x00CC0020

	DPI_AWARENESS_CONTEXT_UNAWARE      = -1
	DPI_AWARENESS_CONTEXT_SYSTEM_AWARE = -2
)

// BitBlt calls the Windows API function.
func BitBlt(hdcDest HDC, nXDest, nYDest, nWidth, nHeight int, hdcSrc HDC, nXSrc, nYSrc int, dwRop uint) {
	ret, _, _ := procBitBlt.Call(
		uintptr(hdcDest), uintptr(nXDest), uintptr(nYDest), uintptr(nWidth), uintptr(nHeight),
		uintptr(hdcSrc), uintptr(nXSrc), uintptr(nYSrc), uintptr(dwRop))
	if ret == 0 {
		panic("BitBlt failed")
	}
}

// CreateCompatibleBitmap calls the Windows API function.
func CreateCompatibleBitmap(hdc HDC, cx, cy int) (HBITMAP, error) {
	ret, _, err := procCreateCompatibleBitmap.Call(uintptr(hdc), uintptr(cx), uintptr(cy))
	if ret == 0 {
		return 0, err
	}
	return HBITMAP(ret), nil
}

// CreateCompatibleDC calls the Windows API function.
func CreateCompatibleDC(hdc HDC) (HDC, error) {
	ret, _, err := procCreateCompatibleDC.Call(uintptr(hdc))
	if ret == 0 {
		return 0, err
	}
	return HDC(ret), nil
}

// DeleteDC calls the Windows API function.
func DeleteDC(hdc HDC) bool {
	ret, _, _ := procDeleteDC.Call(uintptr(hdc))
	return ret != 0
}

// DeleteObject calls the Windows API function.
func DeleteObject(hObject HGDIOBJ) bool {
	ret, _, _ := procDeleteObject.Call(uintptr(hObject))
	return ret != 0
}

// GetDC calls the Windows API function.
func GetDC(hwnd HWND) HDC {
	ret, _, _ := procGetDC.Call(uintptr(hwnd))
	return HDC(ret)
}

// GetDeviceCaps calls the Windows API function.
func GetDeviceCaps(hdc HDC, index int) int {
	ret, _, _ := procGetDeviceCaps.Call(uintptr(hdc), uintptr(index))
	return int(ret)
}

// GetSystemMetrics calls the Windows API function.
func GetSystemMetrics(index int) int {
	ret, _, _ := procGetSystemMetrics.Call(uintptr(index))
	return int(ret)
}

// GetDpiForSystem calls the Windows API function.
func GetDpiForSystem() int {
	ret, _, _ := procGetDpiForSystem.Call()
	return int(ret)
}

// GetDIBits calls the Windows API function.
func GetDIBits(hdc HDC, hBitmap HBITMAP, start uint32, cLines uint32,
	bits unsafe.Pointer, bitmapInfo *BITMAPINFOHEADER, usage uint32) bool {
	ret, _, err := procGetDIBits.Call(uintptr(hdc), uintptr(hBitmap), uintptr(start), uintptr(cLines),
		uintptr(bits), uintptr(unsafe.Pointer(bitmapInfo)), uintptr(usage))
	if ret == 0 {
		fmt.Printf("getdibits err: %v\n", err)
	}
	return ret != 0
}

// ReleaseDC calls the Windows API function.
func ReleaseDC(hwnd HWND, hDC HDC) bool {
	ret, _, _ := procReleaseDC.Call(uintptr(hwnd), uintptr(hDC))
	return ret != 0
}

// SelectObject calls the Windows API function.
func SelectObject(hdc HDC, hgdiobj HGDIOBJ) (HGDIOBJ, error) {
	ret, _, err := procSelectObject.Call(uintptr(hdc), uintptr(hgdiobj))
	if ret == 0 {
		return 0, err
	}
	return HGDIOBJ(ret), nil
}

// SetThreadDpiAwarenessContext calls the Windows API function.
func SetThreadDpiAwarenessContext(value int32) error {
	_, _, err := procSetThreadDpiAwarenessContext.Call(uintptr(value))
	return err
}

// nolint: gochecknoglobals
var (
	libgdi32                         = syscall.NewLazyDLL("gdi32.dll")
	libuser32                        = syscall.NewLazyDLL("user32.dll")
	procSetThreadDpiAwarenessContext = libuser32.NewProc("SetThreadDpiAwarenessContext")
	procGetDpiForSystem              = libuser32.NewProc("GetDpiForSystem")
	procGetSystemMetrics             = libuser32.NewProc("GetSystemMetrics")
	procGetDC                        = libuser32.NewProc("GetDC")
	procReleaseDC                    = libuser32.NewProc("ReleaseDC")
	procDeleteDC                     = libgdi32.NewProc("DeleteDC")
	procBitBlt                       = libgdi32.NewProc("BitBlt")
	procDeleteObject                 = libgdi32.NewProc("DeleteObject")
	procSelectObject                 = libgdi32.NewProc("SelectObject")
	procGetDIBits                    = libgdi32.NewProc("GetDIBits")
	procCreateCompatibleBitmap       = libgdi32.NewProc("CreateCompatibleBitmap")
	procCreateCompatibleDC           = libgdi32.NewProc("CreateCompatibleDC")
	procGetDeviceCaps                = libgdi32.NewProc("GetDeviceCaps")
)
