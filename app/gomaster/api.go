package main

type StartGameRequest struct {
	NumberOfColours   int `json:"numberOfColours" binding:"required"`
	NumberOfPositions int `json:"numberOfPositions" binding:"required"`
}

type StartGameResponse struct {
	GameId string `json:"id" binding:"required"`
}

type QueryGameResponse struct {
	NumberOfColours   int                        `json:"numberOfColours" binding:"required"`
	NumberOfPositions int                        `json:"numberOfPositions" binding:"required"`
	Colours           []Colour                   `json:"colours" binding:"required"`
	Complete          bool                       `json:"complete" binding:"required"`
	Attempts          []QueryGameAttemptResponse `json:"guesses" binding:"required"`
}

type QueryGameAttemptResponse struct {
	Positions             []Colour `json:"positions" binding:"required"`
	Complete              bool     `json:"complete" binding:"required"`
	RightColourRightPlace int      `json:"rightColourRightPlace" binding:"required"`
	RightColourWrongPlace int      `json:"rightColourWrongPlace" binding:"required"`
	WrongColour           int      `json:"wrongColour" binding:"required"`
}

type SubmitGuessRequest struct {
	Positions []Colour `json:"positions" binding:"required"`
}

type SubmitGuessResponse struct {
}
