/*
 * SPDX-License-Identifier: Apache-2.0
 */

package util

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/IBM-Blockchain/fablet/internal/pkg/identity"
	"github.com/golang/protobuf/proto"
)

// MarshalOrPanic marshals the specified Protocol Buffer message into a byte array, and panics on failure.
func MarshalOrPanic(pb proto.Message) []byte {
	res, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	return res
}

const config = `NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/ca.pem
    OrganizationalUnitIdentifier: client
  AdminOUIdentifier:
    Certificate: cacerts/ca.pem
    OrganizationalUnitIdentifier: admin
  PeerOUIdentifier:
    Certificate: cacerts/ca.pem
    OrganizationalUnitIdentifier: peer
  OrdererOUIdentifier:
    Certificate: cacerts/ca.pem
    OrganizationalUnitIdentifier: orderer
`

// CreateMSPDirectory creates an MSP directory on disk suitable for the peer or orderer to use.
func CreateMSPDirectory(directory string, identity *identity.Identity) error {
	directories := []string{
		directory,
		path.Join(directory, "admincerts"),
		path.Join(directory, "cacerts"),
		path.Join(directory, "keystore"),
		path.Join(directory, "signcerts"),
	}
	for _, directory := range directories {
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			return err
		}
	}
	err := ioutil.WriteFile(path.Join(directory, "config.yaml"), []byte(config), 0644)
	if err != nil {
		return err
	}
	privateKey := identity.PrivateKey().Bytes()
	err = ioutil.WriteFile(path.Join(directory, "keystore", "key.pem"), privateKey, 0644)
	if err != nil {
		return err
	}
	certificate := identity.Certificate().Bytes()
	err = ioutil.WriteFile(path.Join(directory, "signcerts", "cert.pem"), certificate, 0644)
	if err != nil {
		return err
	}
	if hasCA := identity.CA() != nil; hasCA {
		ca := identity.CA().Bytes()
		err = ioutil.WriteFile(path.Join(directory, "cacerts", "ca.pem"), ca, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
