package repositories

import (
	"log"
	"testing"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/config"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/database"
)

func TestGetAllChallenges(t *testing.T) {
	config.Load()
	database.Init()
	challenges, err := GetAllChallenges()
	log.Println(challenges, err)
}
