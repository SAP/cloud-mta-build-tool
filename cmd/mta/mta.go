package mta

import "cloud-mta-build-tool/cmd/mta/models"

const (
	PARENT = "../"
)

// SetMtaProp - set mta properties
func SetMtaProp(mtaStruct models.MTA) string {
	mtaDirName := mtaStruct.Id
	parentMtaDir := PARENT + mtaDirName
	return parentMtaDir
}
