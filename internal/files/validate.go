// Package files provide functionality to validate, store, get... files
package files

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/shared/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	"github.com/gabriel-vasile/mimetype"
)

type FileSizeUnit string

const (
	FileSizeUnitKB FileSizeUnit = "KB"
	FileSizeUnitMB FileSizeUnit = "MB"
	FileSizeUnitGB FileSizeUnit = "GB"
)

type AttachmentValidationConfig struct {
	Files              []*pb.Attachment
	AllowedTypes       []string
	MaxSize            int // must be in bytes
	Unit               FileSizeUnit
	ValidateDiminsions bool // set true for images only
	ImgMaxWidth        int
	ImgMaxHeight       int
}

type AttachmentsValidateError struct {
	ID  string // file id
	Err *models.AppErrorError
}

// AttachmentsValidateSizeAndTypes validate attachment type, size, data
func AttachmentsValidateSizeAndTypes(cfg *AttachmentValidationConfig) *AttachmentsValidateError {
	files := cfg.Files
	allowedTypes := cfg.AllowedTypes
	maxSize := cfg.MaxSize
	unit := cfg.Unit
	validateDim := cfg.ValidateDiminsions

	for _, file := range files {
		base := utils.CleanBase64(file.GetBase64())
		data, err := base64.StdEncoding.DecodeString(base)
		if err != nil {
			return &AttachmentsValidateError{ID: file.Id, Err: &models.AppErrorError{ID: "image.data.invalid"}}
		}

		if len(data) > maxSize {
			size := fmt.Sprintf("%.2f", getAppropriateSize(maxSize, unit))
			return &AttachmentsValidateError{ID: file.Id, Err: &models.AppErrorError{ID: "file.size.error", Params: map[string]any{"Max": size, "Unit": string(unit)}}}
		}

		mime := mimetype.Detect(data)
		fmt.Println(mime)
		fmt.Println(len(file.GetBase64()))
		allowed := false
		for _, at := range allowedTypes {
			if strings.HasPrefix(mime.String(), at) {
				allowed = true
			}
		}
		if !allowed {
			return &AttachmentsValidateError{ID: file.Id, Err: &models.AppErrorError{ID: "image.type.unsupported", Params: map[string]any{"Types": strings.Join(allowedTypes, ", ")}}}
		}

		if strings.HasPrefix(mime.String(), "image/") {
			imgCfg, _, err := image.DecodeConfig(bytes.NewBuffer(data))
			if err != nil {
				return &AttachmentsValidateError{ID: file.Id, Err: &models.AppErrorError{ID: "image.data.invalid"}}
			}

			if validateDim {
				if imgCfg.Width > cfg.ImgMaxWidth || imgCfg.Height > cfg.ImgMaxHeight {
					params := map[string]any{"Dimensions": fmt.Sprintf("%d X %d", cfg.ImgMaxWidth, cfg.ImgMaxHeight)}
					return &AttachmentsValidateError{ID: file.Id, Err: &models.AppErrorError{ID: "image.dimensions.error", Params: params}}
				}
			}

			if file.Crop == nil {
				file.Crop = &pb.Crop{}
			}
			file.Crop.Width = float32(imgCfg.Width)
			file.Crop.Height = float32(imgCfg.Height)
		}

		file.FileSize = int64(len(data))
		file.Mime = mime.String()
		file.Data = data
	}
	return nil
}

// getAppropriateSize convert maxSize in bytes to a display friendly one (for user)
func getAppropriateSize(maxSize int, unit FileSizeUnit) float64 {
	switch unit {
	case FileSizeUnitKB:
		return float64(maxSize / 1024)
	case FileSizeUnitMB:
		return float64(maxSize / 1024 / 1024)
	case FileSizeUnitGB:
		return float64(maxSize / 1024 / 1024 / 1024)
	default:
		return float64(maxSize / 1024)
	}
}
