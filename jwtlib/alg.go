package jwtlib

import "github.com/dgrijalva/jwt-go"

type SignALG jwt.SigningMethod

var HS256 SignALG = jwt.SigningMethodHS256
