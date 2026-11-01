package filecontent

import (
	"bytes"
	"log/slog"
	"sync"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

var configureOnce sync.Once

var pdfStrippedEntryKeys = []string{"Metadata", "PieceInfo"}

func stripPDF(content []byte, report AnonymizationReport) ([]byte, AnonymizationReport, error) {
	configureOnce.Do(api.DisableConfigDir)

	cleaned, stripErr := rewritePDFWithoutMetadata(content)
	if stripErr != nil {
		slog.Warn("pdf metadata removal failed; original pdf kept", "error", stripErr)
		report.AddLimitation(LimitationPDFMetadata)
		return content, report, nil
	}

	report.AddStep(StepMetadataRemoved)
	return cleaned, report, nil
}

func rewritePDFWithoutMetadata(content []byte) ([]byte, error) {
	configuration := model.NewDefaultConfiguration()
	configuration.ValidationMode = model.ValidationRelaxed

	pdfContext, readErr := api.ReadValidateAndOptimize(bytes.NewReader(content), configuration)
	if readErr != nil {
		return nil, readErr
	}

	pdfContext.Info = nil
	pdfContext.Title = ""
	pdfContext.Subject = ""
	pdfContext.Author = ""
	pdfContext.Creator = ""
	pdfContext.Producer = ""
	pdfContext.Keywords = ""
	pdfContext.KeywordList = nil
	pdfContext.Properties = nil
	pdfContext.CatalogXMPMeta = nil

	for _, tableEntry := range pdfContext.Table {
		if tableEntry == nil || tableEntry.Free {
			continue
		}
		switch object := tableEntry.Object.(type) {
		case types.Dict:
			deleteMetadataEntries(object)
		case types.StreamDict:
			deleteMetadataEntries(object.Dict)
		}
	}

	var output bytes.Buffer
	if writeErr := api.Write(pdfContext, &output, configuration); writeErr != nil {
		return nil, writeErr
	}
	return output.Bytes(), nil
}

func deleteMetadataEntries(dictionary types.Dict) {
	for _, entryKey := range pdfStrippedEntryKeys {
		dictionary.Delete(entryKey)
	}
}
