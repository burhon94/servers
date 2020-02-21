package operations

import (
	"log"
	"os"
)

func MkDir(dir string) (err error){
	err = os.Mkdir(dir, 0755)
	if err != nil {
		log.Printf("can't create dir: %v", err)
		return
	}
	return
}
