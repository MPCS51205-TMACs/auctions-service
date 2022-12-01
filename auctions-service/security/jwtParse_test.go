package security

import (
	"testing"
	// "time"
)

func TestJwtParse(t *testing.T) {

	secret := "G+KbPeShVmYq3t6w9z$C&F)J@McQfTjW"
	exAdminToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI4NDcwMTVjOC1kODI4LTQwNGUtYjg3OC1lYThlNTRhMzk5ZDkiLCJpc3MiOiJ1c2VyLXNlcnZpY2UiLCJhdWQiOiJtcGNzNTEyMDUiLCJlbWFpbCI6Im1hdHRAbXBjcy5jb20iLCJuYW1lIjoibWF0dCIsImF1dGhvcml0aWVzIjpbIlJPTEVfVVNFUiIsIlJPTEVfQURNSU4iXSwiaWF0IjoxNjY5NDA4MDY2LCJleHAiOjE2NzIwMDAwNjZ9.N1x3fIBUz9CLDtabc9Lig6a4VFmRPdQaJwYX2Vabov0"
	exUserToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJlN2YxNDI0OC03MWM3LTQ5MGQtOWYzOC0yNDdiMjRmNzI4YWEiLCJpc3MiOiJ1c2VyLXNlcnZpY2UiLCJhdWQiOiJtcGNzNTEyMDUiLCJlbWFpbCI6Im1hdHRfQG1wY3MuY29tIiwibmFtZSI6Im1hdHQiLCJhdXRob3JpdGllcyI6WyJST0xFX1VTRVIiXSwiaWF0IjoxNjY5NDA4MDk1LCJleHAiOjE2NzIwMDAwOTV9.j0_T0boKVL0MMpTmI_xSUfc3M25MoWqeo-Sdg9fVelQ"

	jwtParser := NewJwtParser(secret)
	authenticator := NewAuthenticator(jwtParser)
	res1 := authenticator.IsAdmin(exAdminToken) // expect true
	res2 := authenticator.IsUser(exUserToken)   // expect true
	res3 := authenticator.IsAdmin(exUserToken)  // expect false
	res4 := authenticator.IsUser(exUserToken)   // expect true

	if res1 != true {
		t.Errorf("\nRan:%s\nExpected:%t\nGot:%t", "authenticator.IsAdmin(exAdminToken)", true, res1)
	}
	if res2 != true {
		t.Errorf("\nRan:%s\nExpected:%t\nGot:%t", "authenticator.IsAdmin(exAdminToken)", true, res2)
	}
	if res3 == true {
		t.Errorf("\nRan:%s\nExpected:%t\nGot:%t", "authenticator.IsAdmin(exAdminToken)", false, res3)
	}
	if res4 != true {
		t.Errorf("\nRan:%s\nExpected:%t\nGot:%t", "authenticator.IsAdmin(exAdminToken)", true, res4)
	}
}
