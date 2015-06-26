package keys

import (
	"crypto/elliptic"
	"encoding/json"
	"math/big"
)

// TLSCurveID is the type of a TLS identifier for an elliptic curve. See
// http://www.iana.org/assignments/tls-parameters/tls-parameters.xml#tls-parameters-8
type TLSCurveID uint16

// ECDHParams stores elliptic-curve Diffie-Hellman paramters.At any point in
// time, it is unlikely that both ServerPrivate and ClientPrivate will be non-nil.
type ECDHParams struct {
	TLSCurveID   TLSCurveID     `json:"curve_id,omitempty"`
	Curve        elliptic.Curve `json:"-"`
	ServerPublic *ECPoint       `json:"server_public,omitempty"`
}

// ECPoint represents an elliptic curve point and serializes nicely to JSON
type ECPoint struct {
	X *big.Int
	Y *big.Int
}

// MarshalJSON implements the json.Marshler interface
func (p *ECPoint) MarshalJSON() ([]byte, error) {
	aux := struct {
		X *cryptoParameter `json:"x"`
		Y *cryptoParameter `json:"y"`
	}{
		X: &cryptoParameter{Int: p.X},
		Y: &cryptoParameter{Int: p.Y},
	}
	return json.Marshal(&aux)
}

// UnmarshalJSON implements the json.Unmarshler interface
func (p *ECPoint) UnmarshalJSON(b []byte) error {
	aux := struct {
		X *cryptoParameter `json:"x"`
		Y *cryptoParameter `json:"y"`
	}{}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	p.X = aux.X.Int
	p.Y = aux.Y.Int
	return nil
}

// Description returns the description field for the given ID. See
// http://www.iana.org/assignments/tls-parameters/tls-parameters.xml#tls-parameters-8
func (c *TLSCurveID) Description() string {
	if desc, ok := ecIDToName[*c]; ok {
		return desc
	}
	return "unknown"
}

// MarshalJSON implements the json.Marshaler interface
func (c *TLSCurveID) MarshalJSON() ([]byte, error) {
	aux := struct {
		Name string `json:"name"`
		ID   uint16 `json:"id"`
	}{
		Name: c.Description(),
		ID:   uint16(*c),
	}
	return json.Marshal(&aux)
}

//UnmarshalJSON implements the json.Unmarshaler interface
func (c *TLSCurveID) UnmarshalJSON(b []byte) error {
	aux := struct {
		ID uint16 `json:"id"`
	}{}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	*c = TLSCurveID(aux.ID)
	return nil
}