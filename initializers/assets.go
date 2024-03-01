package initializers

import (
	"log"
	"os"
)

var DebugBasePath string = "./assets/"
var ReleaseBasePath string = "/data/assets/"

func CreateAssetsFolder(basePath string) {

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

}
