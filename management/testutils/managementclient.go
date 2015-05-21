// Package testutils contains some test utilities for the Azure SDK
package testutils

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/management"
)

// GetTestClient returns a management Client for testing. Expects
// AZSUBSCRIPTIONID and AZCERTDATA to be present in the environment. AZCERTDATA
// is the base64encoded binary representation of the PEM certificate data.
func GetTestClient(t *testing.T) management.Client {
	subid := os.Getenv("AZSUBSCRIPTIONID")
	certdata := os.Getenv("AZCERTDATA")
	if subid == "" || certdata == "" {
		t.Skip("AZSUBSCRIPTIONID or AZCERTDATA not set, skipping test")
	}
	cert, err := base64.StdEncoding.DecodeString(certdata)
	if err != nil {
		t.Fatal(err)
	}

	client, err := management.NewClient(subid, cert)
	if err != nil {
		t.Fatal(err)
	}
	return logger{client, t}
}

type logger struct {
	management.Client
	*testing.T
}

func chop(d []byte) string {
	const maxlen = 5000

	s := string(d)

	if len(s) > maxlen {
		return s[:maxlen] + "..."
	}
	return s
}

func (l logger) SendAzureGetRequest(url string) ([]byte, error) {
	d, err := l.Client.SendAzureGetRequest(url)
	l.T.Logf("AZURE> GET %s\n", url)
	if err != nil {
		l.T.Logf("   <<< ERROR: %+v\n", err)
	} else {
		l.T.Logf("   <<< %s\n", chop(d))
	}
	return d, err
}

func (l logger) SendAzurePostRequest(url string, data []byte) (management.OperationID, error) {
	oid, err := l.Client.SendAzurePostRequest(url, data)
	logOperation(l.T, "POST", url, data, oid, err)
	return oid, err
}

func (l logger) SendAzurePutRequest(url string, contentType string, data []byte) (management.OperationID, error) {
	oid, err := l.Client.SendAzurePutRequest(url, contentType, data)
	logOperation(l.T, "PUT", url, data, oid, err)
	return oid, err
}

func (l logger) SendAzureDeleteRequest(url string) (management.OperationID, error) {
	oid, err := l.Client.SendAzureDeleteRequest(url)
	logOperation(l.T, "DELETE", url, nil, oid, err)
	return oid, err
}

func logOperation(t *testing.T, method, url string, data []byte, oid management.OperationID, err error) {
	t.Logf("AZURE> %s %s\n", method, url)
	if data != nil {
		t.Logf("   >>> %s\n", chop(data))
	}
	if err != nil {
		t.Logf("   <<< ERROR: %+v\n", err)
	} else {
		t.Logf("   <<< OperationID: %s\n", oid)
	}
}
