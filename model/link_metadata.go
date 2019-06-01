// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/dyatlov/go-opengraph/opengraph"
)

const (
	LINK_METADATA_TYPE_IMAGE     LinkMetadataType = "image"
	LINK_METADATA_TYPE_NONE      LinkMetadataType = "none"
	LINK_METADATA_TYPE_OPENGRAPH LinkMetadataType = "opengraph"
	MAX_IMAGES                   int              = 5
)

type LinkMetadataType string

// checks if original string has more than 300 characters long, if so
// it'll truncate it and add an ellipsis (it will be stripped within
// the client, but in case it isn't we might as well add it so it is
// properly understood)
func truncateText(original string) string {
	if utf8.RuneCountInString(original) > 300 {
		return fmt.Sprintf("%.300s[...]", original)
	}
	return original
}

func firstNImages(images []*opengraph.Image, maxImages int) []*opengraph.Image {
	if maxImages < 0 { // dont break stuff, if it's weird, go for sane defaults
		maxImages = MAX_IMAGES
	}
	numImages := len(images)
	if numImages > maxImages {
		subImages := make([]*opengraph.Image, maxImages)
		subImages = images[0:maxImages]
		return subImages
	}
	return images
}

// TruncateOpenGraph modifies a OG into a smaller version to ensure it
// doesn't grow too big, as that much text won't be displayed by the
// clients anyway. Also remove unwanted fields
func TruncateOpenGraph(ogdata *opengraph.OpenGraph) *opengraph.OpenGraph {
	if ogdata != nil {
		// we might want to truncate url too, but that can have unintended effect
		if ogdata.Title != "" {
			ogdata.Title = truncateText(ogdata.Title)
		}
		if ogdata.Description != "" {
			ogdata.Description = truncateText(ogdata.Description)
		}
		if ogdata.SiteName != "" {
			ogdata.SiteName = truncateText(ogdata.SiteName)
		}
		if ogdata.Article != nil {
			ogdata.Article = nil
		}
		if ogdata.Book != nil {
			ogdata.Book = nil
		}
		if ogdata.Profile != nil {
			ogdata.Profile = nil
		}
		if ogdata.Determiner != "" {
			ogdata.Determiner = ""
		}
		if ogdata.Locale != "" {
			ogdata.Locale = ""
		}
		if ogdata.LocalesAlternate != nil {
			ogdata.LocalesAlternate = make([]string, 0)
		}
		if len(ogdata.Images) > 0 {
			ogdata.Images = firstNImages(ogdata.Images, MAX_IMAGES)
		}
		if len(ogdata.Audios) > 0 {
			ogdata.Audios = make([]*opengraph.Audio, 0)
		}
		if len(ogdata.Videos) > 0 {
			ogdata.Videos = make([]*opengraph.Video, 0)
		}

	}
	return ogdata
}

// LinkMetadata stores arbitrary data about a link posted in a message. This includes dimensions of linked images
// and OpenGraph metadata.
type LinkMetadata struct {
	// Hash is a value computed from the URL and Timestamp for use as a primary key in the database.
	Hash int64

	URL       string
	Timestamp int64
	Type      LinkMetadataType

	// Data is the actual metadata for the link. It should contain data of one of the following types:
	// - *model.PostImage if the linked content is an image
	// - *opengraph.OpenGraph if the linked content is an HTML document
	// - nil if the linked content has no metadata
	Data interface{}
}

func (o *LinkMetadata) PreSave() {
	o.Hash = GenerateLinkMetadataHash(o.URL, o.Timestamp)
}

func (o *LinkMetadata) IsValid() *AppError {
	if o.URL == "" {
		return NewAppError("LinkMetadata.IsValid", "model.link_metadata.is_valid.url.app_error", nil, "", http.StatusBadRequest)
	}

	if o.Timestamp == 0 || !isRoundedToNearestHour(o.Timestamp) {
		return NewAppError("LinkMetadata.IsValid", "model.link_metadata.is_valid.timestamp.app_error", nil, "", http.StatusBadRequest)
	}

	switch o.Type {
	case LINK_METADATA_TYPE_IMAGE:
		if o.Data == nil {
			return NewAppError("LinkMetadata.IsValid", "model.link_metadata.is_valid.data.app_error", nil, "", http.StatusBadRequest)
		}

		if _, ok := o.Data.(*PostImage); !ok {
			return NewAppError("LinkMetadata.IsValid", "model.link_metadata.is_valid.data_type.app_error", nil, "", http.StatusBadRequest)
		}
	case LINK_METADATA_TYPE_NONE:
		if o.Data != nil {
			return NewAppError("LinkMetadata.IsValid", "model.link_metadata.is_valid.data_type.app_error", nil, "", http.StatusBadRequest)
		}
	case LINK_METADATA_TYPE_OPENGRAPH:
		if o.Data == nil {
			return NewAppError("LinkMetadata.IsValid", "model.link_metadata.is_valid.data.app_error", nil, "", http.StatusBadRequest)
		}

		if _, ok := o.Data.(*opengraph.OpenGraph); !ok {
			return NewAppError("LinkMetadata.IsValid", "model.link_metadata.is_valid.data_type.app_error", nil, "", http.StatusBadRequest)
		}
	default:
		return NewAppError("LinkMetadata.IsValid", "model.link_metadata.is_valid.type.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

// DeserializeDataToConcreteType converts o.Data from JSON into properly structured data. This is intended to be used
// after getting a LinkMetadata object that has been stored in the database.
func (o *LinkMetadata) DeserializeDataToConcreteType() error {
	var b []byte
	switch t := o.Data.(type) {
	case []byte:
		// MySQL uses a byte slice for JSON
		b = t
	case string:
		// Postgres uses a string for JSON
		b = []byte(t)
	}

	if b == nil {
		// Data doesn't need to be fixed
		return nil
	}

	var data interface{}
	var err error

	switch o.Type {
	case LINK_METADATA_TYPE_IMAGE:
		image := &PostImage{}

		err = json.Unmarshal(b, &image)

		data = image
	case LINK_METADATA_TYPE_OPENGRAPH:
		og := &opengraph.OpenGraph{}

		json.Unmarshal(b, &og)

		data = og
	}

	if err != nil {
		return err
	}

	o.Data = data

	return nil
}

// FloorToNearestHour takes a timestamp (in milliseconds) and returns it rounded to the previous hour in UTC.
func FloorToNearestHour(ms int64) int64 {
	t := time.Unix(0, ms*int64(1000*1000))

	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location()).UnixNano() / int64(time.Millisecond)
}

// isRoundedToNearestHour returns true if the given timestamp (in milliseconds) has been rounded to the nearest hour in UTC.
func isRoundedToNearestHour(ms int64) bool {
	return FloorToNearestHour(ms) == ms
}

// GenerateLinkMetadataHash generates a unique hash for a given URL and timestamp for use as a database key.
func GenerateLinkMetadataHash(url string, timestamp int64) int64 {
	hash := fnv.New32()

	// Note that we ignore write errors here because the Hash interface says that its Write will never return an error
	binary.Write(hash, binary.LittleEndian, timestamp)
	hash.Write([]byte(url))

	return int64(hash.Sum32())
}
