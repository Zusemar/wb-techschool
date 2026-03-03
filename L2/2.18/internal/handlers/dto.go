package handlers

type createEventRequest struct {
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Text   string `json:"text"`
}

type updateEventRequest struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Text   string `json:"text"`
}

type deleteEventRequest struct {
	ID int `json:"id"`
}

type listEventsRequest struct {
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
}

type eventResponse struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Text   string `json:"text"`
}

type idResponse struct {
	ID int `json:"id"`
}

type errorResponse struct {
	Error string `json:"error"`
}
