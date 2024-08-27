package domain

type SignupRequest struct {
	FirstName string `json:"first_name" bson:"first_name" binding:"required,min=3,max=30"`
	LastName  string `json:"last_name" bson:"last_name" binding:"required,min=3,max=30"`
	Email     string `json:"email" bson:"email" binding:"required,email"`
	Password  string `json:"password" bson:"password" binding:"required,min=4,max=30"`
}
type SignupResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
