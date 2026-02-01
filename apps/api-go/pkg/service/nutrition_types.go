package service

import "log/slog"

type ScanInput struct {
	ImageBase64 string  `json:"image_base64"`
	Description *string `json:"description,omitempty"`
}

// LogValue implements slog.LogValuer for structured logging
func (s *ScanInput) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.Int("image_size_bytes", len(s.ImageBase64)),
	}
	if s.Description != nil {
		attrs = append(attrs, slog.String("description", *s.Description))
	}
	return slog.GroupValue(attrs...)
}

type ScanOutput struct {
	FoodName    string     `json:"food_name"`
	Confidence  float64    `json:"confidence"`
	Macros      *MacroData `json:"macros"`
	ServingSize string     `json:"serving_size"`
}

type MacroData struct {
	Calories int     `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
}
