package main

import "fmt"

type Human struct {
	name     string
	id       int
	location string
}

type validator struct {
	errors map[string]string
}

func new() *validator {
	return &validator{errors: make(map[string]string)}
}

func (v *validator) valid() bool {
	return len(v.errors) == 0
}

func (v *validator) adderror(key, message string) {
	_, exists := v.errors[key]
	fmt.Println("value against a given key:",message)
	fmt.Println("is it there?:",exists)

	if !exists{
		v.errors[key]=message
	}
}

func (v *validator)check(ok bool,key,message string){
	if !ok{
		v.adderror(key,message)
	}
}


func main(){
	//instantiate a new struct and validate fields
	input:=&Human{
		name: "",
		id: 23,
		location: "kisuma",
	}

	v:=new()

	v.check(len(input.name)>=1,"nameError","Name can not be empty")
	v.check(input.id>=50,"IntError","Id can not be less than 50")
	v.check(input.location=="kanairo","locationErr","Location must be Kisuma")

	if !v.valid(){
		fmt.Println("sorry we can not proceed with your request!")
		fmt.Println("Here are the errors:",v)
		return
	}

	fmt.Println("congratulations mate..you passed our test!!")

	




}