package initializers

import (
	"log"
	"os"
)

func CreateAssetsFolder() {

	if _, err := os.Stat("assets/images"); os.IsNotExist(err) {
		if err := os.MkdirAll("assets/images", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("assets/audios"); os.IsNotExist(err) {
		if err := os.MkdirAll("assets/audios", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("assets/videos"); os.IsNotExist(err) {
		if err := os.MkdirAll("assets/videos", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

}
