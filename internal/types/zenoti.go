package types

import (
	"fmt"
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
		Guest            GuestWithId
		Items            []struct {
			Final_sale_price float64
		}
	}

	Guest struct {
		Id            string `json:"guest_id"`
		Personal_info Personal_info
		Code          string
	}

	GuestWithId struct {
		Id            string `json:"id"`
		Personal_info Personal_info
		Code          string
		Referral      ReferralInfo `json:"referral"`
	}

	ReferralInfo struct {
		Referral_source struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"referral_source"`
	}

	Personal_info struct {
		First_name   string     `json:"first_name"`
		Last_name    string     `json:"last_name"`
		Email        string     `json:"email"`
		Mobile_phone Phone_info `json:"mobile_phone"`
		Gender       int        `json:"gender"`
		DateOfBirth  ZenotiTime `json:"date_of_birth"`
		//User_name    string     `json:"user_name"`
		//Middle_name  string     `json:"middle_name"`
		//Work_phone       Phone_info `json:"work_phone"`
		//Home_phone       Phone_info `json:"home_phone"`
		//Is_minor         bool       `json:"is_minor"`
		//Nationality_id   int        `json:"nationality_id"`
		//Anniversary_date ZenotiTime `json:"anniversary_date"`
		//Lock_guest       bool       `json:"lock_guest_custom_data"`
		//Pan              string
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
		Start_time ZenotiTime
		End_time   ZenotiTime
		Price      struct {
			Sales float32
		}
		Creation_date ZenotiTime
		Status        ZenotiStatus
		Therapist     Therapist
	}

	ZenotiStatus int

	ZenotiTime struct {
		time.Time
	}

	GuestFull struct {
		Id            string
		Code          string
		Center_id     string
		Personal_info Personal_info
	}

	Member struct {
		Id       string
		Guest_id GuestWithId
	}

	Therapist struct {
		Id           string
		First_name   string
		Last_name    string
		Display_name string
	}

	BlockSlot struct {
		Block_out_time_id string
		Start_time        ZenotiTime
		End_time          ZenotiTime
		Employee          Therapist
		Notes             string
	}
)

const (
	NoShowed  ZenotiStatus = -2
	Canceled  ZenotiStatus = -1
	Booked    ZenotiStatus = 0
	Closed    ZenotiStatus = 1
	CheckedIn ZenotiStatus = 2
	Confirmed ZenotiStatus = 4
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
	parsed := fmt.Sprintf("%s", t.Time.Format("\"2006-01-02\""))
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
