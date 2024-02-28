package initializers

import (
	"log"
	"os"
)

var DebugBasePath string = "./assets/"
var ReleaseBasePath string = "/data/assets/"

func CreateAssetsFolder(basePath string) {

	if _, err := os.Stat(basePath + "images"); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath+"images", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(basePath + "videos"); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath+"videos", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

}
