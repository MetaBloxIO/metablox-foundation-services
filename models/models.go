package models

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
)

type DIDDocument struct {
	Context            []string             `json:"@context"`
	ID                 string               `json:"id"`
	Created            string               `json:"created"`
	Updated            string               `json:"updated"`
	Version            int                  `json:"version"`
	VerificationMethod []VerificationMethod `json:"verificationMethod"`
	Authentication     string               `json:"authentication"`
	Service            []Service            `json:"service"`
}

type VerificationMethod struct {
	ID           string `json:"id"`
	MethodType   string `json:"type" mapstructure:"type"`
	Controller   string `json:"controller"`
	MultibaseKey string `json:"publicKeyMultibase" mapstructure:"publicKeyMultibase"`
}

type ResolutionOptions struct {
	Accept string `json:"accept"`
}

type RepresentationResolutionOptions struct {
	Accept string `json:"accept"`
}

type ResolutionMetadata struct {
	Error string `json:"error"`
}

type RepresentationResolutionMetadata struct {
	ContentType string `json:"contentType"`
	Error       string `json:"error"`
}

type DocumentMetadata struct {
	Created       string   `json:"created"`
	Updated       string   `json:"updated"`
	Deactivated   string   `json:"deactivated"`
	NextUpdate    string   `json:"nextUpdate"`
	VersionID     string   `json:"versionId"`
	NextVersionID string   `json:"nextVersionId"`
	EquivalentID  []string `json:"equivalentId"`
	CanonicalID   string   `json:"canonicalId"`
}

type VerifiableCredential struct {
	Context           []string    `json:"@context" mapstructure:"@context"`
	ID                string      `json:"id" db:"ID"`
	Type              []string    `json:"type"`
	SubType           string      `json:"-" db:"Type"` //used in place of Type for database operations, as using an array causes issues
	Issuer            string      `json:"issuer" db:"Issuer"`
	IssuanceDate      string      `json:"issuanceDate" db:"IssuanceDate"`
	ExpirationDate    string      `json:"expirationDate" db:"ExpirationDate"`
	Description       string      `json:"description" db:"Description"`
	CredentialSubject interface{} `json:"credentialSubject"`
	Proof             VCProof     `json:"proof"`
	Revoked           bool        `json:"revoked" db:"Revoked"`
}

//This can be a type of input form to set up the VC.
//Temp fields here currently, will be changed in the future
type SubjectInfo struct {
	ID           string `json:"id"`
	GivenName    string `json:"givenName"`
	FamilyName   string `json:"familyName"`
	Gender       string `json:"gender"`
	BirthCountry string `json:"birthCountry"`
	BirthDate    string `json:"birthName"`
}

type VCProof struct {
	Type               string `json:"type"`
	Created            string `json:"created"`
	VerificationMethod string `json:"verificationMethod"`
	ProofPurpose       string `json:"proofPurpose"`
	JWSSignature       string `json:"jws"` //signature is created from a hash of the VC
}

type VPProof struct {
	Type               string `json:"type"`
	Created            string `json:"created"`
	VerificationMethod string `json:"verificationMethod"`
	ProofPurpose       string `json:"proofPurpose"`
	JWSSignature       string `json:"jws" mapstructure:"jws"` //signature is created from a hash of the VP
	Nonce              string `json:"nonce"`                  //random value generated by verifier that must be included in proof
}

type VerifiablePresentation struct {
	Context              []string               `json:"@context" mapstructure:"@context"`
	Type                 []string               `json:"type"`
	VerifiableCredential []VerifiableCredential `json:"verifiableCredential"`
	Holder               string                 `json:"holder"`
	Proof                VPProof                `json:"proof"`
}

type Service struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

type WifiAccessInfo struct {
	ID                   string `json:"id"`
	PlaceholderParameter string `json:"placeholderParameter"`
}

type MiningLicenseInfo struct {
	ID                    string `json:"id"`
	PlaceholderParameter2 string `json:"placeholderParameter2"`
}

func CreateDIDDocument() *DIDDocument {
	return &DIDDocument{}
}

func (doc DIDDocument) RetrieveVerificationMethod(vmID string) (VerificationMethod, error) {
	for _, vm := range doc.VerificationMethod {
		if vm.ID == vmID {
			return vm, nil
		}
	}
	return VerificationMethod{}, errors.New("failed to find verification method with ID " + vmID)
}

func (doc *DIDDocument) AddService(service Service) {
	doc.Service = append(doc.Service, service)
}

func CreateVerifiableCredential() *VerifiableCredential {
	return &VerifiableCredential{}
}

func NewVerifiableCredential(context []string, id string, vctype []string, subtype, issuer, issuanceDate, expirationDate, description string, subject interface{}, proof VCProof, revoked bool) *VerifiableCredential {
	return &VerifiableCredential{
		Context:           context,
		ID:                id,
		Type:              vctype,
		SubType:           subtype,
		Issuer:            issuer,
		IssuanceDate:      issuanceDate,
		ExpirationDate:    expirationDate,
		Description:       description,
		CredentialSubject: subject,
		Proof:             proof,
		Revoked:           revoked,
	}
}

func CreateSubjectInfo() *SubjectInfo {
	return &SubjectInfo{}
}

func CreateWifiAccessInfo() *WifiAccessInfo {
	return &WifiAccessInfo{}
}

func CreateMiningLicenseInfo() *MiningLicenseInfo {
	return &MiningLicenseInfo{}
}

func CreateVCProof() *VCProof {
	return &VCProof{}
}

func CreateVPProof() *VPProof {
	return &VPProof{}
}

func CreateResolutionOptions() *ResolutionOptions {
	return &ResolutionOptions{}
}

func CreatePresentation() *VerifiablePresentation {
	return &VerifiablePresentation{}
}

func NewPresentation(context, presentationType []string, credentials []VerifiableCredential, holder string, proof VPProof) *VerifiablePresentation {
	return &VerifiablePresentation{
		Context:              context,
		Type:                 presentationType,
		VerifiableCredential: credentials,
		Holder:               holder,
		Proof:                proof,
	}
}

func GenerateTestPrivKey() *ecdsa.PrivateKey {
	privKey, _ := crypto.ToECDSA([]byte{165, 190, 153, 12, 246, 178, 211, 170, 147, 144, 51, 73, 48, 27, 20, 79, 61, 110, 201, 118, 99, 219, 50, 252, 135, 12, 107, 237, 245, 95, 170, 17})
	return privKey
}

func GenerateTestDIDDocument() *DIDDocument {
	document := CreateDIDDocument()
	document.Context = append(document.Context, "https://w3id.org/did/v1")
	document.Context = append(document.Context, "https://ns.did.ai/suites/secp256k1-2019/v1/")
	document.ID = "did:metablox:HFXPiudexfvsJBqABNmBp785YwaKGjo95kmDpBxhMMYo"
	document.Created = "2022-03-31T12:53:19-07:00"
	document.Updated = "2022-03-31T12:53:19-07:00"
	document.Version = 1
	document.VerificationMethod = append(document.VerificationMethod, VerificationMethod{ID: "did:metablox:HFXPiudexfvsJBqABNmBp785YwaKGjo95kmDpBxhMMYo#verification", MethodType: "EcdsaSecp256k1VerificationKey2019", Controller: "did:metablox:HFXPiudexfvsJBqABNmBp785YwaKGjo95kmDpBxhMMYo", MultibaseKey: "zR4TQJaWaLA3vvYukULRJoxTsRmqCMsWuEJdDE8CJwRFCUijDGwCBP89xVcWdLRQaEM6b7wD294xCs8byy3CdDMYa"})
	document.Authentication = "did:metablox:HFXPiudexfvsJBqABNmBp785YwaKGjo95kmDpBxhMMYo#verification"
	return document
}

func NewSubjectInfo(id string, givenName, familyName, gender, birthCountry, birthDate string) *SubjectInfo {
	return &SubjectInfo{
		ID:           id,
		GivenName:    givenName,
		FamilyName:   familyName,
		Gender:       gender,
		BirthCountry: birthCountry,
		BirthDate:    birthDate,
	}
}

func GenerateTestSubjectInfo() *SubjectInfo {
	return NewSubjectInfo(
		"did:metablox:HFXPiudexfvsJBqABNmBp785YwaKGjo95kmDpBxhMMYo",
		"John",
		"Jacobs",
		"Male",
		"Canada",
		"2022-03-22",
	)
}

func NewWifiAccessInfo(id, placeholder string) *WifiAccessInfo {
	return &WifiAccessInfo{
		ID:                   id,
		PlaceholderParameter: placeholder,
	}
}

func GenerateTestWifiAccessInfo() *WifiAccessInfo {
	return NewWifiAccessInfo(
		"sampleID",
		"ThisIsAPlaceholder",
	)
}

func NewMiningLicenseInfo(id, placeholder string) *MiningLicenseInfo {
	return &MiningLicenseInfo{
		ID:                    id,
		PlaceholderParameter2: placeholder,
	}
}

func GenerateTestMiningLicenseInfo() *MiningLicenseInfo {
	return NewMiningLicenseInfo(
		"sampleID2",
		"ThisIsAPlaceholder2",
	)
}

func NewVCProof(proofType, created, vm, purpose, sig string) *VCProof {
	return &VCProof{
		Type:               proofType,
		Created:            created,
		VerificationMethod: vm,
		ProofPurpose:       purpose,
		JWSSignature:       sig,
	}
}

func NewVPProof(proofType, created, vm, purpose, sig, nonce string) *VPProof {
	return &VPProof{
		Type:               proofType,
		Created:            created,
		VerificationMethod: vm,
		ProofPurpose:       purpose,
		JWSSignature:       sig,
		Nonce:              nonce,
	}
}

func GenerateTestVC() *VerifiableCredential {
	vcProof := NewVCProof(
		"EcdsaSecp256k1Signature2019",
		"2022-03-31T12:53:19-07:00",
		"did:metablox:HFXPiudexfvsJBqABNmBp785YwaKGjo95kmDpBxhMMYo#verification",
		"Authentication",
		"eyJhbGciOiJFUzI1NiJ9..lvLWxsW_5GIZGCztNs_ioBHHyC4PZ1JP9CQL0NgdTwjf7EHMDgCViLzLwv_FFJtYSEUh7Y67VbIFhM50B5cnxg",
	)

	subjectInfo := GenerateTestSubjectInfo()

	return NewVerifiableCredential(
		[]string{"https://www.w3.org/2018/credentials/v1", "https://ns.did.ai/suites/secp256k1-2019/v1/"},
		"http://metablox.com/credentials/1",
		[]string{"VerifiableCredential", "PermanentResidentCard"},
		"PermanentResidentCard",
		"did:metablox:sampleIssuer",
		"2022-03-31T12:53:19-07:00",
		"2032-03-31T12:53:19-07:00",
		"Government of Example Permanent Resident Card",
		*subjectInfo,
		*vcProof,
		false,
	)
}

func GenerateTestWifiAccessVC() *VerifiableCredential {
	vcProof := NewVCProof(
		"EcdsaSecp256k1Signature2019",
		"2022-03-31T12:53:19-07:00",
		"did:metablox:HFXPiudexfvsJBqABNmBp785YwaKGjo95kmDpBxhMMYo#verification",
		"Authentication",
		"eyJhbGciOiJFUzI1NiJ9..zYvdJMDdwS8IBuXMYCzLSdU_VBn5iG6bYzSIKz366O_KkP0bJ2fV3sUmzzQM7CBBuRSOPH08CAeFzoNXIl0LdA",
	)

	wifiAccessInfo := GenerateTestWifiAccessInfo()

	return NewVerifiableCredential(
		[]string{"https://www.w3.org/2018/credentials/v1", "https://ns.did.ai/suites/secp256k1-2019/v1/"},
		"http://metablox.com/credentials/1",
		[]string{"VerifiableCredential", "WifiAccess"},
		"WifiAccess",
		"did:metablox:sampleIssuer",
		"2022-03-31T12:53:19-07:00",
		"2032-03-31T12:53:19-07:00",
		"Example Wifi Access Credential",
		*wifiAccessInfo,
		*vcProof,
		false,
	)
}

func GenerateTestMiningLicenseVC() *VerifiableCredential {
	vcProof := NewVCProof(
		"EcdsaSecp256k1Signature2019",
		"2022-03-31T12:53:19-07:00",
		"did:metablox:HFXPiudexfvsJBqABNmBp785YwaKGjo95kmDpBxhMMYo#verification",
		"Authentication",
		"eyJhbGciOiJFUzI1NiJ9..zYvdJMDdwS8IBuXMYCzLSdU_VBn5iG6bYzSIKz366O_KkP0bJ2fV3sUmzzQM7CBBuRSOPH08CAeFzoNXIl0LdA",
	)

	miningLicenseInfo := GenerateTestMiningLicenseInfo()

	return NewVerifiableCredential(
		[]string{"https://www.w3.org/2018/credentials/v1", "https://ns.did.ai/suites/secp256k1-2019/v1/"},
		"http://metablox.com/credentials/1",
		[]string{"VerifiableCredential", "MiningLicense"},
		"MiningLicense",
		"did:metablox:sampleIssuer",
		"2022-03-31T12:53:19-07:00",
		"2032-03-31T12:53:19-07:00",
		"Example Mining License Credential",
		*miningLicenseInfo,
		*vcProof,
		false,
	)
}

func CreateService() *Service {
	return &Service{}
}

func GenerateTestPresentation() *VerifiablePresentation {
	vpProof := NewVPProof(
		"EcdsaSecp256k1Signature2019",
		"2022-03-31T12:53:19-07:00",
		"did:metablox:HFXPiudexfvsJBqABNmBp785YwaKGjo95kmDpBxhMMYo#verification",
		"Authentication",
		"eyJhbGciOiJFUzI1NiJ9..bmj6KhHcBkLOHgAZrLqgweE-StyBXvvj6bmZqC6TqiYVtC_tXf076xDAAXzmx160dAqivTzgX-943ZU-VWXDqw",
		"sampleNonce",
	)

	return NewPresentation(
		[]string{"https://www.w3.org/2018/credentials/v1", "https://ns.did.ai/suites/secp256k1-2019/v1/"},
		[]string{"VerifiablePresentation"},
		[]VerifiableCredential{*GenerateTestVC()},
		"did:metablox:HFXPiudexfvsJBqABNmBp785YwaKGjo95kmDpBxhMMYo",
		*vpProof,
	)
}

func GenerateTestWifiAccessPresentation() *VerifiablePresentation {
	vpProof := NewVPProof(
		"EcdsaSecp256k1Signature2019",
		"2022-03-31T12:53:19-07:00",
		"did:metablox:HFXPiudexfvsJBqABNmBp785YwaKGjo95kmDpBxhMMYo#verification",
		"Authentication",
		"eyJhbGciOiJFUzI1NiJ9..bmj6KhHcBkLOHgAZrLqgweE-StyBXvvj6bmZqC6TqiYVtC_tXf076xDAAXzmx160dAqivTzgX-943ZU-VWXDqw",
		"sampleNonce",
	)

	return NewPresentation(
		[]string{"https://www.w3.org/2018/credentials/v1", "https://ns.did.ai/suites/secp256k1-2019/v1/"},
		[]string{"VerifiablePresentation"},
		[]VerifiableCredential{*GenerateTestWifiAccessVC()},
		"did:metablox:HFXPiudexfvsJBqABNmBp785YwaKGjo95kmDpBxhMMYo",
		*vpProof,
	)
}
