package server

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// SpecFormat is the format for the OpenAPI specification
type SpecFormat string

const (
	// SpecFormatJSON is the JSON format for the OpenAPI specification
	SpecFormatJSON SpecFormat = "json"
	// SpecFormatYAML is the YAML format for the OpenAPI specification
	SpecFormatYAML SpecFormat = "yaml"
)

func (f SpecFormat) normalize() SpecFormat {
	switch strings.ToLower(string(f)) {
	case "json", ".json":
		return SpecFormatJSON
	case "yaml", ".yaml":
		return SpecFormatYAML
	default:
		return ""
	}
}

// OpenAPI registers the given controllers and generates the OpenAPI specification to the given path
func (s *Server) OpenAPI(path string, format SpecFormat, controllers ...Controller) error {
	s.RegisterAPI(controllers...)
	return s.generateOpenAPISpecFile(path, format)
}

// generateOpenAPISpecFile generates the OpenAPI specification to the given path
func (s *Server) generateOpenAPISpecFile(path string, format SpecFormat) error {
	slog.Info("writing openapi spec", "path", path, "format", format)
	if format == "" {
		// default to yaml
		format = SpecFormatYAML

		// if path has extension, use this instead of the default
		if path != "" {
			specFormatFromFileExtension := SpecFormat(filepath.Ext(path)).normalize()
			if specFormatFromFileExtension != "" {
				format = specFormatFromFileExtension
			}
		}
	}

	var w io.Writer = os.Stdout

	if path != "" {
		//nolint:gosec // we are not using user input for the path
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		//nolint:errcheck // we don't care about the error when writing the spec
		defer f.Close()
		w = f
	}

	return s.writeOpenAPISpec(w, format)
}

// writeOpenAPISpec writes the OpenAPI specification to the given writer
func (s *Server) writeOpenAPISpec(w io.Writer, format SpecFormat) error {
	openapi := s.huma.OpenAPI()

	var (
		b   []byte
		err error
	)

	switch format {
	case SpecFormatJSON:
		b, err = openapi.MarshalJSON()
	case SpecFormatYAML:
		b, err = openapi.YAML()
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to generate OpenAPI spec: %w", err)
	}

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write OpenAPI spec: %w", err)
	}

	return nil
}
