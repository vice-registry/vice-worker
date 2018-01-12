package client

import (
	"log"
	"errors"

	"github.com/OpenSLX/bwlp-go-client/bwlp"
)

func (handler* SessionHandler) GetImageList(page int32) ([]*bwlp.ImageSummaryRead, error) {
	imageList, err := handler.GetSatClient().GetImageList(handler.SessionData.AuthToken, nil, page)
	if err != nil {
		log.Printf("Failed to retrieve image list: %s\n", err)
		return nil, err
	}
	return imageList, nil
}
func (handler* SessionHandler) GetImageDetails(imageBaseID string) (*bwlp.ImageDetailsRead, error) {
	imageDetails, err := handler.GetSatClient().GetImageDetails(handler.SessionData.AuthToken, bwlp.UUID(imageBaseID))
	if err != nil {
		log.Printf("Failed to retrieve image details for '%s': %s\n", imageBaseID, err)
		return nil, err
	}
	return imageDetails, nil
}
func (handler* SessionHandler) GetImageData(imageBaseID string) (*Transfer, error) {
	// request image details from the satellite
	imageDetails, err := handler.GetSatClient().GetImageDetails(handler.SessionData.AuthToken, bwlp.UUID(imageBaseID))
	if err != nil {
		return nil, err
	}
	// TODO handle versions, for now just use latest one
  var imageVersion *bwlp.ImageVersionDetails = nil
  for _, version := range imageDetails.Versions {
    if version.VersionId == imageDetails.LatestVersionId {
      imageVersion = version
    }
  }
  if imageVersion == nil {
    return nil, errors.New("Latest version not found in image version list, this should not happen :)")
  }
  // Request download of that version
  ti, err := handler.GetSatClient().RequestDownload(handler.SessionData.AuthToken, imageVersion.VersionId)
  if err != nil {
    log.Printf("Error requesting download of image version '%s': %s\n", imageVersion.VersionId, err)
    return nil, err
  }
	return NewTransfer(false, handler.satEndpoint.Hostname, ti, imageVersion.FileSize), nil
}

func (handler* SessionHandler) CreateImage(name string) (*bwlp.UUID, error) {
	newImageBaseId, err := handler.GetSatClient().CreateImage(handler.SessionData.AuthToken, name)
	if err != nil {
		log.Printf("Error creating new base image: %s\n", err)
		return nil, err
	}
	return &newImageBaseId, nil
}

// TODO blockHashes
func (handler* SessionHandler) UploadImageVersion(imageBaseId *bwlp.UUID, fileSize int64, machineDescription []byte) (*Transfer, error) {
	ti, err := handler.GetSatClient().RequestImageVersionUpload(handler.SessionData.AuthToken, *imageBaseId, fileSize, nil, machineDescription)
	if err != nil {
		log.Printf("Error creating version upload: %s\n", err)
		return nil, err
	}
	return NewTransfer(true, handler.satEndpoint.Hostname, ti, fileSize), nil
}

func (handler* SessionHandler) CancelUpload(transfer *Transfer) error {
	if err := handler.GetSatClient().CancelUpload(transfer.Ti.Token); err != nil {
		log.Printf("Error cancelling upload: %s\n", err)
		return err
	}
	log.Printf("Cancelled upload for %s\n", transfer.Ti.Token)
	return nil
}
