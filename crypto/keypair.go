package crypto

import (
	"blockx/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

// PrivateKey 私钥
type PrivateKey struct {
	key *ecdsa.PrivateKey
}

// Sign 用私钥签名
func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

// GeneratePrivateKey 生成私钥
func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return PrivateKey{key: key}
}

// PublicKey 从私钥获取公钥
func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{Key: &k.key.PublicKey}
}

// PublicKey 公钥
type PublicKey struct {
	Key *ecdsa.PublicKey
}

// Bytes 公钥转字节数组
func (k PublicKey) Bytes() []byte {
	return elliptic.Marshal(k.Key, k.Key.X, k.Key.Y)
}

// Address 由公钥生成地址
func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.Bytes())
	return types.AddressFromBytes(h[len(h)-20:])
}

// Signature 签名
type Signature struct {
	R, S *big.Int
}

// Verify 验证签名
func (sig Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.Key, data, sig.R, sig.S)
}
