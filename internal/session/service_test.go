package session

import (
	"testing"

)
func TestStartSession(t *testing.T){
	service:=NewService(nil)
	token,err:=service.StartSession(1)
	if err != nil{
		t.Errorf("expected no error,but got:%v",err)
	}
	if token==""{
		t.Errorf("Expected a token,but got an empty string")
	}
}
func TestValidateSession(t *testing.T){
	service:=NewService(nil)
	_,err:=service.ValidateSession("1234")
	if err==nil{
		t.Errorf("Expected an error for fake token, but got none")
	}
}