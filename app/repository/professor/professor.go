package professor

import (
	"errors"
	"github.com/google/uuid"

	"github.com/armanjr/termustat/app/config"
	"github.com/armanjr/termustat/app/models"
	"github.com/armanjr/termustat/app/utils"
	"gorm.io/gorm"
)

func GetOrCreate(universityID uuid.UUID, rawName string) (uuid.UUID, error) {
	normalizedName := utils.NormalizeProfessor(rawName)
	if normalizedName == "" {
		return uuid.Nil, errors.New("invalid professor name after normalization")
	}

	var prof models.Professor
	err := config.DB.Where(
		"university_id = ? AND normalized_name = ?",
		universityID,
		normalizedName,
	).First(&prof).Error

	if err == nil {
		return prof.ID, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newProf := models.Professor{
			UniversityID:   universityID,
			Name:           rawName,
			NormalizedName: normalizedName,
		}

		if err := config.DB.Create(&newProf).Error; err != nil {
			return uuid.Nil, err
		}
		return newProf.ID, nil
	}

	return uuid.Nil, err
}
