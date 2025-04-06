package models

type KafkaEventType string

const KafkaEventTypeCompanyCreate = "company.create"
const KafkaEventTypeCompanyGet = "company.get"
const KafkaEventTypeCompanyPatch = "company.patch"
const KafkaEventTypeCompanyDelete = "company.delete"

type KafkaEvent struct {
	Type string
	Data interface{}
}
