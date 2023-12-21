package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmtericKey []byte
}

func NewPasetoMaker(symmtericKey string) (Maker, error) {

	if len(symmtericKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("Invalid Key size")
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmtericKey: []byte(symmtericKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmtericKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {

	payload:=&Payload{}

	err:=maker.paseto.Decrypt(token,maker.symmtericKey,payload,nil)

	if err!=nil{
		return nil,errors.New("Invalid Token")
	}

	err=payload.Valid()

	if err!=nil{
		return nil,err
	}

	return payload ,nil
}
