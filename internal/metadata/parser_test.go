package metadata

import (
	"reflect"
	"strings"
	"testing"
)

type testUser struct {
	ID       string `firego:"id" firestore:"-"`
	Name     string `firestore:"name,omitempty"`
	Age      int
	Ignored  string `firestore:"-"`
	internal string
}

func TestParse(t *testing.T) {
	got, err := Parse[*testUser]("users")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if got.Name != "testUser" {
		t.Errorf("Name = %q, want %q", got.Name, "testUser")
	}
	if got.Collection != "users" {
		t.Errorf("Collection = %q, want %q", got.Collection, "users")
	}
	if got.GoType != reflect.TypeOf(testUser{}) {
		t.Errorf("GoType = %v, want %v", got.GoType, reflect.TypeOf(testUser{}))
	}
	if len(got.Fields) != 3 {
		t.Fatalf("len(Fields) = %d, want 3", len(got.Fields))
	}
	if got.IDField == nil || got.IDField.Name != "ID" {
		t.Fatalf("IDField = %#v, want ID field", got.IDField)
	}
	if got.Fields[1].FirestoreName != "name" {
		t.Errorf("Name FirestoreName = %q, want %q", got.Fields[1].FirestoreName, "name")
	}
	if got.Fields[2].FirestoreName != "Age" {
		t.Errorf("Age FirestoreName = %q, want %q", got.Fields[2].FirestoreName, "Age")
	}
}

func TestParseRejectsInvalidModels(t *testing.T) {
	tests := []struct {
		name    string
		parse   func() error
		wantErr string
	}{
		{
			name: "empty collection",
			parse: func() error {
				_, err := Parse[testUser]("")
				return err
			},
			wantErr: "collection must not be empty",
		},
		{
			name: "non-struct model",
			parse: func() error {
				_, err := Parse[string]("values")
				return err
			},
			wantErr: "model must be a struct",
		},
		{
			name: "multiple ID fields",
			parse: func() error {
				type model struct {
					First  string `firego:"id"`
					Second string `firego:"id"`
				}
				_, err := Parse[model]("models")
				return err
			},
			wantErr: "multiple ID fields",
		},
		{
			name: "non-string ID",
			parse: func() error {
				type model struct {
					ID int `firego:"id"`
				}
				_, err := Parse[model]("models")
				return err
			},
			wantErr: "must be a string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.parse()
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want error containing %q", err, tt.wantErr)
			}
		})
	}
}
