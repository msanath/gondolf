package printer

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JSON is a writer to display any proto object in JSON format.
type JSON interface {
	PrintProtoObject(msg proto.Message) error
	PrintProtoObjects(msgs []proto.Message) error
}

type jsonPrinter struct {
	protojson.MarshalOptions
}

func NewJSONPrinter() JSON {
	return &jsonPrinter{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
	}
}

// PrintProtoObject displays any proto object in JSON format.
func (j *jsonPrinter) PrintProtoObject(msg proto.Message) error {
	b, err := j.Marshal(msg)
	if err != nil {
		return fmt.Errorf("protojson serializer failed: %w", err)
	}
	fmt.Println(string(b))
	return nil
}

// PrintProtoObjects displays a list of any proto objects in JSON format.
func (j *jsonPrinter) PrintProtoObjects(msgs []proto.Message) error {
	if len(msgs) == 0 {
		return fmt.Errorf("no messages found")
	}
	mergedJSONList := []json.RawMessage{}
	for _, msg := range msgs {
		b, err := j.Marshal(msg)
		if err != nil {
			return fmt.Errorf("protojson serializer failed: %w", err)
		}
		mergedJSONList = append(mergedJSONList, json.RawMessage(b))
	}
	mergedJSON, err := json.Marshal(mergedJSONList)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}
	fmt.Println(string(mergedJSON))
	return nil
}
