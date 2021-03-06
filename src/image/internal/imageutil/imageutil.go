// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package imageutil contains code shared by image-related packages.
package imageutil

import (
	"image"
	"image/color"
)

// DrawYCbCr draws the YCbCr source image on the RGBA destination image with
// r.Min in dst aligned with sp in src. It returns whether the draw was
// successful. If it returns false, no dst pixels were changed.
func DrawYCbCr(dst *image.RGBA, r image.Rectangle, src *image.YCbCr, sp image.Point) (ok bool) {
	// This function exists in the image/internal/imageutil package because it
	// is needed by both the image/draw and image/jpeg packages, but it doesn't
	// seem right for one of those two to depend on the other.
	//
	// Another option is to have this code be exported in the image package,
	// but we'd need to make sure we're totally happy with the API (for the
	// rest of Go 1 compatibility), and decide if we want to have a more
	// general purpose DrawToRGBA method for other image types. One possibility
	// is:
	//
	// func (src *YCbCr) CopyToRGBA(dst *RGBA, dr, sr Rectangle) (effectiveDr, effectiveSr Rectangle)
	//
	// in the spirit of the built-in copy function for 1-dimensional slices,
	// that also allowed a CopyFromRGBA method if needed.

	x0 := (r.Min.X - dst.Rect.Min.X) * 4
	x1 := (r.Max.X - dst.Rect.Min.X) * 4
	y0 := r.Min.Y - dst.Rect.Min.Y
	y1 := r.Max.Y - dst.Rect.Min.Y
	switch src.SubsampleRatio {
	case image.YCbCrSubsampleRatio444:
		for y, sy := y0, sp.Y; y != y1; y, sy = y+1, sy+1 {
			dpix := dst.Pix[y*dst.Stride:]
			yi := (sy-src.Rect.Min.Y)*src.YStride + (sp.X - src.Rect.Min.X)
			ci := (sy-src.Rect.Min.Y)*src.CStride + (sp.X - src.Rect.Min.X)
			for x := x0; x != x1; x, yi, ci = x+4, yi+1, ci+1 {
				rr, gg, bb := color.YCbCrToRGB(src.Y[yi], src.Cb[ci], src.Cr[ci])
				dpix[x+0] = rr
				dpix[x+1] = gg
				dpix[x+2] = bb
				dpix[x+3] = 255
			}
		}
	case image.YCbCrSubsampleRatio422:
		for y, sy := y0, sp.Y; y != y1; y, sy = y+1, sy+1 {
			dpix := dst.Pix[y*dst.Stride:]
			yi := (sy-src.Rect.Min.Y)*src.YStride + (sp.X - src.Rect.Min.X)
			ciBase := (sy-src.Rect.Min.Y)*src.CStride - src.Rect.Min.X/2
			for x, sx := x0, sp.X; x != x1; x, sx, yi = x+4, sx+1, yi+1 {
				ci := ciBase + sx/2
				rr, gg, bb := color.YCbCrToRGB(src.Y[yi], src.Cb[ci], src.Cr[ci])
				dpix[x+0] = rr
				dpix[x+1] = gg
				dpix[x+2] = bb
				dpix[x+3] = 255
			}
		}
	case image.YCbCrSubsampleRatio420:
		for y, sy := y0, sp.Y; y != y1; y, sy = y+1, sy+1 {
			dpix := dst.Pix[y*dst.Stride:]
			yi := (sy-src.Rect.Min.Y)*src.YStride + (sp.X - src.Rect.Min.X)
			ciBase := (sy/2-src.Rect.Min.Y/2)*src.CStride - src.Rect.Min.X/2
			for x, sx := x0, sp.X; x != x1; x, sx, yi = x+4, sx+1, yi+1 {
				ci := ciBase + sx/2
				rr, gg, bb := color.YCbCrToRGB(src.Y[yi], src.Cb[ci], src.Cr[ci])
				dpix[x+0] = rr
				dpix[x+1] = gg
				dpix[x+2] = bb
				dpix[x+3] = 255
			}
		}

	case image.YCbCrSubsampleRatio440:
		for y, sy := y0, sp.Y; y != y1; y, sy = y+1, sy+1 {
			dpix := dst.Pix[y*dst.Stride:]
			yi := (sy-src.Rect.Min.Y)*src.YStride + (sp.X - src.Rect.Min.X)
			ci := (sy/2-src.Rect.Min.Y/2)*src.CStride + (sp.X - src.Rect.Min.X)
			for x := x0; x != x1; x, yi, ci = x+4, yi+1, ci+1 {
				rr, gg, bb := color.YCbCrToRGB(src.Y[yi], src.Cb[ci], src.Cr[ci])
				dpix[x+0] = rr
				dpix[x+1] = gg
				dpix[x+2] = bb
				dpix[x+3] = 255
			}
		}
	default:
		return false
	}
	return true
}
