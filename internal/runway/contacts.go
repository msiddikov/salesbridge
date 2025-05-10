package runway

import runwayv2 "client-runaway-zenoti/packages/runwayV2"

func GetContactByIdLocationId(locId, contactId string) (contact runwayv2.Contact, err error) {
	client, err := svc.NewClientFromId(locId)
	if err != nil {
		return
	}

	contact, err = client.ContactsGet(contactId)
	return
}
