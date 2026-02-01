package svc_ghl

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

const ghlPublicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAokvo/r9tVgcfZ5DysOSC
Frm602qYV0MaAiNnX9O8KxMbiyRKWeL9JpCpVpt4XHIcBOK4u3cLSqJGOLaPuXw6
dO0t6Q/ZVdAV5Phz+ZtzPL16iCGeK9po6D6JHBpbi989mmzMryUnQJezlYJ3DVfB
csedpinheNnyYeFXolrJvcsjDtfAeRx5ByHQmTnSdFUzuAnC9/GepgLT9SM4nCpv
uxmZMxrJt5Rw+VUaQ9B8JSvbMPpez4peKaJPZHBbU3OdeCVx5klVXXZQGNHOs8gF
3kvoV5rTnXV0IknLBXlcKKAQLZcY/Q9rG6Ifi9c+5vqlvHPCUJFT5XUGG5RKgOKU
J062fRtN+rLYZUV+BjafxQauvC8wSWeYja63VSUruvmNj8xkx2zE/Juc+yjLjTXp
IocmaiFeAO6fUtNjDeFVkhf5LNb59vECyrHD2SQIrhgXpO4Q3dVNA5rw576PwTzN
h/AMfHKIjE4xQA1SZuYJmNnmVZLIZBlQAF9Ntd03rfadZ+yDiOXCCs9FkHibELhC
HULgCsnuDJHcrGNd5/Ddm5hxGQ0ASitgHeMZ0kcIOwKDOzOU53lDza6/Y09T7sYJ
PQe7z0cvj7aE4B+Ax1ZoZGPzpJlZtGXCsu9aTEGEnKzmsFqwcSsnw3JB31IGKAyk
T1hhTiaCeIY/OwwwNUY2yvcCAwEAAQ==
-----END PUBLIC KEY-----`

var ghlPublicKey *rsa.PublicKey
var ghlPublicKeyErr error
var ghlPublicKeyOnce sync.Once

func getGHLPublicKey() (*rsa.PublicKey, error) {
	ghlPublicKeyOnce.Do(func() {
		block, _ := pem.Decode([]byte(ghlPublicKeyPEM))
		if block == nil {
			ghlPublicKeyErr = errors.New("invalid GHL public key PEM")
			return
		}

		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			ghlPublicKeyErr = err
			return
		}

		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			ghlPublicKeyErr = errors.New("GHL public key is not RSA")
			return
		}

		ghlPublicKey = rsaPub
	})
	return ghlPublicKey, ghlPublicKeyErr
}

func WebhookAuthMiddle(c *gin.Context) {
	signature := strings.TrimSpace(c.GetHeader("x-wh-signature"))
	if signature == "" {
		c.Data(lvn.Res(401, "", "missing signature"))
		c.Abort()
		return
	}
	fmt.Printf("signature: %s\n", signature)

	body, err := c.GetRawData()
	if err != nil {
		c.Data(lvn.Res(400, "", "unable to read payload"))
		c.Abort()
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	pubKey, err := getGHLPublicKey()
	if err != nil {
		c.Data(lvn.Res(500, "", "unable to load public key"))
		c.Abort()
		return
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		c.Data(lvn.Res(401, "", "invalid signature encoding"))
		c.Abort()
		return
	}
	fmt.Printf("signature: %s\n", signature)

	hash := sha256.Sum256(body)
	if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], sigBytes); err != nil {
		c.Data(lvn.Res(401, "", "invalid signature"))
		c.Abort()
		return
	}

	c.Next()
}
