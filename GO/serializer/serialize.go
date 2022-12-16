package serializer

import (
	// import protobuf package
	"fmt"

	proto "github.com/golang/protobuf/proto"

	os "os"
)

func WriteProtobufToJSONFile(message proto.Message, filename string) error {
	// serialize message to JSON
	out, err := ProtobufToJSON(message)

	if err != nil {
		return fmt.Errorf("cannot marshal proto message to JSON: %w", err)
	}

	// write to file
	err = os.WriteFile(filename, []byte(out), 0644)
	if err != nil {
		return fmt.Errorf("cannot write to file: %w", err)
	}

	return nil
}

func WriteProtobufToBinaryFile(message proto.Message, filename string) error {
	// serialize message to binary
	out, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary: %w", err)
	}

	// write to file
	err = os.WriteFile(filename, out, 0644)
	if err != nil {
		return fmt.Errorf("cannot write to file: %w", err)
	}

	return nil
}

func ReadProtobufFromBinaryFile(message proto.Message, filename string) error {
	// read from file
	in, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read from file: %w", err)
	}

	// deserialize message from binary
	err = proto.Unmarshal(in, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal proto message from binary: %w", err)
	}

	return nil
}
