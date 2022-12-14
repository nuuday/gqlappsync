package generator

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/nuuday/gqlappsync/generator/test"
)

func TestGoGenerator(t *testing.T) {
	// TODO Add logic that creates the gqlgen.yml & schema.graphql files instead of referencing them at a static directory.
	// Then only run setupSuite if files are missing
	err := Run("test/gqlgen.yml")
	if err != nil {
		t.Error(err)
	}
}

func TestInterfaceImplementation(t *testing.T) {
	var book test.Book = test.TextBook{
		Title: "Clean Code",
		Author: &test.Author{
			Name: "Robert Cecil Martin",
		},
	}
	book = test.SetTypenameRecursively(book)
	if hasAnyEmptyTypename(book) {
		t.Fail()
	}
}
func TestInterfaceImplementationWithUnion(t *testing.T) {
	var book test.Book = test.TextBook{
		Title: "Clean Code",
		Author: &test.Author{
			Name: "Robert Cecil Martin",
		},
		SupplementaryMaterial: []test.MediaItem{test.AudioClip{
			Duration: 120,
		}},
	}

	book = test.SetTypenameRecursively(book)
	if hasAnyEmptyTypename(book) {
		t.Fail()
	}
}
func TestNilInterfaceImplementation(t *testing.T) {
	var book test.Book = nil
	book = test.SetTypenameRecursively(book)
	if hasAnyEmptyTypename(book) {
		t.Fail()
	}
}

func TestStructWithFieldThatIsAnInterfaceImplementation(t *testing.T) {
	var library test.Library = test.Library{
		Books: []test.Book{
			test.TextBook{
				Title: "Clean Code",
				Author: &test.Author{
					Name: "Robert Cecil Martin",
				},
			},
		},
	}

	library = test.SetTypenameRecursively(library)
	if hasAnyEmptyTypename(library) {
		t.Fail()
	}
}
func TestStructWithFieldThatIsAPointerInterfaceImplementation(t *testing.T) {
	var library test.Library = test.Library{
		Books: []test.Book{
			&test.TextBook{
				Title: "Clean Code",
				Author: &test.Author{
					Name: "Robert Cecil Martin",
				},
			},
		},
	}

	library = test.SetTypenameRecursively(library)
	if hasAnyEmptyTypename(library) {
		t.Fail()
	}
}
func hasAnyEmptyTypename[T any](x T) bool {
	switch val := any(x).(type) {
	case []any:
		for _, element := range val {
			if hasAnyEmptyTypename(element) {
				return true
			}
		}
		return false
	default:
		bytes, _ := json.Marshal(x)
		return strings.Contains(string(bytes), "\"__typename\":\"\"") // ! This is a suboptimal way to check if typename is set. Find a better alternative. Maybe recursive reflection?
	}
}
