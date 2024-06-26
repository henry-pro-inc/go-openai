package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrStreamEventEmptyTopic    = errors.New("stream event gets an empty topic")
	ErrStreamEventCallbackPanic = errors.New("occurs a runtime panic during the stream event callback")
)

type AssistantThreadRunStreamResponse struct {
	ID             string             `json:"id"`
	Object         string             `json:"object"`
	Delta          *MessageDelta      `json:"delta,omitempty"`
	Status         string             `json:"status"`
	RequiredAction *RunRequiredAction `json:"required_action,omitempty"`
	LastError      *RunLastError      `json:"last_error,omitempty"`
	FailedAt       *int64             `json:"failed_at,omitempty"`
}

type AssistantThreadRunStream struct {
	*streamReader[AssistantThreadRunStreamResponse]
}

func (c *Client) CreateAssistantThreadRunStream(
	ctx context.Context,
	threadID string,
	request RunRequest,
) (stream *AssistantThreadRunStream, err error) {
	request.Stream = true
	urlSuffix := fmt.Sprintf("/threads/%s/runs", threadID)
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(urlSuffix),
		withBody(request),
		withBetaAssistantVersion(c.config.AssistantVersion),
	)
	if err != nil {
		return nil, err
	}

	resp, err := sendRequestStream[AssistantThreadRunStreamResponse](c, req)
	if err != nil {
		return nil, err
	}

	stream = &AssistantThreadRunStream{
		streamReader: resp,
	}
	return stream, nil
}

func (c *Client) CreateAssistantThreadRunToolStream(
	ctx context.Context,
	threadID string,
	runID string,
	request SubmitToolOutputsRequest,
) (stream *AssistantThreadRunStream, err error) {
	request.Stream = true
	urlSuffix := fmt.Sprintf("/threads/%s/runs/%s/submit_tool_outputs", threadID, runID)
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(urlSuffix),
		withBody(request),
		withBetaAssistantVersion(c.config.AssistantVersion),
	)
	if err != nil {
		return
	}

	resp, err := sendRequestStream[AssistantThreadRunStreamResponse](c, req)
	if err != nil {
		return nil, err
	}

	stream = &AssistantThreadRunStream{
		streamReader: resp,
	}
	return stream, nil
}
