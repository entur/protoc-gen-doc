package extensions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pseudomuto/protoc-gen-doc/extensions"
)

func TestTransformUnwrapsReflectMessages(t *testing.T) {
	extensions.SetTransformer("test.timestamp", func(payload any) any {
		ts, ok := payload.(*timestamppb.Timestamp)
		if !ok {
			return nil
		}
		return ts.AsTime().Format("2006-01-02")
	})

	ts := timestamppb.New(timestamppb.Now().AsTime())

	// Protokit reports message-typed option extensions as protoreflect.Value results, so the
	// payload arrives as a protoreflect.Message rather than as the generated struct pointer that
	// transformers assert on. Both forms must reach the transformer as the generated type.
	require.Equal(
		t,
		map[string]any{"test.timestamp": ts.AsTime().Format("2006-01-02")},
		extensions.Transform(map[string]any{"test.timestamp": ts.ProtoReflect()}),
	)
	require.Equal(
		t,
		map[string]any{"test.timestamp": ts.AsTime().Format("2006-01-02")},
		extensions.Transform(map[string]any{"test.timestamp": ts}),
	)
}
