// Package llms provides interfaces and types for interacting with
// large language models (LLMs) from various providers.
package llms

import "context"

// ContentType represents the type of content in a message part.
type ContentType string

const (
	// ContentTypeText represents plain text content.
	ContentTypeText ContentType = "text"
	// ContentTypeImageURL represents an image referenced by URL.
	ContentTypeImageURL ContentType = "image_url"
	// ContentTypeBinary represents raw binary content.
	ContentTypeBinary ContentType = "binary"
)

// TextContent holds a text message part.
type TextContent struct {
	Text string `json:"text"`
}

// ImageURLContent holds an image URL message part.
type ImageURLContent struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// BinaryContent holds raw binary data with a MIME type.
type BinaryContent struct {
	MIMEType string `json:"mime_type"`
	Data     []byte `json:"data"`
}

// ContentPart represents a single part of a message, which can be
// text, an image URL, or binary data.
type ContentPart interface {
	contentPartType() ContentType
}

func (t TextContent) contentPartType() ContentType     { return ContentTypeText }
func (i ImageURLContent) contentPartType() ContentType { return ContentTypeImageURL }
func (b BinaryContent) contentPartType() ContentType   { return ContentTypeBinary }

// ChatMessageType defines the role of a message in a conversation.
type ChatMessageType string

const (
	// ChatMessageTypeHuman represents a message from the human/user.
	ChatMessageTypeHuman ChatMessageType = "human"
	// ChatMessageTypeAI represents a message from the AI assistant.
	ChatMessageTypeAI ChatMessageType = "ai"
	// ChatMessageTypeSystem represents a system-level instruction message.
	ChatMessageTypeSystem ChatMessageType = "system"
	// ChatMessageTypeGeneric represents a message with an arbitrary role.
	// Note: prefer using a specific type when possible; generic is a fallback.
	ChatMessageTypeGeneric ChatMessageType = "generic"
	// ChatMessageTypeFunction represents a function call result message.
	ChatMessageTypeFunction ChatMessageType = "function"
	// ChatMessageTypeTool represents a tool call result message.
	ChatMessageTypeTool ChatMessageType = "tool"
)

// MessageContent holds a chat message with a role and one or more content parts.
type MessageContent struct {
	Role  ChatMessageType `json:"role"`
	Parts []ContentPart   `json:"parts"`
}

// ContentResponse is the response returned by a model for a content generation request.
type ContentResponse struct {
	// Choices contains the generated response candidates.
	Choices []*ContentChoice `json:"choices"`
}

// ContentChoice represents a single generated response candidate.
type ContentChoice struct {
	// Content is the generated text.
	Content string `json:"content"`
	// StopReason describes why the model stopped generating.
	// Common values: "stop", "length", "content_filter", "tool_calls".
	// Note: not all providers populate this field; check provider docs for details.
	// Some providers (e.g. Anthropic) use "end_turn" instead of "stop".
	StopReason string `json:"stop_reason,omitempty"`
}

// Model is the interface all LLM implementations must satisfy.
type Model interface {
	// Call is a simplified interface for single-turn text generation.
	// Deprecated: prefer GenerateContent for multi-turn and multi-modal use.
	Call(ctx context.Context, prompt string, options ...CallOption) (string, error)
	// GenerateContent generates a response for the given messages.
	GenerateContent(ctx context.Context, messages []MessageContent, options ...CallOption) (*ContentResponse, error)
}
