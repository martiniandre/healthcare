package validator

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	emailRegex      = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneBRRegex    = regexp.MustCompile(`^\(\d{2}\) \d{4,5}-\d{4}$`)
	phoneE164Regex  = regexp.MustCompile(`^\+55\d{10,11}$`)
	icd10Regex      = regexp.MustCompile(`^[A-Z]\d{2}(\.\d{1,2})?$`)
	loincRegex      = regexp.MustCompile(`^\d{3,6}-\d$`)
	cpfRegex        = regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
	crmRegex        = regexp.MustCompile(`^(CRM|COREN)(-[A-Z]{2})?[\s\-]?\d{1,6}$`)
)

var validDICOMModalities = map[string]bool{
	"CR": true, "CT": true, "MR": true, "US": true, "DX": true,
	"MG": true, "NM": true, "PT": true, "XA": true, "RF": true,
	"OT": true, "SC": true, "ES": true,
}

var validClinicalStatuses = map[string]bool{
	"active": true, "inactive": true, "resolved": true,
}

var validGenders = map[string]bool{
	"male": true, "female": true, "other": true, "unknown": true,
}

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

func IsValidPhoneBR(phone string) bool {
	return phoneBRRegex.MatchString(phone)
}

func IsValidPhoneE164(phone string) bool {
	return phoneE164Regex.MatchString(phone)
}

func IsValidPhone(phone string) bool {
	return IsValidPhoneBR(phone) || IsValidPhoneE164(phone)
}

func IsValidICD10(code string) bool {
	return icd10Regex.MatchString(strings.ToUpper(strings.TrimSpace(code)))
}

func IsValidLOINC(code string) bool {
	return loincRegex.MatchString(strings.TrimSpace(code))
}

func IsValidDICOMModality(modality string) bool {
	return validDICOMModalities[strings.ToUpper(strings.TrimSpace(modality))]
}

func IsValidClinicalStatus(status string) bool {
	return validClinicalStatuses[strings.ToLower(strings.TrimSpace(status))]
}

func IsValidGender(gender string) bool {
	return validGenders[strings.ToLower(strings.TrimSpace(gender))]
}

func IsValidCRMNumber(license string) bool {
	if strings.TrimSpace(license) == "" {
		return true
	}
	return crmRegex.MatchString(strings.TrimSpace(license))
}

func IsValidCPF(cpf string) bool {
	digitsOnly := regexp.MustCompile(`\D`).ReplaceAllString(cpf, "")
	if len(digitsOnly) != 11 {
		return false
	}
	allSame := true
	for _, character := range digitsOnly {
		if character != rune(digitsOnly[0]) {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}
	firstDigitSum := 0
	for position := 0; position < 9; position++ {
		digit, _ := strconv.Atoi(string(digitsOnly[position]))
		firstDigitSum += digit * (10 - position)
	}
	firstRemainder := (firstDigitSum * 10) % 11
	if firstRemainder == 10 || firstRemainder == 11 {
		firstRemainder = 0
	}
	firstVerifier, _ := strconv.Atoi(string(digitsOnly[9]))
	if firstRemainder != firstVerifier {
		return false
	}
	secondDigitSum := 0
	for position := 0; position < 10; position++ {
		digit, _ := strconv.Atoi(string(digitsOnly[position]))
		secondDigitSum += digit * (11 - position)
	}
	secondRemainder := (secondDigitSum * 10) % 11
	if secondRemainder == 10 || secondRemainder == 11 {
		secondRemainder = 0
	}
	secondVerifier, _ := strconv.Atoi(string(digitsOnly[10]))
	return secondRemainder == secondVerifier
}

func IsPastDate(dateStr string) bool {
	parsedDate, parseErr := time.Parse("2006-01-02", dateStr)
	if parseErr != nil {
		return false
	}
	return parsedDate.Before(time.Now())
}

func IsValidObservationRange(loincCode string, value float64) bool {
	ranges := map[string][2]float64{
		"8867-4":  {0, 300},
		"8310-5":  {30.0, 45.0},
		"55284-4": {0, 300},
		"59408-5": {0, 100},
		"9279-1":  {0, 60},
		"29463-7": {0, 500},
		"8302-2":  {0, 250},
	}
	boundaries, exists := ranges[loincCode]
	if !exists {
		return value > 0
	}
	return value >= boundaries[0] && value <= boundaries[1]
}

var _ = cpfRegex
