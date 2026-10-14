package models

import utilsModels "github.com/vucongthanh92/go-base-utils/models"

// ErrorDTO aliases the shared Courier API error contract. Keeping the domain
// alias avoids leaking a package-level dependency into payment business DTOs
// while ensuring handlers, middleware and use cases produce one response shape.
type ErrorDTO = utilsModels.ErrorDTO
