package utils

type stream struct {
	label  string
	values []float64
}

// SourceData structures data relevant for mgJSON files
type SourceData struct {
	timing  []float64
	streams []stream
}

// For now, only the fields we are using are specified
type mgjson struct {
	Version                string `json:"version"`
	Creator                string `json:"creator"`
	DynamicSamplesPresentB bool   `json:"dynamicSamplesPresentB"`
	DynamicDataInfo        struct {
		UseTimecodeB bool `json:"useTimecodeB"`
		UtcInfo      struct {
			PrecisionLength int  `json:"precisionLength"`
			IsGMT           bool `json:"isGMT"`
		} `json:"utcInfo"`
	} `json:"dynamicDataInfo"`
	DataOutline []struct {
		ObjectType  string `json:"objectType"`
		DisplayName string `json:"displayName"`
		DampleSetID string `json:"sampleSetID"`
		DataType    struct {
			Type                   string `json:"type"`
			NumberStringProperties struct {
				Pattern struct {
					DigitsInteger int `json:"digitsInteger"`
					DigitsDecimal int `json:"digitsDecimal"`
				} `json:"pattern"`
				Range struct {
					Occuring struct {
						Min int `json:"min"`
						Max int `json:"max"`
					} `json:"occuring"`
					Legal struct {
						Min int `json:"min"`
						Max int `json:"max"`
					} `json:"legal"`
				} `json:"range"`
			} `json:"numberStringProperties"`
		} `json:"dataType"`
		Interpolation         string `json:"interpolation"`
		HasExpectedFrequencyB bool   `json:"hasExpectedFrequecyB"`
		SampleCount           int    `json:"sampleCount"`
		MatchName             string `json:"matchName"`
	} `json:"dataOutline"`
	DataDynamicSamples []struct {
		SampleSetID string `json:"sampleSetID"`
		Samples     []struct {
			Time  string `json:"time"`
			Value string `json:"value"`
		} `json:"samples"`
	} `json:"dataDynamicSamples"`
}

func FormatMgjson(sd SourceData, creator string) mgjson {

	//Hardcode non configurable values (for now)
	data := mgjson{
		Version:                "MGJSON2.0.0",
		Creator:                creator,
		DynamicSamplesPresentB: true,
	}

	data.DynamicDataInfo.UseTimecodeB = false
	data.DynamicDataInfo.UtcInfo.PrecisionLength = 3
	data.DynamicDataInfo.UtcInfo.IsGMT = true

	return data
}
