package testdata

import (
	"client-runaway-zenoti/internal/types"
	"encoding/json"
	"fmt"
)

func Appointments() []types.Appointment {
	data := `[{
        "appointment_id": "bf2c1177-f110-434d-b6df-28fd018aba43",
        "appointment_segment_id": null,
        "parent_service_name": null,
        "appointment_group_id": "bf2c1177-f110-434d-b6df-28fd018aba43",
        "invoice_id": "00000000-0000-0000-0000-000000000000",
        "service": null,
        "start_time": "2022-06-01T15:00:00",
        "start_time_utc": "2022-06-01T19:00:00",
        "end_time": "2022-06-01T16:00:00",
        "end_time_utc": "2022-06-01T20:00:00",
        "status": 10,
        "source": 0,
        "progress": 0,
        "locked": false,
        "invoice_locked": false,
        "has_active_membership_for_auto_pay": false,
        "auto_pay_authorize_status": 0,
        "has_unexpired_packages": false,
        "guest": null,
        "therapist": {
            "id": "8a8d6f7d-dc49-480c-a4aa-fa2f0ef5baf9",
            "first_name": "Callie",
            "last_name": "Greene",
            "nick_name": null,
            "display_name": null,
            "email": "cgreene.vioburlington@gmail.com",
            "gender": 0,
            "vanity_image_url": ""
        },
        "room": null,
        "equipment": null,
        "service_custom_data_indicator": null,
        "notes": "Mock Day\nV/O Signature Facial w/ dermaplaning\nModel: Karen Russo",
        "price": {
            "currency_id": 0,
            "sales": 0.0,
            "tax": 0.0,
            "final": 0.0,
            "final1": 0.0,
            "discount": 0.0,
            "tip": 0.0,
            "ssg": null,
            "rounding_correction": 0.0
        },
        "actual_start_time": null,
        "actual_completed_time": null,
        "checkin_time": null,
        "therapist_preference_type": null,
        "form_id": null,
        "blockout": {
            "id": 1309,
            "name": "Meeting  ",
            "code": "Meeting  ",
            "description": null,
            "duration": 60,
            "active": true,
            "count_in_utilization": true,
            "bookable": false,
            "indicator_color": "#8D8DCF",
            "text_color": "#000000"
        },
        "creation_date": "2022-05-18T21:20:00",
        "creation_date_utc": "2022-05-19T01:20:00",
        "created_by_id": null,
        "closed_by_id": null,
        "show_in_calender": 1,
        "email_link": null,
        "sms_link": null,
        "appointment_category_id": null,
        "invoice_processed_in_integrations": 0,
        "parallel_group_id": null,
        "available_therapists": null,
        "available_rooms": null,
        "available_times": null,
        "selected_therapist_id": null,
        "selected_room_id": null,
        "selected_time": "0001-01-01T00:00:00",
        "requested_alternative": 0,
        "group_invoice_id": null,
        "group_name": null,
        "canUpdateTherapist": false,
        "package_id": "00000000-0000-0000-0000-000000000000",
        "virtual_room_link": null,
        "error": null
    },
    {
        "appointment_id": "d3a8d58d-539a-45a1-b12e-dbf46070d1aa",
        "appointment_segment_id": null,
        "parent_service_name": "Lifting Thread - Rejuvenation Lifting",
        "appointment_group_id": "e891318c-5267-4008-b03d-dc1cef59bce3",
        "invoice_id": "abe82145-513d-453f-a92a-bd97f16faf66",
        "service": {
            "id": "79dd666b-c5e4-4d88-a159-6ee9ce170769",
            "segment_id": "00000000-0000-0000-0000-000000000000",
            "name": "Lifting Thread - Rejuvenation Lifting",
            "is_addon": false,
            "has_addons": true,
            "parent_appointment_id": null,
            "business_unit": {
                "guid": "7a081b11-f8fc-408c-bc62-5dbd965699b2",
                "name": "Default",
                "id": 210
            },
            "category": {
                "id": "06d120a5-e120-4357-abf5-f3491bae6183",
                "name": "Injectables"
            },
            "sub_category": {
                "id": "9b6945a8-587f-478d-b75a-70b001f41e67",
                "name": "Threads"
            },
            "override_product_consumption": false,
            "override_default_product_consumption": false,
            "is_virtual_service": false
        },
        "start_time": "2022-06-01T15:30:00",
        "start_time_utc": "2022-06-01T19:30:00",
        "end_time": "2022-06-01T16:30:00",
        "end_time_utc": "2022-06-01T20:30:00",
        "status": 0,
        "source": 0,
        "progress": 0,
        "locked": false,
        "invoice_locked": false,
        "has_active_membership_for_auto_pay": true,
        "auto_pay_authorize_status": -3,
        "has_unexpired_packages": false,
        "guest": {
            "id": "d8a9f935-ae60-4a51-a0c0-87c1173a39cd",
            "first_name": "Linnea",
            "last_name": "Laverty",
            "gender": 0,
            "mobile": {
                "country_id": 0,
                "number": null,
                "display_number": "+1 2076894988"
            },
            "email": "neanea126@yahoo.com",
            "indicator": "0@0@0@0@0@0@0@x@0@0@0@0@0#0@0@0@0",
            "lp_tier_info": "0@x",
            "is_virtual_user": false,
            "GuestIndicatorValue": {
                "HighSpender": null,
                "Member": 0,
                "LowFeedback": null,
                "RegularGuest": null,
                "FirstTimer": null,
                "ReturningCustomer": null,
                "NoShow": null,
                "HasActivePackages": null,
                "HasProfileAlerts": null,
                "OtherCenterGuest": null,
                "HasCTA": null,
                "Dues": null,
                "CardOnFile": null,
                "AutoPayEnabled": null,
                "RecurrenceAppointment": null,
                "RebookedAppointment": null,
                "hasAddOns": true,
                "LpTier": null,
                "IsSurpriseVisit": null,
                "CustomDataIndicator": null,
                "IsGuestBirthday": null
            }
        },
        "therapist": {
            "id": "80e77bca-9743-439d-8cf2-dd4af579e02a",
            "first_name": "Melissa",
            "last_name": "Lake",
            "nick_name": null,
            "display_name": null,
            "email": "mlake.vioburlington@gmail.com",
            "gender": 0,
            "vanity_image_url": ""
        },
        "room": null,
        "equipment": null,
        "service_custom_data_indicator": "1#1#0#0#0#1#0",
        "notes": "Prepaid at Grand Opening Event. SO",
        "price": {
            "currency_id": 0,
            "sales": 0.0,
            "tax": 0.0,
            "final": 0.0,
            "final1": 0.0,
            "discount": 0.0,
            "tip": 0.0,
            "ssg": null,
            "rounding_correction": 0.0
        },
        "actual_start_time": null,
        "actual_completed_time": null,
        "checkin_time": null,
        "therapist_preference_type": 0,
        "form_id": null,
        "blockout": null,
        "creation_date": "2022-05-19T17:46:00",
        "creation_date_utc": "2022-05-19T21:46:00",
        "created_by_id": "3bb483df-a44f-418a-ba7d-db4ebe0754c5",
        "closed_by_id": null,
        "show_in_calender": 1,
        "email_link": null,
        "sms_link": null,
        "appointment_category_id": null,
        "invoice_processed_in_integrations": 0,
        "parallel_group_id": null,
        "available_therapists": null,
        "available_rooms": null,
        "available_times": null,
        "selected_therapist_id": null,
        "selected_room_id": null,
        "selected_time": "0001-01-01T00:00:00",
        "requested_alternative": 0,
        "group_invoice_id": null,
        "group_name": null,
        "canUpdateTherapist": true,
        "package_id": "00000000-0000-0000-0000-000000000000",
        "virtual_room_link": null,
        "error": null
    },
    {
        "appointment_id": "6711cd9d-26eb-48b4-9e69-04279821c092",
        "appointment_segment_id": null,
        "parent_service_name": null,
        "appointment_group_id": "6711cd9d-26eb-48b4-9e69-04279821c092",
        "invoice_id": "00000000-0000-0000-0000-000000000000",
        "service": null,
        "start_time": "2022-06-01T16:00:00",
        "start_time_utc": "2022-06-01T20:00:00",
        "end_time": "2022-06-01T17:00:00",
        "end_time_utc": "2022-06-01T21:00:00",
        "status": 10,
        "source": 0,
        "progress": 0,
        "locked": false,
        "invoice_locked": false,
        "has_active_membership_for_auto_pay": false,
        "auto_pay_authorize_status": 0,
        "has_unexpired_packages": false,
        "guest": null,
        "therapist": {
            "id": "8a8d6f7d-dc49-480c-a4aa-fa2f0ef5baf9",
            "first_name": "Callie",
            "last_name": "Greene",
            "nick_name": null,
            "display_name": null,
            "email": "cgreene.vioburlington@gmail.com",
            "gender": 0,
            "vanity_image_url": ""
        },
        "room": null,
        "equipment": null,
        "service_custom_data_indicator": null,
        "notes": "Mock Day\nMicroneedling\nModel: Save for Callie's friend\n",
        "price": {
            "currency_id": 0,
            "sales": 0.0,
            "tax": 0.0,
            "final": 0.0,
            "final1": 0.0,
            "discount": 0.0,
            "tip": 0.0,
            "ssg": null,
            "rounding_correction": 0.0
        },
        "actual_start_time": null,
        "actual_completed_time": null,
        "checkin_time": null,
        "therapist_preference_type": null,
        "form_id": null,
        "blockout": {
            "id": 1309,
            "name": "Meeting  ",
            "code": "Meeting  ",
            "description": null,
            "duration": 60,
            "active": true,
            "count_in_utilization": true,
            "bookable": false,
            "indicator_color": "#8D8DCF",
            "text_color": "#000000"
        },
        "creation_date": "2022-05-18T21:20:00",
        "creation_date_utc": "2022-05-19T01:20:00",
        "created_by_id": null,
        "closed_by_id": null,
        "show_in_calender": 1,
        "email_link": null,
        "sms_link": null,
        "appointment_category_id": null,
        "invoice_processed_in_integrations": 0,
        "parallel_group_id": null,
        "available_therapists": null,
        "available_rooms": null,
        "available_times": null,
        "selected_therapist_id": null,
        "selected_room_id": null,
        "selected_time": "0001-01-01T00:00:00",
        "requested_alternative": 0,
        "group_invoice_id": null,
        "group_name": null,
        "canUpdateTherapist": false,
        "package_id": "00000000-0000-0000-0000-000000000000",
        "virtual_room_link": null,
        "error": null
    },
    {
        "appointment_id": "a03d30e7-5772-495e-899d-3939c5a38b60",
        "appointment_segment_id": null,
        "parent_service_name": "HydraFacial",
        "appointment_group_id": "03a14ae9-0f6d-4152-902d-cdab891bd42c",
        "invoice_id": "e812a818-ff2f-4848-ad63-b6a21131b51a",
        "service": {
            "id": "f57bc93b-7e5c-472d-aa03-b8de3e044cfb",
            "segment_id": "00000000-0000-0000-0000-000000000000",
            "name": "HydraFacial",
            "is_addon": false,
            "has_addons": true,
            "parent_appointment_id": null,
            "business_unit": {
                "guid": "7a081b11-f8fc-408c-bc62-5dbd965699b2",
                "name": "Default",
                "id": 210
            },
            "category": {
                "id": "54087b1b-c6ed-4009-aac7-84242f1b0768",
                "name": "Spa Services      "
            },
            "sub_category": {
                "id": "bb71cfcb-bb63-4644-bd0c-369530652b58",
                "name": "HydraFacial"
            },
            "override_product_consumption": false,
            "override_default_product_consumption": false,
            "is_virtual_service": false
        },
        "start_time": "2022-06-01T16:00:00",
        "start_time_utc": "2022-06-01T20:00:00",
        "end_time": "2022-06-01T17:00:00",
        "end_time_utc": "2022-06-01T21:00:00",
        "status": 0,
        "source": 0,
        "progress": 0,
        "locked": false,
        "invoice_locked": false,
        "has_active_membership_for_auto_pay": true,
        "auto_pay_authorize_status": -3,
        "has_unexpired_packages": false,
        "guest": {
            "id": "7c843fad-d4f7-454c-9531-663425a5b702",
            "first_name": "Tara",
            "last_name": "Saieh",
            "gender": 0,
            "mobile": {
                "country_id": 0,
                "number": null,
                "display_number": "+1 6173863120"
            },
            "email": "tara.chalvire@gmail.com",
            "indicator": "0@0@0@0@1@0@0@x@0@0@0@0@0#0@0@0@0",
            "lp_tier_info": "0@x",
            "is_virtual_user": false,
            "GuestIndicatorValue": {
                "HighSpender": null,
                "Member": 0,
                "LowFeedback": null,
                "RegularGuest": null,
                "FirstTimer": true,
                "ReturningCustomer": null,
                "NoShow": null,
                "HasActivePackages": null,
                "HasProfileAlerts": null,
                "OtherCenterGuest": null,
                "HasCTA": null,
                "Dues": null,
                "CardOnFile": null,
                "AutoPayEnabled": null,
                "RecurrenceAppointment": null,
                "RebookedAppointment": null,
                "hasAddOns": true,
                "LpTier": null,
                "IsSurpriseVisit": null,
                "CustomDataIndicator": null,
                "IsGuestBirthday": null
            }
        },
        "therapist": {
            "id": "fc0a76e1-12a9-45a9-a3d6-f87f2cd8cf63",
            "first_name": "Linda",
            "last_name": "Babigian",
            "nick_name": null,
            "display_name": null,
            "email": "lbabigian.vioburlington@gmail.com",
            "gender": 0,
            "vanity_image_url": ""
        },
        "room": {
            "id": "2f2852c3-0439-4805-8974-5274ac6a7a8a",
            "name": "Spa Services 5"
        },
        "equipment": {
            "id": "3f6814ad-01e5-45ae-980d-77956e5bfe01",
            "name": "HydraFacial"
        },
        "service_custom_data_indicator": "1#1#0#0#0#1#0",
        "notes": "this apointment is for the daughter ",
        "price": {
            "currency_id": 0,
            "sales": 200.0,
            "tax": 0.0,
            "final": 0.0,
            "final1": 0.0,
            "discount": 0.0,
            "tip": 0.0,
            "ssg": null,
            "rounding_correction": 0.0
        },
        "actual_start_time": null,
        "actual_completed_time": null,
        "checkin_time": null,
        "therapist_preference_type": 3,
        "form_id": null,
        "blockout": null,
        "creation_date": "2022-05-04T18:44:00",
        "creation_date_utc": "2022-05-04T22:44:00",
        "created_by_id": "b94a83cc-ac9b-4a49-931f-1b184204dc4a",
        "closed_by_id": null,
        "show_in_calender": 1,
        "email_link": null,
        "sms_link": null,
        "appointment_category_id": null,
        "invoice_processed_in_integrations": 0,
        "parallel_group_id": null,
        "available_therapists": null,
        "available_rooms": null,
        "available_times": null,
        "selected_therapist_id": null,
        "selected_room_id": null,
        "selected_time": "0001-01-01T00:00:00",
        "requested_alternative": 0,
        "group_invoice_id": null,
        "group_name": null,
        "canUpdateTherapist": true,
        "package_id": "00000000-0000-0000-0000-000000000000",
        "virtual_room_link": null,
        "error": null
    },
    {
        "appointment_id": "4493e803-b2d8-4f0c-914e-3b48ed082e8c",
        "appointment_segment_id": null,
        "parent_service_name": null,
        "appointment_group_id": "4493e803-b2d8-4f0c-914e-3b48ed082e8c",
        "invoice_id": "00000000-0000-0000-0000-000000000000",
        "service": null,
        "start_time": "2022-06-01T17:00:00",
        "start_time_utc": "2022-06-01T21:00:00",
        "end_time": "2022-06-01T18:00:00",
        "end_time_utc": "2022-06-01T22:00:00",
        "status": 10,
        "source": 0,
        "progress": 0,
        "locked": false,
        "invoice_locked": false,
        "has_active_membership_for_auto_pay": false,
        "auto_pay_authorize_status": 0,
        "has_unexpired_packages": false,
        "guest": null,
        "therapist": {
            "id": "8a8d6f7d-dc49-480c-a4aa-fa2f0ef5baf9",
            "first_name": "Callie",
            "last_name": "Greene",
            "nick_name": null,
            "display_name": null,
            "email": "cgreene.vioburlington@gmail.com",
            "gender": 0,
            "vanity_image_url": ""
        },
        "room": null,
        "equipment": null,
        "service_custom_data_indicator": null,
        "notes": "Mock Day\nDeLuxe HydraFacial\nModel: Pam Miller",
        "price": {
            "currency_id": 0,
            "sales": 0.0,
            "tax": 0.0,
            "final": 0.0,
            "final1": 0.0,
            "discount": 0.0,
            "tip": 0.0,
            "ssg": null,
            "rounding_correction": 0.0
        },
        "actual_start_time": null,
        "actual_completed_time": null,
        "checkin_time": null,
        "therapist_preference_type": null,
        "form_id": null,
        "blockout": {
            "id": 1309,
            "name": "Meeting  ",
            "code": "Meeting  ",
            "description": null,
            "duration": 60,
            "active": true,
            "count_in_utilization": true,
            "bookable": false,
            "indicator_color": "#8D8DCF",
            "text_color": "#000000"
        },
        "creation_date": "2022-05-18T21:21:00",
        "creation_date_utc": "2022-05-19T01:21:00",
        "created_by_id": null,
        "closed_by_id": null,
        "show_in_calender": 1,
        "email_link": null,
        "sms_link": null,
        "appointment_category_id": null,
        "invoice_processed_in_integrations": 0,
        "parallel_group_id": null,
        "available_therapists": null,
        "available_rooms": null,
        "available_times": null,
        "selected_therapist_id": null,
        "selected_room_id": null,
        "selected_time": "0001-01-01T00:00:00",
        "requested_alternative": 0,
        "group_invoice_id": null,
        "group_name": null,
        "canUpdateTherapist": false,
        "package_id": "00000000-0000-0000-0000-000000000000",
        "virtual_room_link": null,
        "error": null
    }]`

	result := []types.Appointment{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		fmt.Println(err)
	}

	return result
}

func Collections() []types.Collection {
	data := `[
        {
            "invoice_id": "8f14ea37-dce8-4bc2-913a-aa0f36889f8a",
            "invoice_no": "2777",
            "reciept_no": "R201907241",
            "created_date": "2019-07-24T00:00:00",
            "status": 4,
            "total_collection": 1131.35,
            "gross_amount": 1131.35,
            "net_amount": 1131.35,
            "discount": 0.0,
            "rounding_adjustment": 0.0,
            "cashback": 0.0,
            "guest_id": "d8a9f935-ae60-4a51-a0c0-87c1173a39cd",
            "items": [
                {
                    "id": "69fab682-1c16-43ba-8f13-ef12b998bc31",
                    "name": "60 Min Hot Stone Custom Massage Session",
                    "code": "60HSCMS",
                    "type": "Service",
                    "quantity": 1,
                    "final_sale_price": 1076.35,
                    "discount": 0.0,
                    "cashback_redemption": 0.0,
                    "therapist_id": "d833bde1-ace2-45fd-923c-85c873ce7600",
                    "taxes": [
                        {
                            "type": "Sales Tax",
                            "amount": 48.925,
                            "tax_percentage": 5.0,
                            "item_percentage": 100.0
                        },
                        {
                            "type": "Luxury Tax",
                            "amount": 48.925,
                            "tax_percentage": 5.0,
                            "item_percentage": 100.0
                        }
                    ],
                    "payments": [
                        {
                            "type": "CC",
                            "detail_type": "Amex",
                            "amount": 95.1385,
                            "tax": 8.649,
                            "tip": 0.0,
                            "ssg": 0.0
                        },
                        {
                            "type": "CC",
                            "detail_type": "Visa",
                            "amount": 95.1385,
                            "tax": 8.649,
                            "tip": 0.0,
                            "ssg": 0.0
                        },
                        {
                            "type": "LP",
                            "detail_type": null,
                            "amount": 95.1385,
                            "tax": 8.649,
                            "tip": 0.0,
                            "ssg": 0.0
                        },
                        {
                            "type": "CC",
                            "detail_type": "Mastercard",
                            "amount": 95.1386,
                            "tax": 8.649,
                            "tip": 0.0,
                            "ssg": 0.0
                        },
                        {
                            "type": "custom",
                            "detail_type": null,
                            "amount": 95.1386,
                            "tax": 8.649,
                            "tip": 0.0,
                            "ssg": 0.0
                        },
                        {
                            "type": "CASH",
                            "detail_type": null,
                            "amount": 600.6573,
                            "tax": 54.6052,
                            "tip": 220.0,
                            "ssg": 0.0
                        }
                    ]
                },
                {
                    "id": "60f9d3fe-8b8f-4341-a730-22ee5445b782",
                    "name": "MOROCCANOIL Soften & Shine Set",
                    "code": "400200001",
                    "type": "Product",
                    "quantity": 1,
                    "final_sale_price": 55.0,
                    "discount": 0.0,
                    "cashback_redemption": 0.0,
                    "therapist_id": "46cd4a55-4ded-4c42-bccb-421755a84845",
                    "taxes": [
                        {
                            "type": "Sales Tax",
                            "amount": 2.5,
                            "tax_percentage": 5.0,
                            "item_percentage": 100.0
                        },
                        {
                            "type": "Luxury Tax",
                            "amount": 2.5,
                            "tax_percentage": 5.0,
                            "item_percentage": 100.0
                        }
                    ],
                    "payments": [
                        {
                            "type": "CC",
                            "detail_type": "Mastercard",
                            "amount": 4.8614,
                            "tax": 0.4419,
                            "tip": 0.0,
                            "ssg": 0.0
                        },
                        {
                            "type": "custom",
                            "detail_type": null,
                            "amount": 4.8614,
                            "tax": 0.4419,
                            "tip": 0.0,
                            "ssg": 0.0
                        },
                        {
                            "type": "CC",
                            "detail_type": "Amex",
                            "amount": 4.8615,
                            "tax": 0.442,
                            "tip": 0.0,
                            "ssg": 0.0
                        },
                        {
                            "type": "CC",
                            "detail_type": "Visa",
                            "amount": 4.8615,
                            "tax": 0.442,
                            "tip": 0.0,
                            "ssg": 0.0
                        },
                        {
                            "type": "LP",
                            "detail_type": null,
                            "amount": 4.8615,
                            "tax": 0.442,
                            "tip": 0.0,
                            "ssg": 0.0
                        },
                        {
                            "type": "CASH",
                            "detail_type": null,
                            "amount": 30.6927,
                            "tax": 2.7902,
                            "tip": 0.0,
                            "ssg": 0.0
                        }
                    ]
                }
            ]
        },
        {
            "invoice_id": "692f33f3-f3cb-4204-a799-d6797a9d5194",
            "invoice_no": "2778",
            "reciept_no": "",
            "created_date": "2019-07-24T00:00:00",
            "status": 0,
            "total_collection": 55.0,
            "gross_amount": 55.0,
            "net_amount": 55.0,
            "discount": 0.0,
            "rounding_adjustment": 0.0,
            "cashback": 0.0,
            "guest_id": "7c843fad-d4f7-454c-9531-663425a5b702",
            "items": [
                {
                    "id": "60f9d3fe-8b8f-4341-a730-22ee5445b782",
                    "name": "MOROCCANOIL Soften & Shine Set",
                    "code": "400200001",
                    "type": "Product",
                    "quantity": 1,
                    "final_sale_price": 55.0,
                    "discount": 0.0,
                    "cashback_redemption": 0.0,
                    "therapist_id": "46cd4a55-4ded-4c42-bccb-421755a84845",
                    "taxes": [
                        {
                            "type": "Sales Tax",
                            "amount": 2.5,
                            "tax_percentage": 5.0,
                            "item_percentage": 100.0
                        },
                        {
                            "type": "Luxury Tax",
                            "amount": 2.5,
                            "tax_percentage": 5.0,
                            "item_percentage": 100.0
                        }
                    ],
                    "payments": [
                        {
                            "type": "CASH",
                            "detail_type": null,
                            "amount": 55.0,
                            "tax": 5.0,
                            "tip": 0.0,
                            "ssg": 0.0
                        }
                    ]
                }
            ]
        }
    ]`

	result := []types.Collection{}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		fmt.Println(err)
	}

	return result
}
