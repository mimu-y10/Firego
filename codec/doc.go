// Package codec encodes Go models into Firestore documents and decodes
// Firestore documents back into Go models.
package codec

type Codec interface {
	Encode(v any) (map[string]any, error)
	Decode(data map[string]any, dst any) error
}
