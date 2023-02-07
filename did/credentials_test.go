package did

import (
	"os"
	"testing"

	"github.com/MetaBloxIO/metablox-foundation-services/dao"
	"github.com/MetaBloxIO/metablox-foundation-services/errval"
	"github.com/MetaBloxIO/metablox-foundation-services/models"
	"github.com/MetaBloxIO/metablox-foundation-services/settings"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestCreateMiningLicenseVC(t *testing.T) {
	err := settings.Init()
	assert.Nil(t, err)
	err = dao.TestDBInit()
	assert.Nil(t, err)
	privData, err := os.ReadFile("../wallet/issuer")
	assert.Nil(t, err)
	IssuerPrivateKey, err = crypto.ToECDSA(privData)
	assert.Nil(t, err)
	t.Cleanup(dao.Close)
	t.Cleanup(dao.DeleteTestCredentialsTable)
	t.Cleanup(dao.DeleteTestMiningLicenseTable)

	err = dao.CreateTestCredentialsTable()
	assert.Nil(t, err)
	err = dao.CreateTestMiningLicenseTable()
	assert.Nil(t, err)

	issuerDocument := models.GenerateTestDIDDocument()
	issuerDocument.ID = "did:metablox:sampleIssuer"
	miningLicenseInfo := models.GenerateTestMiningLicenseInfo()
	issuerPrivKey := models.GenerateTestPrivKey()
	vc, err := CreateMiningLicenseVC(issuerDocument, miningLicenseInfo, issuerPrivKey)
	assert.Nil(t, err)
	sampleVC := models.GenerateTestMiningLicenseVC()
	assert.Equal(t, sampleVC.Context, vc.Context)
	assert.Equal(t, sampleVC.Type, vc.Type)
	assert.Equal(t, sampleVC.Issuer, vc.Issuer)
	assert.Equal(t, sampleVC.Description, vc.Description)
	assert.Equal(t, sampleVC.Proof.Type, vc.Proof.Type)
	assert.Equal(t, sampleVC.Proof.ProofPurpose, vc.Proof.ProofPurpose)
	assert.Equal(t, sampleVC.Proof.VerificationMethod, vc.Proof.VerificationMethod)
	createdSubjectInfo := vc.CredentialSubject.(models.MiningLicenseInfo)
	assert.Equal(t, miningLicenseInfo.ID, createdSubjectInfo.ID)
	assert.Equal(t, miningLicenseInfo.CredentialID, createdSubjectInfo.CredentialID)
	assert.Equal(t, miningLicenseInfo.Name, createdSubjectInfo.Name)
	assert.Equal(t, miningLicenseInfo.Model, createdSubjectInfo.Model)
	assert.Equal(t, miningLicenseInfo.Serial, createdSubjectInfo.Serial)
	vc, err = CreateMiningLicenseVC(issuerDocument, miningLicenseInfo, issuerPrivKey)
	assert.Nil(t, err)
	assert.Equal(t, sampleVC.Context, vc.Context)
	assert.Equal(t, sampleVC.Type, vc.Type)
	assert.Equal(t, sampleVC.Issuer, vc.Issuer)
	assert.Equal(t, sampleVC.Description, vc.Description)
	assert.Equal(t, sampleVC.Proof.Type, vc.Proof.Type)
	assert.Equal(t, sampleVC.Proof.ProofPurpose, vc.Proof.ProofPurpose)
	assert.Equal(t, sampleVC.Proof.VerificationMethod, vc.Proof.VerificationMethod)
	createdSubjectInfo = vc.CredentialSubject.(models.MiningLicenseInfo)
	assert.Equal(t, miningLicenseInfo.ID, createdSubjectInfo.ID)
	assert.Equal(t, miningLicenseInfo.CredentialID, createdSubjectInfo.CredentialID)
	assert.Equal(t, miningLicenseInfo.Name, createdSubjectInfo.Name)
	assert.Equal(t, miningLicenseInfo.Model, createdSubjectInfo.Model)
	assert.Equal(t, miningLicenseInfo.Serial, createdSubjectInfo.Serial)
}

func TestCreateWifiAccessVC(t *testing.T) {
	err := settings.Init()
	assert.Nil(t, err)
	err = dao.TestDBInit()
	assert.Nil(t, err)
	t.Cleanup(dao.Close)
	t.Cleanup(dao.DeleteTestCredentialsTable)
	t.Cleanup(dao.DeleteTestWifiAccessTable)

	err = dao.CreateTestCredentialsTable()
	assert.Nil(t, err)
	err = dao.CreateTestWifiAccessTable()
	assert.Nil(t, err)

	issuerDocument := models.GenerateTestDIDDocument()
	issuerDocument.ID = "did:metablox:sampleIssuer"
	wifiAccessInfo := models.GenerateTestWifiAccessInfo()
	issuerPrivKey := models.GenerateTestPrivKey()
	vc, err := CreateWifiAccessVC(issuerDocument, wifiAccessInfo, issuerPrivKey)
	assert.Nil(t, err)
	sampleVC := models.GenerateTestWifiAccessVC()
	assert.Equal(t, sampleVC.Context, vc.Context)
	assert.Equal(t, sampleVC.Type, vc.Type)
	assert.Equal(t, sampleVC.Issuer, vc.Issuer)
	assert.Equal(t, sampleVC.Description, vc.Description)
	assert.Equal(t, sampleVC.Proof.Type, vc.Proof.Type)
	assert.Equal(t, sampleVC.Proof.ProofPurpose, vc.Proof.ProofPurpose)
	assert.Equal(t, sampleVC.Proof.VerificationMethod, vc.Proof.VerificationMethod)
	createdSubjectInfo := vc.CredentialSubject.(models.WifiAccessInfo)
	assert.Equal(t, wifiAccessInfo.ID, createdSubjectInfo.ID)
	assert.Equal(t, wifiAccessInfo.CredentialID, createdSubjectInfo.CredentialID)
	assert.Equal(t, wifiAccessInfo.Type, createdSubjectInfo.Type)
}

func TestRenewVC(t *testing.T) {
	err := settings.Init()
	assert.Nil(t, err)
	err = dao.TestDBInit()
	assert.Nil(t, err)
	t.Cleanup(dao.Close)
	t.Cleanup(dao.DeleteTestCredentialsTable)

	err = dao.CreateTestCredentialsTable()
	assert.Nil(t, err)

	vc := models.GenerateTestVC()
	vc.ExpirationDate = "2022-03-31T12:53:19-07:00"
	err = dao.InsertSampleIntoCredentials(vc, "PermanentResidentCard")
	assert.Nil(t, err)
	issuerPrivKey := models.GenerateTestPrivKey()

	err = RenewVC(vc, issuerPrivKey)
	assert.Nil(t, err)
	assert.Equal(t, "2023-03-31T12:53:19-07:00", vc.ExpirationDate)

	pubKey, err := crypto.UnmarshalPubkey(vc.Proof.PublicKeyString)
	assert.Nil(t, err)
	success, err := VerifyVCSecp256k1(vc, pubKey)
	assert.Nil(t, err)
	assert.True(t, success)

	dbVC, err := dao.RetrieveSampleFromCredentials("1")
	assert.Nil(t, err)
	assert.Equal(t, "2023-03-31 12:53:19", dbVC.ExpirationDate)
}

func TestRevokeVC(t *testing.T) {
	err := settings.Init()
	assert.Nil(t, err)
	err = dao.TestDBInit()
	assert.Nil(t, err)
	t.Cleanup(dao.Close)
	t.Cleanup(dao.DeleteTestCredentialsTable)

	err = dao.CreateTestCredentialsTable()
	assert.Nil(t, err)

	vc := models.GenerateTestWifiAccessVC()
	err = dao.InsertSampleIntoCredentials(vc, "PermanentResidentCard")
	assert.Nil(t, err)

	err = RevokeVC(vc)
	assert.Nil(t, err)
	assert.True(t, vc.Revoked)

	dbVC, err := dao.RetrieveSampleFromCredentials("1")
	assert.Nil(t, err)
	assert.True(t, dbVC.Revoked)
}

func TestVerifyVC(t *testing.T) {
	vc := models.GenerateTestVC()

	success, err := VerifyVCSecp256k1(vc, &models.GenerateTestPrivKey().PublicKey)
	assert.Nil(t, err)
	assert.True(t, success)
	vc.Type = append(vc.Type, "Modified")
	success, err = VerifyVCSecp256k1(vc, &models.GenerateTestPrivKey().PublicKey)
	assert.Equal(t, errval.ErrJWSAuthentication, err)
	assert.False(t, success)
}