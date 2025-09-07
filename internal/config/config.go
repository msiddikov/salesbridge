package config

import (
	"strings"
	"time"

	"github.com/Lavina-Tech-LLC/lavinagopackage/v2/conf"
)

var (
	Confs   Conf
	testing = false
)

type (
	Conf struct {
		Locations []Location
		Settings  Settings
		RC        RingCentral
		DB        gormDB
		GAjson    string
	}
	Settings struct {
		SrvAddress   string
		SrvDomain    string
		Cert         string
		Key          string
		ClientId     string
		ClientSecret string
		Token        string
		RefreshToken string

		CRAgencyAPI string
	}

	RingCentral struct {
		ClientId     string
		ClientSecret string
		ServerURL    string
		Username     string
		Extension    string
		Password     string
		SyncDate     time.Time
		Locations    []string
	}

	gormDB struct {
		Host     string
		Port     string
		User     string
		Password string
		DbName   string
	}

	Location struct {
		Name       string
		Id         string
		RWID       string
		Api        string
		Pipeline   string
		PipelineId string
		Book       string
		BookId     string
		Sales      string
		SalesId    string
		Leads      string
		LeadsId    string
		NoShows    string
		NoShowsId  string
		Zenoti     struct {
			Api                 string
			CenterId            string
			Url                 string
			ReferralId          string
			Appointments        []string
			CollectionsSyncDay  time.Time
			AppointmentsSyncDay time.Time
		}
	}
)

func init() {
	arg := "conf/"
	if strings.Contains(conf.GetPath(), "tests") {
		arg = "../conf/"
	}
	Confs = conf.Get[Conf](arg)
}

func SetSalesSyncDate(l Location, t time.Time) {
	for k, loc := range Confs.Locations {
		if l.Id == loc.Id {
			Confs.Locations[k].Zenoti.CollectionsSyncDay = t
			SaveConf()
			return
		}
	}
}

func GetSalesSyncDate(l Location) time.Time {
	for _, loc := range Confs.Locations {
		if l.Id == loc.Id {
			return loc.Zenoti.CollectionsSyncDay
		}
	}
	return time.Now()
}

func SetTesting() {
	testing = true
}

func IsTesting() bool {
	return testing
}

func GetLocations() []Location {
	return Confs.Locations
}

func GetLocationById(id string) Location {
	for _, l := range Confs.Locations {
		if l.Id == id {
			return l
		}
	}
	return Location{}
}
func GetLocationByRWID(id string) Location {
	for _, l := range Confs.Locations {
		if l.RWID == id {
			return l
		}
	}
	return Location{}
}

func UpdateLocations(l []Location) {
	Confs.Locations = l
	SaveConf()
}

func SaveConf() {
	conf.Set(Confs)
}
