package libs

import (
	"log"
)

func DbErr(err error) {
	if err != nil {
		log.Fatal("database error: %s",err)
	}	
}

func ResponseErr(err error) {
	if err != nil {
		log.Fatal("response error: %s",err)
	}
}

