---
trigger: glob
glob: apps/api-go/**
---

# Backend Rules (VitalStack Go API)

## Architecture Reference

Follow the layered architecture defined in `apps/api-go/architecture.md`.

## Testing

Every new function or endpoint **must** have a corresponding test in a `_test.go` file. Use interfaces and dependency injection to make code testable.

## Strict Layering

```
HTTP Request → Gin Router → Huma API → Controller → Service
```

- Higher layers may only import lower layers, never vice versa.
- Do NOT skip layers (e.g., controller must not directly access database).
- Do NOT import database drivers into HTTP handlers.

## DTOs (Data Transfer Objects)

- Use separate DTOs per layer — never reuse models from another layer.
- Create conversion functions on DTOs to convert from/to another layer's DTO:
  ```go
  func (dto *ControllerDTO) ToServiceDTO() *ServiceDTO { ... }
  func FromServiceDTO(s *ServiceDTO) *ControllerDTO { ... }
  ```

## Type-Safe REST (Huma)

- Always use Huma's typed approach with clear Go structs.
- **Never** use `map[string]interface{}`.
- Define structs with `doc`, `json`, and `example` tags for a perfect OpenAPI 3.1 spec.
- Handle errors gracefully with meaningful Huma error responses.

## Documentation

Every exported type, struct, and function **must** have a descriptive comment above it:
```go
// NutritionService handles business logic for food analysis and macro estimation.
type NutritionService struct { ... }

// ScanFood analyzes a food image and returns estimated macro nutrients.
func (s *NutritionService) ScanFood(ctx context.Context, req ScanRequest) (*ScanResult, error) { ... }
```

## Linting

Always run `golangci-lint` to ensure the API has no linting issues. Fix all warnings before committing.

## OpenAPI Specification Sync

After **any** controller is added, updated, or deleted:
1. Run `make openapi` to regenerate the OpenAPI spec and TypeScript client.
2. This executes `gen-openapi-spec` (Go spec generation) followed by `gen-api-client` (TypeScript client generation).
3. Verify the TypeScript client in `apps/web` reflects the API changes.

## Architecture Documentation

Update `apps/api-go/architecture.md` whenever meaningful architecture changes are made (new layers, new packages, new patterns).
