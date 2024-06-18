package model

type ConnectionMessage struct{
	ViewerUserId string `json:"viewerid" binding:"required"`
	SdpConnectionString string `json:"sdp" binding:"required"`
}