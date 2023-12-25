package server

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/RSO-project-Prepih/get-photo-info/database"
	"github.com/RSO-project-Prepih/get-photo-info/handlers"

	pb "github.com/RSO-project-Prepih/get-photo-info/github.com/RSO-project-Prepih/get-photo-info"
)

type photoServer struct {
	pb.UnimplementedPhotoServiceServer
}

func (s *photoServer) GetPhotoInfo(ctx context.Context, req *pb.PhotoRequest) (*pb.PhotoResponse, error) {
	log.Printf("Received photo: %v", req.GetPhoto())

	metadata, err := handlers.GetPhotoInfo(req.GetPhoto())
	if err != nil {
		return nil, err
	}

	// Check coordinates and decide if allowed
	allowed := checkCoordinates(metadata)

	// Optionally save EXIF data to the database
	if allowed {
		err := saveMetadataToDatabase(metadata, req.GetImageId())
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

	// allowed coordinates
	firstCoordinate := []float64{45.714683, 14.058524}
	secondCoordinate := []float64{45.880323, 14.053031}
	thirdCoordinate := []float64{45.887015, 14.302970}
	fourthCoordinate := []float64{45.714683, 14.294730}

	box := [][]float64{firstCoordinate, secondCoordinate, thirdCoordinate, fourthCoordinate}

	// Parse latitude and longitude
	latitude, err := strconv.ParseFloat(metadata["Latitude"], 64)
	if err != nil {
		log.Println("Error parsing latitude:", err)
		return false
	}
	longitude, err := strconv.ParseFloat(metadata["Longitude"], 64)
	if err != nil {
		log.Println("Error parsing longitude:", err)
		return false
	}

	// Check if the coordinates are in the box allowed range
	if latitude < box[0][0] && latitude > box[3][0] && longitude > box[0][1] && longitude < box[1][1] {
		log.Println("Coordinates are in the allowed range")
		return true
	}

	return false
}

func saveMetadataToDatabase(metadata map[string]string, imageID string) error {
	db := database.NewDBConnection()
	defer db.Close()

	// Prepare your query with the correct table and column names
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

	// Execute the statement
	_, err = stmt.Exec(
		imageID,
		metadata["Camera Model"],
		metadata["Date"],
		metadata["Latitude"],
		metadata["Longitude"],
		metadata["Altitude"],
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
