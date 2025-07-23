package validator

import "regexp"

//Declare a regular expression for sanity checking the format of email addresses

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)


//New validator type which contains a map of validation errors
type Validator struct{
	Errors map[string]string
}

//New is a helper which creates a new validator instance with an empty errors map
func New() *Validator{
	return &Validator{Errors: make(map[string]string)}
}

//Valid returns true if the errors map doesn't contain any entries
func (v *Validator)Valid()bool{
	return len(v.Errors)==0
}

//AddError adds an error message to the map (as long as no entry exists for the given key)
func (v *Validator) AddError(key,message string){
	_,exists:=v.Errors[key]
	if !exists{
		v.Errors[key]=message
	}
}

//Check adds an error message to the map only if a validation check is not ok
func(v *Validator)Check(ok bool,key, message string){
	if !ok{
		v.AddError(key,message)
	}
}

//If the number of unique value appears twice
//the map just keeps one of them
//if the number of unique items (in the map) is equal to the original slice length, all
//values were unique
func Unique[T comparable](values []T)bool{
	uniqueValues:=make(map[T]bool)

	for _,value:=range values{
		uniqueValues[value]=true
	}
	 return len(uniqueValues)==len(values)
}