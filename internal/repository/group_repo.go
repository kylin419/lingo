package repository

type GroupRepository interface{
	SaveLanguages(sourceID string,langs []string)error
	GetLanguages(sourceID string)([]string,error)
}