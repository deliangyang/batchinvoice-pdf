package core

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	qrcodeEncoder "github.com/skip2/go-qrcode"
)

// QRCodeScanner 二维码扫描器
type QRCodeScanner struct {
	reader gozxing.Reader
}

// NewQRCodeScanner 创建二维码扫描器
func NewQRCodeScanner() *QRCodeScanner {
	return &QRCodeScanner{
		reader: qrcode.NewQRCodeReader(),
	}
}

// ScanImage 扫描图像中的二维码
func (s *QRCodeScanner) ScanImage(img image.Image) ([]string, error) {
	// 转换为 gozxing 可用的 BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return nil, fmt.Errorf("failed to create binary bitmap: %w", err)
	}

	// 扫描二维码
	result, err := s.reader.Decode(bmp, nil)
	if err != nil {
		// 没有找到二维码不算错误
		return []string{}, nil
	}

	return []string{result.GetText()}, nil
}

// ScanMultipleImages 扫描多个图像中的二维码
func (s *QRCodeScanner) ScanMultipleImages(images []image.Image) ([]string, error) {
	var allQRCodes []string

	for i, img := range images {
		qrcodes, err := s.ScanImage(img)
		if err != nil {
			log.Printf("Failed to scan image %d: %v", i, err)
			continue
		}
		allQRCodes = append(allQRCodes, qrcodes...)
	}

	return allQRCodes, nil
}

// ExtractQRCodesFromPDF 从PDF数据中提取所有二维码
func ExtractQRCodesFromPDF(pdfData []byte) ([]string, error) {
	// 创建PDF处理器
	pdfProcessor := NewPDFProcessor()

	// 提取图像
	images, err := pdfProcessor.ConvertPDFToImages(pdfData)
	if err != nil {
		return nil, fmt.Errorf("failed to extract images from PDF: %w", err)
	}

	// 创建二维码扫描器
	scanner := NewQRCodeScanner()

	// 扫描所有图像
	qrcodes, err := scanner.ScanMultipleImages(images)
	if err != nil {
		return nil, fmt.Errorf("failed to scan QR codes: %w", err)
	}

	log.Printf("Found %d QR codes in PDF", len(qrcodes))
	return qrcodes, nil
}

// GenerateQRCode 生成二维码图片
func GenerateQRCode(data string, width, height int) ([]byte, error) {
	// 使用 go-qrcode 库生成二维码
	qr, err := qrcodeEncoder.New(data, qrcodeEncoder.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	// 生成图片
	img := qr.Image(width)

	// 将图片编码为 PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), nil
}
