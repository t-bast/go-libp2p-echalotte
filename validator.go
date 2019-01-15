package echalotte

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	pb "github.com/t-bast/go-libp2p-echalotte/pb"

	"gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	"gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
)

// Errors used by validators.
const (
	ErrInvalidKeyFormat       = "invalid DHT key format"
	ErrInvalidNamespace       = "invalid DHT key namespace"
	ErrInvalidSenderSignature = "invalid sender signature"
)

const (
	// EncryptionNamespace is the namespace used for storing encryption public
	// keys on a DHT for node-to-node encryption.
	EncryptionNamespace = "enc"
)

// PublicKeyValidator validates public keys used for node-to-node encryption
// before storing them in the DHT.
type PublicKeyValidator struct{}

// CreateKey returns a namespaced DHT key for the given peer's encryption key.
func (pkv PublicKeyValidator) CreateKey(peerID peer.ID) string {
	return fmt.Sprintf("/%s/%s", EncryptionNamespace, peerID.Pretty())
}

// CreateRecord creates a record for a Curve25519 encryption key.
// This record is suitable for storage on a DHT.
func (pkv PublicKeyValidator) CreateRecord(signingKey crypto.PrivKey, encryptionKey *[32]byte) ([]byte, error) {
	publicKey := &pb.PublicKey{
		Type:      pb.KeyType_Curve25519,
		CreatedAt: ptypes.TimestampNow(),
		Data:      encryptionKey[:],
	}

	toSign, err := proto.Marshal(publicKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	publicKey.Signature, err = signingKey.Sign(toSign)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	publicKey.SignatureKey, err = signingKey.GetPublic().Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	serialized, err := proto.Marshal(publicKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return serialized, nil
}

// Validate the node-to-node encryption record.
func (pkv PublicKeyValidator) Validate(key string, value []byte) error {
	peerID, err := pkv.getPeerID(key)
	if err != nil {
		return err
	}

	var publicKey pb.PublicKey
	err = proto.Unmarshal(value, &publicKey)
	if err != nil {
		return errors.WithStack(err)
	}

	signatureKey, err := crypto.UnmarshalPublicKey(publicKey.SignatureKey)
	if err != nil {
		return errors.WithStack(err)
	}

	if !peerID.MatchesPublicKey(signatureKey) {
		return errors.New(ErrInvalidSenderSignature)
	}

	signature := publicKey.Signature

	publicKey.SignatureKey = nil
	publicKey.Signature = nil

	signedBytes, err := proto.Marshal(&publicKey)
	if err != nil {
		return errors.WithStack(err)
	}

	ok, err := signatureKey.Verify(signedBytes, signature)
	if err != nil {
		return errors.Wrap(err, ErrInvalidSenderSignature)
	}
	if !ok {
		return errors.New(ErrInvalidSenderSignature)
	}

	// No need to validate that the point is on the curve because we only use
	// curve25519 for now which has twist security.
	// If we support more elliptic curves, we might need to check here that the
	// public key received is a valid curve point.

	return nil
}

// Select the most recently published encryption key.
// We likely need to work on a more robust story around key revokation.
func (pkv PublicKeyValidator) Select(_ string, values [][]byte) (int, error) {
	i := 0
	createdAt := int64(0)

	for index, value := range values {
		var publicKey pb.PublicKey
		err := proto.Unmarshal(value, &publicKey)
		if err != nil {
			continue
		}

		if publicKey.CreatedAt.GetSeconds() > createdAt {
			i = index
			createdAt = publicKey.CreatedAt.GetSeconds()
		}
	}

	return i, nil
}

// getPeerID takes a key in the form `/enc/$peerID` and extracts the peer ID.
func (pkv PublicKeyValidator) getPeerID(key string) (peer.ID, error) {
	if len(key) == 0 || key[0] != '/' {
		return "", errors.New(ErrInvalidKeyFormat)
	}

	key = key[1:]

	i := strings.IndexByte(key, '/')
	if i <= 0 {
		return "", errors.New(ErrInvalidKeyFormat)
	}

	ns := key[:i]
	if ns != EncryptionNamespace {
		return "", errors.New(ErrInvalidNamespace)
	}

	peerID, err := peer.IDB58Decode(key[i+1:])
	if err != nil {
		return "", errors.WithStack(err)
	}

	return peerID, nil
}
