package tripmanagement

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"go-transport-hub/dtos"
	"go-transport-hub/utils"
)

func (tr *TripSheetObj) UploadTripSheetImages(imageFor, tripSheetNumber string, file multipart.File, fileHeader *multipart.FileHeader) (*dtos.UploadTripSheetResponse, error) {

	if imageFor == "" {
		return nil, errors.New("imageFor should not be empty. imageFor:driver_license_image, lr_gate_image")
	}

	imageTypes := [...]string{
		"driver_license_image", "lr_gate_image",
	}
	exists := false
	for _, imgType := range imageTypes {
		if imgType == imageFor {
			exists = true
			break
		}
	}
	if !exists {
		return nil, errors.New("imageFor should not be empty. imageFor:driver_license_image")
	}
	baseDirectory := os.Getenv("BASE_DIRECTORY")
	uploadPath := os.Getenv("IMAGE_DIRECTORY")
	if uploadPath == "" || baseDirectory == "" {
		tr.l.Error("ERROR:  BASE_DIRECTORY &  IMAGE_DIRECTORY found")
		return nil, errors.New("BASE_DIRECTORY & IMAGE_DIRECTORY path not found")
	}
	tripNumber := strings.ReplaceAll(tripSheetNumber, "/", "_")
	tripNumber = strings.ReplaceAll(tripNumber, "-", "_")

	imageDirectory := filepath.Join(uploadPath, "tripsheet", tripNumber)
	fullPath := filepath.Join(baseDirectory, imageDirectory)
	tr.l.Infof("trip sheet path: %v", fullPath)
	err := os.MkdirAll(fullPath, os.ModePerm) // os.ModePerm sets permissions to 0777
	if err != nil {
		tr.l.Error("ERROR: MkdirAll failed for path: ", fullPath, " error: ", err)
		return nil, fmt.Errorf("failed to create directory %s: %w", fullPath, err)
	}
	
	// Verify directory was created
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		tr.l.Error("ERROR: Directory does not exist after MkdirAll: ", fullPath)
		return nil, fmt.Errorf("directory was not created: %s", fullPath)
	}
	
	extension := strings.Split(fileHeader.Filename, ".")
	lengthExt := len(extension)

	imageName := fmt.Sprintf("%v_%v.%v", imageFor, tripNumber, extension[lengthExt-1])
	tsImageFullPath := filepath.Join(fullPath, imageName)
	tr.l.Info("tsImageFullPath: ", imageFor, tsImageFullPath)

	//
	out, err := os.Create(tsImageFullPath)
	if err != nil {
		if utils.CheckFileExists(tsImageFullPath) {
			tr.l.Error("updating  is already exitis: ", tripNumber, uploadPath, err)
		} else {
			tr.l.Error("tsImageFullPath create error: ", tripNumber, uploadPath, err)
			defer out.Close()
			return nil, err
		}
	}
	defer out.Close()

	//
	_, err = io.Copy(out, file)
	if err != nil {
		tr.l.Error("tripNumber upload Copy error: ", tripNumber, uploadPath, err)
		return nil, err
	}

	//
	imageDirectory = filepath.Join(imageDirectory, imageName)
	tr.l.Info("##### to be stored path: ", imageFor, imageDirectory)

	// if tripSheetNumber is exist update
	// Note: if row not found - 0 row(s) affected Rows matched: 0  Changed: 0  Warnings: 0 (no error retun not matched)
	updateQuery := fmt.Sprintf(`UPDATE trip_sheet_header SET %v = '%v' WHERE trip_sheet_num = '%v'`, imageFor, imageDirectory, tripSheetNumber)
	errU := tr.tripSheetDao.UpdateTripSheetImagePath(updateQuery)
	if errU != nil {
		errS := errU.Error()
		tr.l.Error("ERROR: UpdateTripSheetImagePath", tripSheetNumber, errU, errS)
		//return nil, errU
	}

	tr.l.Info("Image uploaded successfully: ", tripSheetNumber)
	roleResponse := dtos.UploadTripSheetResponse{}
	roleResponse.Message = fmt.Sprintf("Image uploaded successfully : %v", tripSheetNumber)
	roleResponse.ImagePath = imageDirectory
	roleResponse.TripSheetNumber = tripSheetNumber
	return &roleResponse, nil
}
