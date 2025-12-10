package zenotiv1

import (
	"time"
)

type (
	Center struct {
		Id   string
		Name string
	}

	Sale struct {
		Invoice_no string
		Guest      struct {
			Guest_id   string
			Guest_name string
			Guest_code string
		}
		Sold_on     time.Time
		Serviced_on time.Time
		Item        struct {
			Type int
			Name string
			Code string
		}
	}

	Collection struct {
		Invoice_id       string
		Total_collection float64
		Guest_id         string
		Created_Date     ZenotiTime
		Guest            Guest
		Items            []struct {
			Final_sale_price float64
		}
	}

	Guest struct {
		Id            string
		Center_id     string `json:"center_id"`
		Personal_info Personal_info
		Address_info  Address_Info `json:"address_info"`
		Code          string
		Referral      ReferralInfo     `json:"referral"`
		Preferences   GuestPreferences `json:"preferences"`
		Tags          []string
		Gender        int
		DateOfBirth   ZenotiDate `json:"date_of_birth"`
	}

	GuestPreferences struct {
		Receive_Transactional_SMS   bool `json:"receive_transactional_sms"`
		Receive_Marketing_SMS       bool `json:"receive_marketing_sms"`
		Receive_Transactional_Email bool `json:"receive_transactional_email"`
		Receive_Marketing_Email     bool `json:"receive_marketing_email"`
	}

	ReferralInfo struct {
		Referral_source struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"referral_source"`
	}

	Address_Info struct {
		Address_1   string `json:"address_1"`
		Address_2   string `json:"address_2"`
		City        string `json:"city"`
		Country_id  int    `json:"country_id"`
		State_id    int    `json:"state_id"`
		State_other string `json:"state_other"`
		Zip_code    string `json:"zip_code"`
	}

	Personal_info struct {
		First_name   string     `json:"first_name"`
		Last_name    string     `json:"last_name"`
		Email        string     `json:"email"`
		Mobile_phone Phone_info `json:"mobile_phone"`
		Gender       int        `json:"gender"`
		DateOfBirth  ZenotiTime `json:"date_of_birth"`
	}

	Phone_info struct {
		Country_code int    `json:"country_code"`
		Number       string `json:"number"`
	}

	Appointment struct {
		Id         string `json:"appointment_id"`
		Invoice_id string
		Guest      struct {
			Id         string
			First_name string
			Last_name  string
			Mobile     struct {
				Number string `json:"display_number"`
			}
			Email string
		}
		Start_time     ZenotiTime
		End_time       ZenotiTime
		Start_time_utc ZenotiTime
		End_time_utc   ZenotiTime
		Price          struct {
			Sales float32
		}
		Creation_date ZenotiTime
		Status        ZenotiStatus
		Therapist     Therapist
		BlockOut      BlockOut
		Service       AppointmentServiceItem
	}

	AppointmentGroup struct {
		Appointment_Group_Id string
		Invoice_id           string
		Invoice_status       InvoiceStatus
		Invoice              Invoice
		Appointment_Services []AppointmentService
		Appointment_Packages []AppointmentPackage

		Price struct {
			Sales float32
		}
	}

	AppointmentGroupWebhookData struct {
		Appointment_Group_Id string `json:"appointment_group_id"`
		Invoice_id           string `json:"invoice_id"`
		Center_Id            string `json:"center_id"`
		Guest                AppointmentGroupGuest
		Appointments         []AppointmentGroupAppointment
	}

	AppointmentGroupAppointment struct {
		Appointment_Id       string
		Start_time           ZenotiTime
		End_time             ZenotiTime
		Start_time_in_center ZenotiTime
		End_time_in_center   ZenotiTime
		Status               ZenotiStatus
		Service_id           string
		Service_Name         string
	}

	AppointmentGroupGuest struct {
		Id        string
		FirstName string
		LastName  string
		Email     string
	}

	AppointmentService struct {
		Appointment_Id     string
		Invoice_item_id    string
		Start_time         ZenotiTime
		End_time           ZenotiTime
		Appointment_status ZenotiStatus
		Service            AppointmentServiceItem
		Quantity           int
		Room               Room
	}

	AppointmentServiceItem struct {
		Id           string
		Name         string
		Display_name string
		Price        struct {
			Sales float32
		}
	}

	AppointmentPackage struct {
		Id    string
		Name  string
		Price struct {
			Sales float32
		}
		Service              AppointmentService
		Appointment_services []AppointmentService
	}

	Booking struct {
		Id    string
		Error struct {
			StatusCode int
			Message    string
		}
	}

	ReservationInvoice struct {
		Invoice_id string
		Guest      Guest
	}

	Reservation struct {
		Reservation_id string `json:"reservation_id"`
		Expiry_time    ZenotiTime
		Invoices       []ReservationInvoice
		Invoice        ReservationInvoice
		Error          struct {
			StatusCode int
			Message    string
		}
	}

	Invoice struct {
		Invoice struct {
			Id             string
			Invoice_number string
			Invoice_date   ZenotiTime
			Total_price    struct {
				Currency_id         uint
				Net_price           float64
				Tax                 float64
				Rounding_adjustment float64
				Total_price         float64
				Sum_total           float64
			}
		} `json:"invoice"`
		Guest struct {
			Id           string
			First_name   string
			Last_name    string
			Mobile_phone string
		}
		Appointments []Appointment
	}

	ZenotiTime struct {
		Time time.Time `jsont:"time,omitempty"`
	}

	ZenotiDate struct {
		Time time.Time `jsont:"time,omitempty"`
	}

	GuestFull struct {
		Id            string
		Code          string
		Center_id     string
		Personal_info Personal_info
	}

	Member struct {
		Id       string
		Guest_id Guest
	}

	Therapist struct {
		Id           string
		First_name   string
		Last_name    string
		Display_name string
	}

	BlockOut struct {
		Block_out_time_id   string
		Start_time          ZenotiTime
		End_time            ZenotiTime
		Employee            Therapist
		Notes               string
		Block_out_time_type struct {
			Name string
		}
	}

	Service struct {
		Id          string
		Name        string
		Description string
		Duration    float32
	}

	Room struct {
		Id          string
		Name        string
		Code        string
		Description string
	}

	PageInfo struct {
		Total int
		Page  int
		Size  int
	}

	// Webhook definitions
	InvoiceWebhookData struct {
		Id           string
		Invoice_Date ZenotiTime
		Center_Id    string
		Total_Price  struct {
			Sum_Total float64
		}

		Guest GuestWebhookData
	}

	GuestWebhookData struct {
		Id           string
		First_Name   string
		Last_Name    string
		Mobile_Phone string
		Email        string
	}

	WebhookDataPayload struct {
		Invoice InvoiceWebhookData
	}

	WebhookData struct {
		Id         string
		Event_type string
		Data       WebhookDataPayload
	}

	PaymentOption struct {
		Id       string
		CardLogo string `json:"card_logo"`
		LastFour string `json:"last_four"`
		ExpiryOn string `json:"expiry_on"`
	}
)

var (
	DefaultReferral = ReferralInfo{
		Referral_source: struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}{
			Id:   "e3b71e3a-28e2-4458-a374-6fc556c83475",
			Name: "Internet",
		},
	}
)

func (t *ZenotiTime) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		t.Time = time.Now()
		return
	}

	date, err := time.Parse(`"2006-01-02T15:04:05"`, string(b))
	if err == nil {
		t.Time = date
		return
	}
	date, err = time.Parse(`"2006-01-02"`, string(b))
	if err == nil {
		t.Time = date
		return
	}

	return err
}

func (t ZenotiTime) MarshalJSON() ([]byte, error) {
	parsed := t.Time.Format("\"2006-01-02T15:04:05\"")
	return []byte(parsed), nil
}

func (t *ZenotiDate) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		t.Time = time.Now()
		return
	}

	date, err := time.Parse(`"2006-01-02"`, string(b))
	if err == nil {
		t.Time = date
		return
	}
	date, err = time.Parse(`"2006-01-02"`, string(b))
	if err == nil {
		t.Time = date
		return
	}

	return err
}

func (t ZenotiDate) MarshalJSON() ([]byte, error) {
	parsed := t.Time.Format("\"2006-01-02\"")
	return []byte(parsed), nil
}

func GetDefaultReferral() ReferralInfo {
	return ReferralInfo{
		struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}{
			Id:   "e3b71e3a-28e2-4458-a374-6fc556c83475",
			Name: "Internet",
		},
	}
}
