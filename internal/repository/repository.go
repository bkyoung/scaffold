package repository

import "io"

// Creater defines methods for creating a source code repository
type Creater interface {
    Create() error
}

// Cloner defines methods for checking out a source code repository
type Cloner interface {
    Clone(io.Writer) error
}

// SCMRepository defines behavior that a source code repository needs to possess
type SCMRepository interface {
    Creater
    Cloner
    URL() (string, error)
}