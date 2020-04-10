package controller

// Controller is 外部出力装置へのinterface
type Controller interface {
	// Reply replies a message
	Reply() error
}
