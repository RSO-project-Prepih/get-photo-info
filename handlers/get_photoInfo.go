package handlers

import (
	"bytes"
	"fmt"
	"log"

	"github.com/rwcarlsen/goexif/exif"
)

func GetPhotoInfo(photoBytes []byte) (map[string]string, error) {
	// Extract EXIF metadata, check coordinates, and store in the database
	// Return PhotoResponse

	reader := bytes.NewReader(photoBytes)
	exifData, err := exif.Decode(reader)
	if err != nil {
		log.Println("Error decoding EXIF data:", err)
		return nil, err
	}

	metadata := make(map[string]string)

	// Extract camera model
	if cameraModel, err := exifData.Get(exif.Model); err == nil {
		metadata["Camera Model"] = cameraModel.String()
	}

	// Extract date and time
	if date, err := exifData.Get(exif.DateTime); err == nil {
		metadata["Date"] = date.String()
	}

	// Extract GPS coordinates
	if latitude, longitude, err := exifData.LatLong(); err == nil {
		metadata["Latitude"] = fmt.Sprintf("%v", latitude)
		metadata["Longitude"] = fmt.Sprintf("%v", longitude)
	}

	// Extract altitude
	if altitude, err := exifData.Get(exif.GPSAltitude); err == nil {
		metadata["Altitude"] = altitude.String()
	}

	// Extract exposure time
	if exposureTime, err := exifData.Get(exif.ExposureTime); err == nil {
		metadata["Exposure Time"] = exposureTime.String()
	}

	// Lens model
	if lensModel, err := exifData.Get(exif.LensModel); err == nil {
		metadata["Lens Model"] = lensModel.String()
	}

	return metadata, nil
}
