package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Record struct {
	CourseID  string `json:"course_id"`
	Name      string `json:"name"`
	Weight    string `json:"weight"`
	Capacity  string `json:"capacity"`
	Gender    string `json:"gender"`
	Professor string `json:"professor"`
	Faculty   string `json:"faculty"`
	Time1     string `json:"time1"`
	Time2     string `json:"time2"`
	Time3     string `json:"time3"`
	Time4     string `json:"time4"`
	Time5     string `json:"time5"`
	TimeExam  string `json:"time_exam"`
	DateExam  string `json:"date_exam"`
}

func main() {
	router := gin.New()
	router.POST("/process", processUploadedFile)
	log.Println("Starting engine...")
	if err := router.Run(":80"); err != nil {
		log.Fatal("Failed to start engine server", err)
	}
}

func processUploadedFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file uploaded",
		})
		return
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".html") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only HTML files are allowed",
		})
		return
	}

	tempDir, err := os.MkdirTemp("", "course_processing")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create temporary directory",
		})
		return
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save uploaded file",
		})
		return
	}

	records, err := processHTMLFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process HTML file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File processed successfully",
		"records": records,
	})
}

// processAllCourses is used to process all the files in a directory
func processAllCourses() error {
	files, err := os.ReadDir("courses")
	if err != nil {
		return fmt.Errorf("error reading courses directory: %w", err)
	}

	if err := os.MkdirAll("export", os.ModePerm); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	combinedFile, err := os.OpenFile("export/combined.sql", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening combined.sql: %w", err)
	}
	defer combinedFile.Close()

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".html" {
			continue
		}

		faculty := strings.TrimSuffix(file.Name(), ".html")
		filePath := filepath.Join("courses", file.Name())

		records, err := processHTMLFile(filePath)
		if err != nil {
			log.Printf("Skipping %s due to error: %v", filePath, err)
			continue
		}

		if err := generateJSONOutput(faculty, records); err != nil {
			log.Printf("Error generating JSON for %s: %v", faculty, err)
		}

		sqlContent, err := generateSQLInsert(faculty, records)
		if err != nil {
			log.Printf("Error generating SQL for %s: %v", faculty, err)
			continue
		}

		if err := writeSQLFile(faculty, sqlContent); err != nil {
			log.Printf("Error writing SQL file for %s: %v", faculty, err)
		}

		if _, err := combinedFile.WriteString(sqlContent + "\n"); err != nil {
			log.Printf("Error appending to combined.sql: %v", err)
		}
	}

	return nil
}

func processHTMLFile(path string) ([]Record, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	var records []Record

	doc.Find("table tr").Each(func(_ int, row *goquery.Selection) {
		if row.HasClass("DTitle") {
			return
		}

		cells := row.Find("td")
		if cells.Length() < 19 {
			return
		}

		record := Record{
			Faculty:   cleanText(cells.Eq(2).Text()),
			CourseID:  cleanText(cells.Eq(6).Text()),
			Name:      cleanText(cells.Eq(7).Text()),
			Weight:    cleanText(cells.Eq(8).Text()),
			Capacity:  cleanText(cells.Eq(10).Text()),
			Gender:    cleanText(cells.Eq(13).Text()),
			Professor: cleanText(cells.Eq(14).Text()),
		}

		processTimeInfo(&record, cells.Eq(15).Text(), "")
		records = append(records, record)
	})

	return records, nil
}

func processTimeInfo(record *Record, timeStr, examStr string) {
	examStr = strings.Replace(examStr, "تاريخ: ", "امتحان(", 1)
	examStr = strings.Replace(examStr, " ساعت:", ") ساعت :", 1)
	timeStr = strings.ReplaceAll(timeStr+" "+examStr, "\n", "Q")

	processed := strings.NewReplacer(
		"پنج شنبه", "d5/", "پنجشنبه", "d5/",
		"چهار شنبه", "d4/", "چهارشنبه", "d4/",
		"سه شنبه", "d3/", "سهشنبه", "d3/",
		"دو شنبه", "d2/", "دوشنبه", "d2/",
		"يك شنبه", "d1/", "يكشنبه", "d1/",
		"شنبه", "d0/",
		"نيمه۱ ت", "", "نيمه۲ ت", "",
		"درس(ت):", "c", "درس(ع):", "c", "درس (ت):", "c", "درس (ع):", "c",
		"حلتمرين(ت):", "c", "حلتمرین(ت):", "c", "حل تمرين(ت):", "c", "حل تمرین(ت):", "c", "حل تمرين (ت):", "c", "حل تمرین (ت):", "c",
		"امتحان", "e",
		"ساعت", "",
		" ", "", "\t", "T",
	).Replace(timeStr)

	if idx := strings.Index(processed, "e"); idx != -1 {
		processExamInfo(record, processed[idx:])
		processed = processed[:idx]
	}

	parts := strings.Split(processed, "c")
	for i := 1; i < len(parts) && i <= 5; i++ {
		timeSlot := cleanText(parts[i])
		if locIdx := strings.Index(timeSlot, "مکان"); locIdx != -1 {
			timeSlot = timeSlot[:locIdx]
		}

		switch i {
		case 1:
			record.Time1 = timeSlot
		case 2:
			record.Time2 = timeSlot
		case 3:
			record.Time3 = timeSlot
		case 4:
			record.Time4 = timeSlot
		case 5:
			record.Time5 = timeSlot
		}
	}
}

func processExamInfo(record *Record, examStr string) {
	if dateStart := strings.Index(examStr, "("); dateStart != -1 {
		if dateEnd := strings.Index(examStr[dateStart:], ")"); dateEnd != -1 {
			date := examStr[dateStart+1 : dateStart+dateEnd]
			date = strings.ReplaceAll(date, ".", "/")
			record.DateExam = cleanText(date)
		}
	}

	if timeStart := strings.Index(examStr, "):"); timeStart != -1 {
		record.TimeExam = cleanText(examStr[timeStart+2:])
	}
}

func cleanText(text string) string {
	replacer := strings.NewReplacer(
		"۰", "0", "۱", "1", "۲", "2", "۳", "3", "۴", "4",
		"۵", "5", "۶", "6", "۷", "7", "۸", "8", "۹", "9",
		"ي", "ی", "ك", "ک", "ئ", "ی", "ء", "", "٠", "0", "١", "1",
		"٢", "2", "٣", "3", "٤", "4", "٥", "5", "٦", "6", "٧", "7",
		"٨", "8", "٩", "9", "‌", " ", "‏", "", "\u200c", " ",
	)
	return strings.TrimSpace(replacer.Replace(text))
}

func generateJSONOutput(faculty string, records []Record) error {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	jsonPath := filepath.Join("export", faculty+".json")
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing JSON file: %w", err)
	}

	return nil
}

func generateSQLInsert(faculty string, records []Record) (string, error) {
	if len(records) == 0 {
		return "", nil
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("\nINSERT INTO %s VALUES\n", faculty))

	for i, record := range records {
		if i > 0 {
			builder.WriteString(",\n")
		}
		builder.WriteString(fmt.Sprintf(
			"(NULL,'%s','%s',%s,%s,'%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')",
			escapeSQL(record.CourseID),
			escapeSQL(record.Name),
			record.Weight,
			record.Capacity,
			escapeSQL(record.Gender),
			escapeSQL(record.Professor),
			escapeSQL(record.Faculty),
			escapeSQL(record.Time1),
			escapeSQL(record.Time2),
			escapeSQL(record.Time3),
			escapeSQL(record.Time4),
			escapeSQL(record.Time5),
			escapeSQL(record.TimeExam),
			escapeSQL(record.DateExam),
		))
	}

	builder.WriteString(";")
	return builder.String(), nil
}

func escapeSQL(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func writeSQLFile(faculty, content string) error {
	if content == "" {
		return nil
	}

	sqlPath := filepath.Join("export", faculty+".sql")
	return os.WriteFile(sqlPath, []byte(content), 0644)
}
