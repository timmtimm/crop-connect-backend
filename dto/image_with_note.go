package dto

type ImageAndNote struct {
	ImageURL string `bson:"imageURL" json:"imageURL"`
	Note     string `bson:"note" json:"note"`
}
