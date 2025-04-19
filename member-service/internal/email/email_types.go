package email

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Image   string `json:"image"` // Base64編碼的圖片，若無圖像可不傳遞
}
