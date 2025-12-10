package zenotiv1

import "time"

func getClient() Client { //Fairlawn
	locationId := "VtYagR1iXZB7C52IvNay"
	centerId := "dc4f223d-495a-48a9-878b-a2f194349f9e"
	apiKey := "e594defe147c4e98a56aea07b81abdc98038296566884a9297f1d342b48c9584"
	host := hostProd

	return Client{
		cfg: config{
			locationId: locationId,
			centerId:   centerId,
			apiKey:     apiKey,
			created:    time.Now(),
			host:       host,
		},
		service: Service{},
	}
}

func getFairlawnClient() Client { //Fairlawn
	locationId := "TpwQvq1uDohQXHFebMQj"
	centerId := "5bcf46f1-3db2-418b-a100-b1477f9af7dc"
	apiKey := "893ccafc65d44472adca82df5f447217e513715d6de345ed9a46c2bd41e4168d"

	return Client{
		cfg: config{
			locationId: locationId,
			centerId:   centerId,
			apiKey:     apiKey,
			created:    time.Now(),
			host:       hostProd,
		},
		service: Service{},
	}
}

func getTrainingClient() Client {
	locationId := "TpwQvq1uDohQXHFebMQj"
	centerId := "cd06fc08-fe13-40d7-942b-d43f476a3d40"
	apiKey := "893ccafc65d44472adca82df5f447217e513715d6de345ed9a46c2bd41e4168d"

	return Client{
		cfg: config{
			locationId: locationId,
			centerId:   centerId,
			apiKey:     apiKey,
			created:    time.Now(),
			host:       hostProd,
		},
		service: Service{},
	}
}

func getTribecaClient() Client {
	locationId := "aTVEb9wmnz5wD53YMFwc"
	centerId := "abae4b6a-7eaa-49be-a615-b765f8ff2999"
	apiKey := "3ddc6aacd0d74a4fb5841eaf525e8dbb1f805af2e8db42b09dc635fd5ee6a264"

	return Client{
		cfg: config{
			locationId: locationId,
			centerId:   centerId,
			apiKey:     apiKey,
			created:    time.Now(),
		},
		service: Service{},
	}
}

func getNaplesClient() Client {
	locationId := "VTuXXGb0flmRx2YOE5O0"
	centerId := "540bad63-a729-4720-8fa9-29f017c82e74"
	apiKey := "12687dc090584bfc963e5b43ba36c7ea285fec1e8afa438797d803a86050df92"

	return Client{
		cfg: config{
			locationId: locationId,
			centerId:   centerId,
			apiKey:     apiKey,
			created:    time.Now(),
			host:       hostProd,
		},
		service: Service{},
	}
}

func getSkinlabClient() Client {
	locationId := "aTVEb9wmnz5wD53YMFwc"
	centerId := "48c67965-eb5a-4926-9c67-d816da3c266f"
	apiKey := "3ddc6aacd0d74a4fb5841eaf525e8dbb1f805af2e8db42b09dc635fd5ee6a264"

	return Client{
		cfg: config{
			locationId: locationId,
			centerId:   centerId,
			apiKey:     apiKey,
			created:    time.Now(),
		},
		service: Service{},
	}
}
