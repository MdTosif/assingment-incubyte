package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tofiquem/assingment/pkg/database"
	"github.com/tofiquem/assingment/pkg/models"
)

type SeedData struct {
	FirstNames  []string
	LastNames   []string
	JobTitles   []string
	Countries   []string
	Departments []string
}

func main() {
	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	// Load seed data
	seedData, err := loadSeedData()
	if err != nil {
		log.Fatalf("Failed to load seed data: %v", err)
	}

	// Generate and insert 10,000 employees
	log.Println("Starting to seed 10,000 employees...")
	start := time.Now()

	employees := make([]*models.Employee, 0, 10000)
	for i := 0; i < 10000; i++ {
		employee := generateEmployee(seedData, i)
		employees = append(employees, employee)

		// Batch insert every 1000 employees for performance
		if (i+1)%1000 == 0 {
			if err := database.DB.CreateInBatches(employees, 1000).Error; err != nil {
				log.Printf("Failed to insert batch %d: %v", (i+1)/1000, err)
			} else {
				log.Printf("Inserted %d employees", i+1)
			}
			employees = employees[:0] // Clear the slice
		}
	}

	// Insert any remaining employees
	if len(employees) > 0 {
		if err := database.DB.CreateInBatches(employees, len(employees)).Error; err != nil {
			log.Fatalf("Failed to insert final batch: %v", err)
		}
	}

	duration := time.Since(start)
	log.Printf("Successfully seeded 10,000 employees in %v", duration)
}

func loadSeedData() (*SeedData, error) {
	data := &SeedData{}

	// Get seed data directory from environment variable or use default
	seedDir := os.Getenv("SEED_DATA_DIR")
	if seedDir == "" {
		seedDir = "seed"
	}

	// Load first names
	firstNames, err := readFileLines(filepath.Join(seedDir, "first_names.txt"))
	if err != nil {
		return nil, fmt.Errorf("failed to load first names: %v", err)
	}
	data.FirstNames = firstNames

	// Load last names
	lastNames, err := readFileLines(filepath.Join(seedDir, "last_names.txt"))
	if err != nil {
		return nil, fmt.Errorf("failed to load last names: %v", err)
	}
	data.LastNames = lastNames

	// Define job titles
	data.JobTitles = []string{
		"Software Engineer",
		"Senior Software Engineer",
		"Lead Software Engineer",
		"Principal Software Engineer",
		"Software Architect",
		"DevOps Engineer",
		"Senior DevOps Engineer",
		"Site Reliability Engineer",
		"Product Manager",
		"Senior Product Manager",
		"UX Designer",
		"Senior UX Designer",
		"UI Designer",
		"Frontend Developer",
		"Senior Frontend Developer",
		"Backend Developer",
		"Senior Backend Developer",
		"Full Stack Developer",
		"Senior Full Stack Developer",
		"Data Scientist",
		"Senior Data Scientist",
		"Machine Learning Engineer",
		"QA Engineer",
		"Senior QA Engineer",
		"Technical Writer",
		"Project Manager",
		"Scrum Master",
		"Business Analyst",
		"Systems Administrator",
		"Network Engineer",
		"Security Engineer",
		"Cloud Engineer",
		"Database Administrator",
		"Mobile Developer",
		"Senior Mobile Developer",
		"Engineering Manager",
		"Director of Engineering",
		"CTO",
		"VP of Engineering",
		"HR Manager",
		"Senior HR Manager",
		"Recruiter",
		"Senior Recruiter",
		"Marketing Manager",
		"Senior Marketing Manager",
		"Sales Manager",
		"Senior Sales Manager",
		"Account Manager",
		"Senior Account Manager",
		"Financial Analyst",
		"Senior Financial Analyst",
		"Operations Manager",
		"Senior Operations Manager",
		"Customer Success Manager",
		"Technical Support Engineer",
		"Senior Technical Support Engineer",
	}

	// Define countries
	data.Countries = []string{
		"United States",
		"United Kingdom",
		"Canada",
		"Germany",
		"France",
		"India",
		"Japan",
		"China",
		"Australia",
		"Netherlands",
		"Sweden",
		"Norway",
		"Denmark",
		"Finland",
		"Switzerland",
		"Austria",
		"Belgium",
		"Ireland",
		"Spain",
		"Italy",
		"Poland",
		"Czech Republic",
		"Hungary",
		"Romania",
		"Bulgaria",
		"Greece",
		"Portugal",
		"Turkey",
		"Israel",
		"UAE",
		"Saudi Arabia",
		"South Africa",
		"Brazil",
		"Argentina",
		"Mexico",
		"Chile",
		"Colombia",
		"Peru",
		"South Korea",
		"Singapore",
		"Malaysia",
		"Thailand",
		"Indonesia",
		"Philippines",
		"Vietnam",
		"New Zealand",
		"Russia",
		"Ukraine",
		"Egypt",
		"Nigeria",
		"Kenya",
		"Ghana",
		"Morocco",
	}

	// Define departments
	data.Departments = []string{
		"Engineering",
		"Product",
		"Design",
		"Marketing",
		"Sales",
		"HR",
		"Finance",
		"Operations",
		"Customer Success",
		"Legal",
		"IT",
		"Data Science",
		"Security",
		"Infrastructure",
		"Mobile",
		"Web",
		"Backend",
		"Frontend",
		"DevOps",
		"QA",
	}

	return data, nil
}

func readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}

func generateEmployee(data *SeedData, index int) *models.Employee {
	rand.Seed(time.Now().UnixNano() + int64(index))

	firstName := data.FirstNames[rand.Intn(len(data.FirstNames))]
	lastName := data.LastNames[rand.Intn(len(data.LastNames))]
	jobTitle := data.JobTitles[rand.Intn(len(data.JobTitles))]
	country := data.Countries[rand.Intn(len(data.Countries))]
	department := data.Departments[rand.Intn(len(data.Departments))]

	// Generate realistic salary based on job title and country
	baseSalary := getBaseSalary(jobTitle)
	countryMultiplier := getCountryMultiplier(country)
	salary := baseSalary * countryMultiplier * (0.8 + rand.Float64()*0.4) // ±20% variation

	// Generate hire date within the last 5 years
	hireDate := time.Now().AddDate(-rand.Intn(5), -rand.Intn(12), -rand.Intn(30))

	return &models.Employee{
		FirstName:  firstName,
		LastName:   lastName,
		Email:      fmt.Sprintf("%s.%s%d@company.com", strings.ToLower(firstName), strings.ToLower(lastName), index),
		JobTitle:   jobTitle,
		Country:    country,
		Salary:     salary,
		Department: department,
		HireDate:   hireDate,
	}
}

func getBaseSalary(jobTitle string) float64 {
	salaryRanges := map[string]float64{
		"Software Engineer":                 80000,
		"Senior Software Engineer":          120000,
		"Lead Software Engineer":            140000,
		"Principal Software Engineer":       160000,
		"Software Architect":                150000,
		"DevOps Engineer":                   90000,
		"Senior DevOps Engineer":            130000,
		"Site Reliability Engineer":         120000,
		"Product Manager":                   110000,
		"Senior Product Manager":            140000,
		"UX Designer":                       85000,
		"Senior UX Designer":                110000,
		"UI Designer":                       75000,
		"Frontend Developer":                80000,
		"Senior Frontend Developer":         115000,
		"Backend Developer":                 85000,
		"Senior Backend Developer":          120000,
		"Full Stack Developer":              90000,
		"Senior Full Stack Developer":       125000,
		"Data Scientist":                    120000,
		"Senior Data Scientist":             150000,
		"Machine Learning Engineer":         130000,
		"QA Engineer":                       70000,
		"Senior QA Engineer":                95000,
		"Technical Writer":                  75000,
		"Project Manager":                   95000,
		"Scrum Master":                      85000,
		"Business Analyst":                  80000,
		"Systems Administrator":             75000,
		"Network Engineer":                  80000,
		"Security Engineer":                 100000,
		"Cloud Engineer":                    95000,
		"Database Administrator":            85000,
		"Mobile Developer":                  85000,
		"Senior Mobile Developer":           115000,
		"Engineering Manager":               140000,
		"Director of Engineering":           180000,
		"CTO":                               250000,
		"VP of Engineering":                 200000,
		"HR Manager":                        85000,
		"Senior HR Manager":                 110000,
		"Recruiter":                         65000,
		"Senior Recruiter":                  85000,
		"Marketing Manager":                 90000,
		"Senior Marketing Manager":          120000,
		"Sales Manager":                     95000,
		"Senior Sales Manager":              125000,
		"Account Manager":                   80000,
		"Senior Account Manager":            105000,
		"Financial Analyst":                 85000,
		"Senior Financial Analyst":          110000,
		"Operations Manager":                90000,
		"Senior Operations Manager":         115000,
		"Customer Success Manager":          80000,
		"Technical Support Engineer":        70000,
		"Senior Technical Support Engineer": 90000,
	}

	if salary, exists := salaryRanges[jobTitle]; exists {
		return salary
	}

	return 80000 // Default salary
}

func getCountryMultiplier(country string) float64 {
	multipliers := map[string]float64{
		"United States":  1.0,
		"United Kingdom": 0.8,
		"Canada":         0.85,
		"Germany":        0.9,
		"France":         0.85,
		"India":          0.3,
		"Japan":          0.9,
		"China":          0.4,
		"Australia":      0.95,
		"Netherlands":    0.9,
		"Sweden":         0.95,
		"Norway":         1.1,
		"Denmark":        1.0,
		"Finland":        0.95,
		"Switzerland":    1.2,
		"Austria":        0.9,
		"Belgium":        0.85,
		"Ireland":        0.9,
		"Spain":          0.7,
		"Italy":          0.75,
		"Poland":         0.5,
		"Czech Republic": 0.55,
		"Hungary":        0.5,
		"Romania":        0.45,
		"Bulgaria":       0.4,
		"Greece":         0.6,
		"Portugal":       0.65,
		"Turkey":         0.35,
		"Israel":         0.85,
		"UAE":            0.8,
		"Saudi Arabia":   0.7,
		"South Africa":   0.4,
		"Brazil":         0.5,
		"Argentina":      0.45,
		"Mexico":         0.4,
		"Chile":          0.5,
		"Colombia":       0.35,
		"Peru":           0.3,
		"South Korea":    0.8,
		"Singapore":      0.85,
		"Malaysia":       0.4,
		"Thailand":       0.3,
		"Indonesia":      0.25,
		"Philippines":    0.25,
		"Vietnam":        0.2,
		"New Zealand":    0.85,
		"Russia":         0.4,
		"Ukraine":        0.3,
		"Egypt":          0.2,
		"Nigeria":        0.15,
		"Kenya":          0.15,
		"Ghana":          0.12,
		"Morocco":        0.2,
	}

	if multiplier, exists := multipliers[country]; exists {
		return multiplier
	}

	return 0.8 // Default multiplier
}
