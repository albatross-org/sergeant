package server

import "github.com/albatross-org/sergeant"

// CardJSON is the response returned when a client asks for a card.
type CardJSON struct {
	Path        string `json:"path"`
	QuestionImg string `json:"questionImg"`
	AnswerImg   string `json:"answerImg"`
	ID          string `json:"id"`
}

// cardToJSON converts a *sergeant.Card into the JSON format ready to be accepted by the client.
// If an error is returned, it's due to an issue with converting the card's contents to a data URI.
func cardToJSON(card *sergeant.Card) (CardJSON, error) {
	questionImg, err := card.QuestionImage()
	if err != nil {
		return CardJSON{}, err
	}

	answerImg, err := card.AnswerImage()
	if err != nil {
		return CardJSON{}, err
	}

	return CardJSON{
		Path:        card.PathParent(),
		QuestionImg: questionImg,
		AnswerImg:   answerImg,
		ID:          card.ID,
	}, nil
}

// CardUpdateJSON is what is sent to the server when a client wants to update a card.
type CardUpdateJSON struct {
	ID       string `json:"id"`
	Answer   string `json:"answer"`
	Duration int    `json:"duration"`
}
