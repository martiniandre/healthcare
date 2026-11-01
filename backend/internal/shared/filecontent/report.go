package filecontent

const (
	StepMetadataRemoved     = "metadata_removed"
	StepNoMetadataFound     = "no_metadata_found"
	StepTrailingDataRemoved = "trailing_data_removed"
	StepFileNameNeutralized = "filename_neutralized"
	LimitationBurnedInText  = "burned_in_text_not_removed"
	LimitationPDFMetadata   = "pdf_metadata_not_removed"
)

type AnonymizationReport struct {
	Applied     bool     `json:"applied"`
	Steps       []string `json:"steps"`
	Limitations []string `json:"limitations"`
}

func newReport() AnonymizationReport {
	return AnonymizationReport{
		Steps:       []string{},
		Limitations: []string{LimitationBurnedInText},
	}
}

func (report *AnonymizationReport) AddStep(step string) {
	report.Steps = append(report.Steps, step)
}

func (report *AnonymizationReport) AddLimitation(limitation string) {
	report.Limitations = append(report.Limitations, limitation)
}

func (report *AnonymizationReport) Finalize() {
	report.Applied = len(report.Steps) > 0
}
