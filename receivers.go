package ghipu

import (
	"net/http"
	"net/url"
)

type ReceiversService struct {
	client *httpClient
}

func NewReceiversService(secret string, receiverID int) *ReceiversService {
	client := httpClient{
		client: &http.Client{},
		secret: secret,
		recid:  receiverID,
	}
	return &ReceiversService{&client}
}

// Payment returns informatin of the payment with the given id.
func (rs *ReceiversService) Receivers(r *Receiver) (*ReceiverCreateResponse, error) {
	resp, err := rs.client.PostForm("/receivers", r.Params())
	if err != nil {
		return nil, err
	}
	var rcr ReceiverCreateResponse
	if err := unmarshalJSON(resp.Body, &rcr); err != nil {
		return nil, err
	}
	return &rcr, nil
}

type Receiver struct {
	AdminFirstName        string `json:"admin_first_name"`
	AdminLastName         string `json:"admin_last_name"`
	AdminEmail            string `json:"admin_email"`
	CountryCode           string `json:"country_code"`
	BussinessIdentifier   string `json:"bussiness_identifier"`
	BussinessCategory     string `json:"bussiness_category"`
	BussinessName         string `json:"bussiness_name"`
	BussinessPhone        string `json:"bussiness_phone"`
	BussinessAddressLine1 string `json:"bussiness_address_line_1"`
	BussinessAddressLine2 string `json:"bussiness_address_line_2"`
	BussinessAddressLine3 string `json:"bussiness_address_line_3"`
	ContactFullName       string `json:"contact_full_name"`
	ContactJobTitle       string `json:"contact_job_title"`
	ContactEmail          string `json:"contact_email"`
	ContactPhone          string `json:"contact_phone"`
}

func (r *Receiver) Params() url.Values {
	return url.Values{
		"admin_first_name":         {r.AdminFirstName},
		"admin_last_name":          {r.AdminLastName},
		"admin_email":              {r.AdminEmail},
		"country_code":             {r.AdminFirstName},
		"bussiness_identifier":     {r.BussinessIdentifier},
		"bussiness_category":       {r.BussinessCategory},
		"bussiness_name":           {r.BussinessName},
		"bussiness_phone":          {r.BussinessPhone},
		"bussiness_address_line_1": {r.BussinessAddressLine1},
		"bussiness_address_line_2": {r.BussinessAddressLine2},
		"bussiness_address_line_3": {r.BussinessAddressLine3},
		"contact_full_name":        {r.ContactFullName},
		"contact_job_title":        {r.ContactJobTitle},
		"contact_email":            {r.ContactEmail},
		"contact_phone":            {r.ContactPhone},
	}
}

type ReceiverCreateResponse struct {
	ReceiverID string `json:"receiver_id"`
	Secret     string `json:"secret"`
}
