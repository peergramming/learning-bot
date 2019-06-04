package models

import (
	"hash"
)

type Report struct {
	User          string
	Repository    string
	Commit        string
	ArchiveSHA256 hash.Hash
}
