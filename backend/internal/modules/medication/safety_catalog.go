package medication

import (
	"sort"
	"strings"
	"unicode"
)

const (
	minimumContainsTermLength = 5
	minimumExactTermLength    = 3

	SubjectTypeIngredient = "ingredient"
	SubjectTypeClass      = "class"

	AllergyMatchIngredient = "ingredient"
	AllergyMatchClass      = "class"
)

var accentReplacer = strings.NewReplacer(
	"á", "a", "à", "a", "â", "a", "ã", "a", "ä", "a", "å", "a",
	"é", "e", "è", "e", "ê", "e", "ë", "e",
	"í", "i", "ì", "i", "î", "i", "ï", "i",
	"ó", "o", "ò", "o", "ô", "o", "õ", "o", "ö", "o",
	"ú", "u", "ù", "u", "û", "u", "ü", "u",
	"ç", "c", "ñ", "n",
)

func NormalizeText(rawText string) string {
	lowerText := accentReplacer.Replace(strings.ToLower(rawText))
	var builder strings.Builder
	builder.Grow(len(lowerText))
	for _, character := range lowerText {
		if (character >= 'a' && character <= 'z') || unicode.IsDigit(character) {
			builder.WriteRune(character)
		}
	}
	return builder.String()
}

type CatalogClass struct {
	Code        string
	Name        string
	SearchTerms []string
}

type CatalogEntry struct {
	Code             string
	Name             string
	ActiveIngredient string
	ClassCode        string
	SearchTerms      []string
}

type CatalogInteraction struct {
	SubjectAType   string
	SubjectA       string
	SubjectBType   string
	SubjectB       string
	Severity       string
	DescriptionKey string
}

type DrugProfile struct {
	Code       string
	Name       string
	Ingredient string
	ClassCode  string
	ClassName  string
	Cataloged  bool
}

type AllergenProfile struct {
	Ingredients map[string]bool
	Classes     map[string]bool
}

type termReference struct {
	term  string
	index int
}

type CatalogSnapshot struct {
	classesByCode map[string]CatalogClass
	classList     []CatalogClass
	entries       []CatalogEntry
	entryTerms    []termReference
	classTerms    []termReference
	interactions  []CatalogInteraction
}

func NewCatalogSnapshot(classes []CatalogClass, entries []CatalogEntry, interactions []CatalogInteraction) *CatalogSnapshot {
	snapshot := &CatalogSnapshot{
		classesByCode: make(map[string]CatalogClass, len(classes)),
		classList:     classes,
		entries:       entries,
		interactions:  interactions,
	}
	for classIndex, catalogClass := range classes {
		snapshot.classesByCode[NormalizeText(catalogClass.Code)] = catalogClass
		classTerms := append([]string{catalogClass.Code, catalogClass.Name}, catalogClass.SearchTerms...)
		for _, term := range classTerms {
			snapshot.classTerms = append(snapshot.classTerms, termReference{term: NormalizeText(term), index: classIndex})
		}
	}
	for entryIndex, entry := range entries {
		entryTerms := append([]string{entry.Name, entry.ActiveIngredient}, entry.SearchTerms...)
		for _, term := range entryTerms {
			snapshot.entryTerms = append(snapshot.entryTerms, termReference{term: NormalizeText(term), index: entryIndex})
		}
	}
	return snapshot
}

func termMatches(normalizedText string, term string) bool {
	if len(term) < minimumExactTermLength {
		return false
	}
	if normalizedText == term {
		return true
	}
	return len(term) >= minimumContainsTermLength && strings.Contains(normalizedText, term)
}

func (snapshot *CatalogSnapshot) entryByCode(normalizedCode string) (CatalogEntry, bool) {
	if normalizedCode == "" {
		return CatalogEntry{}, false
	}
	for _, entry := range snapshot.entries {
		if NormalizeText(entry.Code) == normalizedCode {
			return entry, true
		}
	}
	return CatalogEntry{}, false
}

func (snapshot *CatalogSnapshot) bestEntryByText(normalizedText string) (CatalogEntry, bool) {
	bestLength := 0
	bestIndex := -1
	for _, reference := range snapshot.entryTerms {
		if len(reference.term) > bestLength && termMatches(normalizedText, reference.term) {
			bestLength = len(reference.term)
			bestIndex = reference.index
		}
	}
	if bestIndex < 0 {
		return CatalogEntry{}, false
	}
	return snapshot.entries[bestIndex], true
}

func (snapshot *CatalogSnapshot) ResolveMedication(medicationCode string, medicationName string) DrugProfile {
	profile := DrugProfile{Code: medicationCode, Name: medicationName}
	matchedEntry, found := snapshot.entryByCode(NormalizeText(medicationCode))
	if !found {
		matchedEntry, found = snapshot.bestEntryByText(NormalizeText(medicationName))
	}
	if !found {
		snapshot.applyClassFromText(&profile, NormalizeText(medicationName))
		return profile
	}
	profile.Cataloged = true
	profile.Ingredient = NormalizeText(matchedEntry.ActiveIngredient)
	profile.ClassCode = NormalizeText(matchedEntry.ClassCode)
	if matchedClass, hasClass := snapshot.classesByCode[profile.ClassCode]; hasClass {
		profile.ClassName = matchedClass.Name
	}
	return profile
}

func (snapshot *CatalogSnapshot) applyClassFromText(profile *DrugProfile, normalizedText string) {
	bestLength := 0
	for _, reference := range snapshot.classTerms {
		if len(reference.term) > bestLength && termMatches(normalizedText, reference.term) {
			matchedClass := snapshot.classList[reference.index]
			bestLength = len(reference.term)
			profile.ClassCode = NormalizeText(matchedClass.Code)
			profile.ClassName = matchedClass.Name
		}
	}
}

func (snapshot *CatalogSnapshot) ResolveAllergen(allergenCode string, allergenDisplay string) AllergenProfile {
	profile := AllergenProfile{Ingredients: map[string]bool{}, Classes: map[string]bool{}}
	addEntry := func(entry CatalogEntry) {
		profile.Ingredients[NormalizeText(entry.ActiveIngredient)] = true
		if entry.ClassCode != "" {
			profile.Classes[NormalizeText(entry.ClassCode)] = true
		}
	}
	if codeEntry, found := snapshot.entryByCode(NormalizeText(allergenCode)); found {
		addEntry(codeEntry)
	}
	normalizedDisplay := NormalizeText(allergenDisplay)
	if normalizedDisplay == "" {
		return profile
	}
	for _, reference := range snapshot.entryTerms {
		if termMatches(normalizedDisplay, reference.term) {
			addEntry(snapshot.entries[reference.index])
		}
	}
	for _, reference := range snapshot.classTerms {
		if termMatches(normalizedDisplay, reference.term) {
			profile.Classes[NormalizeText(snapshot.classList[reference.index].Code)] = true
		}
	}
	return profile
}

func (snapshot *CatalogSnapshot) MatchAllergy(allergenCode string, allergenDisplay string, medication DrugProfile) string {
	allergen := snapshot.ResolveAllergen(allergenCode, allergenDisplay)
	if medication.Ingredient != "" && allergen.Ingredients[medication.Ingredient] {
		return AllergyMatchIngredient
	}
	if medication.ClassCode != "" && allergen.Classes[medication.ClassCode] {
		return AllergyMatchClass
	}
	normalizedAllergenCode := NormalizeText(allergenCode)
	if normalizedAllergenCode != "" && normalizedAllergenCode == NormalizeText(medication.Code) {
		return AllergyMatchIngredient
	}
	normalizedAllergenDisplay := NormalizeText(allergenDisplay)
	if len(normalizedAllergenDisplay) >= minimumContainsTermLength && strings.Contains(NormalizeText(medication.Name), normalizedAllergenDisplay) {
		return AllergyMatchIngredient
	}
	return ""
}

func subjectMatches(subjectType string, subject string, medication DrugProfile) bool {
	normalizedSubject := NormalizeText(subject)
	switch subjectType {
	case SubjectTypeIngredient:
		return medication.Ingredient != "" && medication.Ingredient == normalizedSubject
	case SubjectTypeClass:
		return medication.ClassCode != "" && medication.ClassCode == normalizedSubject
	default:
		return false
	}
}

func (snapshot *CatalogSnapshot) FindInteractions(requested DrugProfile, active DrugProfile) []CatalogInteraction {
	var matches []CatalogInteraction
	for _, interaction := range snapshot.interactions {
		forwardMatch := subjectMatches(interaction.SubjectAType, interaction.SubjectA, requested) &&
			subjectMatches(interaction.SubjectBType, interaction.SubjectB, active)
		reverseMatch := subjectMatches(interaction.SubjectAType, interaction.SubjectA, active) &&
			subjectMatches(interaction.SubjectBType, interaction.SubjectB, requested)
		if forwardMatch || reverseMatch {
			matches = append(matches, interaction)
		}
	}
	return matches
}

func severityRank(severity string) int {
	switch severity {
	case SeverityContraindicated:
		return 3
	case SeverityMajor:
		return 2
	case SeverityModerate:
		return 1
	default:
		return 0
	}
}

func sortAlerts(alerts []SafetyAlert) {
	sort.SliceStable(alerts, func(firstIndex int, secondIndex int) bool {
		if alerts[firstIndex].Blocking != alerts[secondIndex].Blocking {
			return alerts[firstIndex].Blocking
		}
		return severityRank(alerts[firstIndex].Severity) > severityRank(alerts[secondIndex].Severity)
	})
}
