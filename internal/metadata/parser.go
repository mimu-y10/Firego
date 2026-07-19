package metadata

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mimu-y10/firego/schema"
)

// Parse creates schema metadata for model T and the given collection.
func Parse[T any](collection string) (*schema.Schema, error) {
	if collection == "" {
		return nil, fmt.Errorf("metadata: collection must not be empty")
	}

	modelType := reflect.TypeOf((*T)(nil)).Elem()
	for modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("metadata: model must be a struct")
	}

	// Preallocate enough capacity for every struct field while keeping the
	// length at zero, because unexported and firestore:"-" fields are skipped
	// and only mapped fields should be appended to the result.
	fields := make([]schema.Field, 0, modelType.NumField())
	idFieldIndex := -1

	for i := 0; i < modelType.NumField(); i++ {
		structField := modelType.Field(i)

		if !structField.IsExported() {
			continue
		}

		isID := structField.Tag.Get("firego") == "id"
		firestoreName := parseFirestoreName(structField)

		// A field ignored by the Firestore codec is only relevant when it holds
		// the document ID.
		if firestoreName == "-" && !isID {
			continue
		}

		if isID {
			if idFieldIndex >= 0 {
				return nil, fmt.Errorf("metadata: multiple ID fields")
			}

			if structField.Type.Kind() != reflect.String {
				return nil, fmt.Errorf("metadata: ID field %q must be a string", structField.Name)
			}

			idFieldIndex = len(fields)
		}

		fields = append(fields, schema.Field{
			Name:          structField.Name,
			FirestoreName: firestoreName,
			GoType:        structField.Type,
			StructIndex:   structField.Index,
			IsID:          isID,
		})
	}

	result := &schema.Schema{
		Name:       modelType.Name(),
		Collection: collection,
		GoType:     modelType,
		Fields:     fields,
	}

	if idFieldIndex >= 0 {
		result.IDField = &result.Fields[idFieldIndex]
	}

	return result, nil
}

func parseFirestoreName(field reflect.StructField) string {
	tag := field.Tag.Get("firestore")
	if tag == "" {
		return field.Name
	}

	name, _, _ := strings.Cut(tag, ",")
	if name == "" {
		return field.Name
	}

	return name
}
