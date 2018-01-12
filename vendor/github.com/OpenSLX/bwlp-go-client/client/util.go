package client

import (
	"github.com/OpenSLX/bwlp-go-client/bwlp"
)

type BwlpSpecifics struct {
	Details *bwlp.ImageDetailsRead
	Version *bwlp.ImageVersionDetails
}

func GetLatestVersionDetails(imageDetails *bwlp.ImageDetailsRead) *bwlp.ImageVersionDetails {
	if imageDetails == nil {
		return nil
	}
	for _, version := range imageDetails.Versions {
		if version.VersionId == imageDetails.LatestVersionId {
			return version
		}
	}
	return nil
}
