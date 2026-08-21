package observation

type vitalSignMetricDefinition struct {
	FieldKey    string
	LoincCode   string
	CodeDisplay string
	ValueUnit   string
}

var vitalSignMetricDefinitions = []vitalSignMetricDefinition{
	{FieldKey: "heart_rate", LoincCode: "8867-4", CodeDisplay: "Frequência Cardíaca", ValueUnit: "bpm"},
	{FieldKey: "body_temperature", LoincCode: "8310-5", CodeDisplay: "Temperatura Corporal", ValueUnit: "°C"},
	{FieldKey: "systolic_blood_pressure", LoincCode: "8480-6", CodeDisplay: "Pressão Arterial Sistólica", ValueUnit: "mmHg"},
	{FieldKey: "diastolic_blood_pressure", LoincCode: "8462-4", CodeDisplay: "Pressão Arterial Diastólica", ValueUnit: "mmHg"},
	{FieldKey: "oxygen_saturation", LoincCode: "59408-5", CodeDisplay: "Saturação de Oxigênio", ValueUnit: "%"},
	{FieldKey: "respiratory_rate", LoincCode: "9279-1", CodeDisplay: "Frequência Respiratória", ValueUnit: "irpm"},
	{FieldKey: "weight_kg", LoincCode: "29463-7", CodeDisplay: "Peso Corporal", ValueUnit: "kg"},
	{FieldKey: "height_cm", LoincCode: "8302-2", CodeDisplay: "Altura Corporal", ValueUnit: "cm"},
}

var VitalSignPanelSize = len(vitalSignMetricDefinitions)
