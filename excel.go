package mlb

import "github.com/xuri/excelize/v2"

func conditionalFormattingHigherGreen(sheet string, f *excelize.File, ref string) error {
	return f.SetConditionalFormat(sheet, ref,
		[]excelize.ConditionalFormatOptions{
			{
				Type:     "3_color_scale",
				Criteria: "=",
				MinType:  "min",
				MidType:  "percentile",
				MaxType:  "max",
				MinColor: "#F8696B",
				MidColor: "#FCF55F",
				MaxColor: "#63BE7B",
			},
		},
	)
}

func conditionalFormattingLowerGreen(sheet string, f *excelize.File, ref string) error {
	return f.SetConditionalFormat(sheet, ref,
		[]excelize.ConditionalFormatOptions{
			{
				Type:     "3_color_scale",
				Criteria: "=",
				MinType:  "min",
				MidType:  "percentile",
				MaxType:  "max",
				MaxColor: "#F8696B",
				MidColor: "#FCF55F",
				MinColor: "#63BE7B",
			},
		},
	)
}
