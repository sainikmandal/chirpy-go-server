package main

import db "github.com/sainikmandal/chirpy-go-server/internal/database"

type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func dbUserToUserResponse(user db.User) UserResponse {
	return UserResponse{
		ID:          user.ID.String(),
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
}
