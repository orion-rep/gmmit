package gemini

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/stretchr/testify/assert"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r) //nolint:errcheck
	os.Stdout = old
	return buf.String()
}

func makeResponse(text string) *genai.GenerateContentResponse {
	return &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{genai.Text(text)},
				},
			},
		},
	}
}

func TestModelResponseToString_SingleCandidate(t *testing.T) {
	resp := makeResponse("hello world")
	result := ModelResponseToString(resp)
	assert.Contains(t, result, "hello world")
}

func TestModelResponseToString_TrailingNewline(t *testing.T) {
	resp := makeResponse("no trailing newline")
	result := ModelResponseToString(resp)
	assert.NotEqual(t, '\n', result[len(result)-1])
}

func TestModelResponseToString_MarkdownFences(t *testing.T) {
	resp := makeResponse("```json\n{\"key\":\"value\"}\n```")
	result := ModelResponseToString(resp)
	assert.NotContains(t, result, "```json")
}

func TestModelResponseToString_NilContent(t *testing.T) {
	resp := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{Content: nil},
		},
	}
	result := ModelResponseToString(resp)
	assert.Equal(t, "", result)
}

func TestModelResponseToString_MultipleParts(t *testing.T) {
	resp := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text("part one"),
						genai.Text("part two"),
					},
				},
			},
		},
	}
	result := ModelResponseToString(resp)
	assert.Contains(t, result, "part one")
	assert.Contains(t, result, "part two")
}

func TestPrintModelResponse_NilContent_NoPanic(t *testing.T) {
	resp := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{Content: nil},
		},
	}
	assert.NotPanics(t, func() {
		captureStdout(func() { PrintModelResponse(resp) })
	})
}

func TestPrintModelResponse_PrintsLines(t *testing.T) {
	resp := makeResponse("line one\nline two")
	out := captureStdout(func() { PrintModelResponse(resp) })
	assert.Contains(t, out, "line one")
	assert.Contains(t, out, "line two")
}
