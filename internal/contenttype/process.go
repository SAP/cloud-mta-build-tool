package contenttype

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// GetContentTypes - gets content types associated with files extensions from the configuration config_type_cgf.yaml
func GetContentTypes() (*ContentTypes, error) {
	contentTypes := ContentTypes{}
	err := yaml.Unmarshal(ContentTypeConfig, &contentTypes)
	if err != nil {
		return &contentTypes, errors.Wrap(err, "unmarshalling of the content types failed")
	}
	return &contentTypes, nil
}

// GetContentType - get content type by file extension
func GetContentType(cfg *ContentTypes, extension string) (string, error) {
	for _, ct := range cfg.ContentTypes {
		if ct.Extension == extension {
			return ct.ContentType, nil
		}
	}
	return "", fmt.Errorf("content type for the %s extension is not defined", extension)
}
