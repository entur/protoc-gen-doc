// Package extensions implements a system for working with extended options.
package extensions

import "google.golang.org/protobuf/reflect/protoreflect"

// Transformer functions for transforming payloads of an extension option into
// something that can be rendered by a template.
type Transformer func(payload any) any

var transformers = make(map[string]Transformer)

// SetTransformer sets the transformer function for the given extension name
func SetTransformer(extensionName string, f Transformer) {
	transformers[extensionName] = f
}

// Transform the extensions using the registered transformers.
func Transform(extensions map[string]any) map[string]any {
	if extensions == nil {
		return nil
	}
	out := make(map[string]any, len(extensions))
	for name, payload := range extensions {
		transform, ok := transformers[name]
		if !ok {
			// No transformer registered, skip.
			continue
		}
		transformedPayload := transform(concrete(payload))
		if transformedPayload == nil {
			// Transformer returned nothing, skip.
			continue
		}
		out[name] = transformedPayload
	}
	return out
}

// concrete unwraps a message-typed extension payload into the generated Go type transformers
// assert on. Protokit reports option extensions as protoreflect.Value results, so message
// payloads arrive as a protoreflect.Message rather than as the generated struct pointer;
// Interface() hands back the latter when the extension type is registered (and a
// *dynamicpb.Message when it is only known from the descriptor set, which transformers ignore).
func concrete(payload any) any {
	if m, ok := payload.(protoreflect.Message); ok {
		return m.Interface()
	}
	return payload
}
