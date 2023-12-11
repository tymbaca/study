package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	users = map[string]string{
		"tymbaca": "tiko77",
	}

	secret      = "secret-key"
	expTimeFlag = flag.Int64(
		"exptime",
		1000,
		"Interval in seconds after which JWT will expire",
	)
)

const (
	HS256 = "HS256"
	HS512 = "HS512"
)

// Header and Payload contains JSON strings with corresponding values
type Token struct {
	header    JWTMap
	payload   JWTMap
	signature hash.Hash
}

// alg can be HS256, HS368 or HS512
func NewToken(alg string, exp int64, name string) Token {
	h := NewHeader(HS256)
	p := NewPayload(exp, name)
	s := NewSignature(h, p)
	return Token{header: h, payload: p, signature: s}
}

func NewDefaultExp(alg, name string) Token {
	exp := time.Now().Unix() + *expTimeFlag
	t := NewToken(alg, exp, name)
	return t
}

// String returns final JWT token: Header, Payload, Signature encoded in
// base64 (base64.RawURLEncoding) separated with '.'
// e.g. "eyJhbGciOiJIUzI1NiIsCJ9.eyJzdWWF0IjoxNTE2MjM5MDIyfQ.cThIIoqwxjxSJyQQ"
func (t *Token) String() string {
	signHash := NewSignature(t.header, t.payload)
	sign := base64.RawURLEncoding.EncodeToString(signHash.Sum(nil))
	res := strings.Join(
		[]string{t.header.Base64(), t.payload.Base64(), sign},
		".",
	)
	return res
}

type JWTMap map[string]any

func (h *JWTMap) Bytes() []byte {
	data, err := json.Marshal(h)
	// It must not panic. I guess...
	if err != nil {
		panic(err)
	}
	return data
}

func (h *JWTMap) Base64() string {
	data := h.Bytes()
	return base64.RawURLEncoding.EncodeToString(data)
}

func NewHeader(alg string) JWTMap {
	m := make(JWTMap)
	m["alg"] = alg
	m["typ"] = "JWT"
	return m
}

func NewPayload(exp int64, name string) JWTMap {
	m := make(JWTMap)
	m["exp"] = exp
	m["name"] = name
	return m
}

func NewSignature(head, payl JWTMap) hash.Hash {
	joinded := head.Base64() + "." + payl.Base64()
	// "alg" key must contain string
	alg, _ := head["alg"].(string)
	var hf func() hash.Hash
	switch alg {
	case HS256:
		hf = sha256.New
	case HS512:
		hf = sha512.New
	default:
		panic("Unknown algorithm: " + alg)
	}

	hm := hmac.New(hf, []byte(secret))
	hm.Write([]byte(joinded))
	return hm
}

// Auth uses HS256 (HMAC + SHA256)
func Login(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	checkpass, ok := users[user]
	if !ok || pass != checkpass {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := NewToken(HS256, 100, user)
	w.Header().Add("Authorization", "Bearer "+token.String())
	return
}

func Auth(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")

	// Clean header
	auth, ok := strings.CutPrefix(auth, "Bearer ")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSW Token must start with 'Bearer '"))
		return
	}

	parts := strings.Split(auth, ".")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(auth))
		return
	}

	header, payload, signature := parts[0], parts[1], parts[2]

	hm := hmac.New(sha256.New, []byte("secret-key"))
	checkString := []byte(header + "." + payload)
	hm.Write(checkString)
	checkSum := hm.Sum(nil)
	check := base64.RawURLEncoding.EncodeToString(checkSum)

	message := fmt.Sprintf(
		"Your signature: '%s'. Verified signature: '%s'",
		signature,
		check,
	)
	w.Write([]byte(message))

	// ADD EXP DATE CHECK

	//=====================================================================

	// headerData, err := base64.RawURLEncoding.DecodeString(header)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// payloadData, err := base64.RawURLEncoding.DecodeString(payload)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	// w.Write(headerData)
	// w.Write(payloadData)
	// _ = signature
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", Login)
	mux.HandleFunc("/auth", Auth)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
