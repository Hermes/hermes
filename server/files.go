package server 

import (
	"fmt"
	"github.com/secondbit/wendy"
)

const (
	VaultPath = "~/.Hermes/Vault")

type VaultID string

type Files struct {
	Map map[string]map[string][]FilePart
}

type FilePart struct {
	SplitID string
	Path string
}

func NewFilesMap() Files {
	Files := Files{}
	Files.Map = make(map[string]map[string][]FilePart)
	return Files
}

func (m *Files) Insert(pushFile PushJSON) error {
	files := m[pushFile.VaultID]
	if files == nil {		//Already has fileparts
		files := make(map[string][]FilePart)
	}
	parts := files[pushFile.Filename]
	part := FilePart{pushFile.SplitID, VaultPath + pushFile.Filename}
	parts = append(parts, part)

	m[pushFile.VaultID] = files
}

func (m *Files) HasVaultID(VaultID string) bool {
	if m[VaultID] != nil {
		return true
	}
	return false
}

func (m *Files) GetFilesWithVaultID(VaultID string) []string {
	filenames = make([]Filename, 0)

	for key,values := range m {
		if key == VaultID {
			for key, _ := range m[key] {
				filenames = append(filenames, key)
			}
		}
	}
	
	//Found no files with VaultID
	return filenames
}

func (m *Files) HasVaultIDAndFile(VaultID, Filename string) bool {
	if m[VaultID][Filename] != nil {
		return true
	}
	return false
}

func (m *Files) GetFileParts(VaultID, Filename string) []FilePart {
	return m[VaultID][Filename]
}