package git

import (
	"strings"
	"github.com/pkg/errors"
	"sort"
	"path"
)

type ListFile struct{
	Lines []ListFileLine
}

type ListFileLine struct{
	Fields 		[]string
	Type 		string
	Name    	string
	FullName        string
	ObjectId 	string
}

func (this *ListFile) Init(rawString string) error {
	lines := strings.Split(rawString, "\n")

	if len(lines) == 0 {
		return errors.New("Empty string")
	}

	listFileLines := make([]ListFileLine, 0)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		listFileLine := ListFileLine{
			Fields: fields,
		}
		listFileLine.Type = listFileLine.Fields[1]
		listFileLine.FullName = listFileLine.Fields[3]
		listFileLine.Name = path.Base(listFileLine.FullName)
		listFileLine.ObjectId = listFileLine.Fields[2]
		listFileLines = append(listFileLines, listFileLine)
	}

	sort.Sort(ByFileType(listFileLines))
	this.Lines = listFileLines
	return nil
}

/**
 Sort By File Type
 */
type ByFileType []ListFileLine

func (a ByFileType) Len() int {
	return len(a)
}
func (a ByFileType) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByFileType) Less(i, j int) bool {
	if a[i].Type == "blob" {
		return false
	}
	return true
}