package dto

type TotalDocument struct {
	Total int `bson:"total"`
}

type ImageAndNote struct {
	ImageURL string `bson:"imageURL" json:"imageURL"`
	Note     string `bson:"note" json:"note"`
}