package server

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/RSO-project-Prepih/get-photo-info/database"
	"github.com/RSO-project-Prepih/get-photo-info/handlers"

	pb "github.com/RSO-project-Prepih/get-photo-info/github.com/RSO-project-Prepih/get-photo-info"
)

type photoServer struct {
	pb.UnimplementedPhotoServiceServer
}

// GetPhotoInfo godoc
// @Summary Get photo information
// @Description Retrieves information about a photo
// @Tags photo
// @Accept  json
// @Produce  json
// @Param photo body pb.PhotoRequest true "Photo Request"
// @Success 200 {object} string
// @Router /photo/info [post]
func (s *photoServer) GetPhotoInfo(ctx context.Context, req *pb.PhotoRequest) (*pb.PhotoResponse, error) {
	log.Printf("Received request for GetPhotoInfo with image id: %s", req.GetImageId())
	log.Printf("Received photo image id: %s", req.GetImageId())

	metadata, err := handlers.GetPhotoInfo(req.GetPhoto())
	if err != nil {
		return nil, err
	}

	// Check coordinates and decide if allowed
	allowed := checkCoordinates(metadata)

	if allowed {

		err = saveMetadataToDatabase(metadata, req.GetImageId())
		if err != nil {
			log.Println("Error saving metadata to the database:", err)
			return nil, err
		}

		// Format metadata
		exifData, err := formatMetadata(metadata)
		if err != nil {
			log.Println("Error formatting metadata:", err)
			return nil, err
		}

		return &pb.PhotoResponse{
			Allowed:  allowed,
			ExifData: exifData,
		}, nil

	}

	return &pb.PhotoResponse{
		Allowed: false,
	}, nil
}

func NewServer() *photoServer {
	return &photoServer{}
}

func formatMetadata(metadata map[string]string) (string, error) {
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		log.Println("Error marshalling metadata:", err)
		return "", err
	}
	return string(jsonData), nil
}

func checkCoordinates(metadata map[string]string) bool {
	log.Println("Checking coordinates...")

	// Define the bounding box
	minLatitude := 45.714683
	maxLatitude := 45.887015
	minLongitude := 14.053031
	maxLongitude := 14.302970

	// Parse latitude and longitude from metadata
	latitude, err := strconv.ParseFloat(metadata["Latitude"], 64)
	log.Println("Latitude:", latitude)
	if err != nil {
		log.Println("Error parsing latitude:", err)
		return false
	}
	longitude, err := strconv.ParseFloat(metadata["Longitude"], 64)
	log.Println("Longitude:", longitude)
	if err != nil {
		log.Println("Error parsing longitude:", err)
		return false
	}

	// Check if the coordinates are within the bounding box
	if latitude >= minLatitude && latitude <= maxLatitude && longitude >= minLongitude && longitude <= maxLongitude {
		log.Println("Coordinates are in the allowed range")
		return true
	}

	log.Println("Coordinates are not in the allowed range")
	return false
}

func saveMetadataToDatabase(metadata map[string]string, imageID string) error {
	log.Println("Saving metadata to the database...")
	db := database.NewDBConnection()
	defer db.Close()

	query := `
        INSERT INTO photo_metadata (
            image_id, camera_model, date, latitude, longitude, 
            altitude, exposure_time, lens_model
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	// Prepare the statement for execution
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing database statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Function to safely convert string to float64 for SQL
	safeConvert := func(s string) interface{} {
		if s == "" {
			return nil
		}
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Printf("Error converting string to float64: %v", err)
			return nil
		}
		return val
	}

	// Execute the statement
	_, err = stmt.Exec(
		imageID,
		metadata["Camera Model"],
		strings.Trim(metadata["Date"], "\""),
		safeConvert(metadata["Latitude"]),
		safeConvert(metadata["Longitude"]),
		safeConvert(metadata["Altitude"]),
		metadata["Exposure Time"],
		metadata["Lens Model"],
	)
	if err != nil {
		log.Printf("Error executing database statement: %v", err)
		return err
	}

	log.Println("Metadata saved to the database successfully")
	return nil
}
