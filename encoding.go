package cqrskit

import "encoding/json"

// JSONEncoder implements the cqrskit.Encoder to encode EventCommits types.
type JSONEncoder struct{}

// Encode attempts to encode a EventCommit into a json byte slice.
// It returns an error if it failed.
func (JSONEncoder) Encode(commit EventCommit) ([]byte, error) {
	return json.Marshal(commit)
}

// JSONEncoder implements the cqrskit.Decoder to decode byte slices to EventCommits types.
type JSONDecoder struct{}

// Decode attempts to decode byte slice of json into a EventCommit.
// It returns an error if it failed.
func (JSONDecoder) Decode(data []byte) (EventCommit, error) {
	var commit EventCommit
	err := json.Unmarshal(data, &commit)
	return commit, err
}
