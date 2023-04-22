package dto

type TotalDocument struct {
	TotalDocument int `bson:"totalDocument"`
}

type ImageAndNote struct {
	ImageURL string `bson:"imageURL" json:"imageURL"`
	Note     string `bson:"note" json:"note"`
}
