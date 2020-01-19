package conttype

import (
	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

// GetContentTypes - gets content types associated with files extensions from the configuration config_type_cgf.yaml
func GetContentTypes() (*ContentTypes, error) {
	contentTypes := ContentTypes{}
	err := yaml.UnmarshalStrict(ContentTypeConfig, &contentTypes)
	if err != nil {
		return &contentTypes, errors.Wrap(err, unmarshalFailed)
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
	return "", errors.Errorf(ContentTypeUndefinedMsg, extension)
}
