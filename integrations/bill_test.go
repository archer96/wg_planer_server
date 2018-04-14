package integrations

import (
	"net/http"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/wgplaner/wg_planer_server/models"
)

func TestCreateBill(t *testing.T) {
	prepareTestEnv(t)
	authValid := "1234567890fakefirebaseid0001"
	req := NewRequest(t, "POST", authValid, "/group/00112233-4455-6677-8899-aabbccddeeff/bills/create")
	MakeRequest(t, req, http.StatusOK)
}

func TestGetBillUnauthorizedGroup(t *testing.T) {
	prepareTestEnv(t)
	authValid := "1234567890fakefirebaseid0001"
	groupUnauthorized := "00112233-4455-6677-8899-aabbccddeef0"
	req := NewRequest(t, "GET", authValid, "/group/"+groupUnauthorized+"/bills")
	MakeRequest(t, req, http.StatusUnauthorized)
}

func TestGetBills(t *testing.T) {
	prepareTestEnv(t)
	var (
		billList  = models.BillList{}
		authValid = "1234567890fakefirebaseid0001"
		req       = NewRequest(t, "GET", authValid, "/group/00112233-4455-6677-8899-aabbccddeeff/bills")
		resp      = MakeRequest(t, req, http.StatusOK)
	)
	DecodeJSON(t, resp, &billList)

	assert.Len(t, billList.Bills, 1)
	assert.Equal(t, int64(1), billList.Count)
	assert.Len(t, billList.Bills[0].BillItems, 2)
	assert.Len(t, billList.Bills[0].SentTo, 2)
	assert.Equal(t, "todo", *billList.Bills[0].State)
	assert.Equal(t, int64(370), billList.Bills[0].Sum)
	assert.Equal(t, strfmt.UUID("00112233-4455-6677-8899-aabbccddeeff"), billList.Bills[0].GroupUID)
}
