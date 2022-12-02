package security

import (
	"fmt"
	"testing"
	// "time"
)

func TestJwtParse(t *testing.T) {

	secret := "G+KbPeShVmYq3t6w9z$C&F)J@McQfTjW"
	exAdminToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI4NDcwMTVjOC1kODI4LTQwNGUtYjg3OC1lYThlNTRhMzk5ZDkiLCJpc3MiOiJ1c2VyLXNlcnZpY2UiLCJhdWQiOiJtcGNzNTEyMDUiLCJlbWFpbCI6Im1hdHRAbXBjcy5jb20iLCJuYW1lIjoibWF0dCIsImF1dGhvcml0aWVzIjpbIlJPTEVfVVNFUiIsIlJPTEVfQURNSU4iXSwiaWF0IjoxNjY5NDA4MDY2LCJleHAiOjE2NzIwMDAwNjZ9.N1x3fIBUz9CLDtabc9Lig6a4VFmRPdQaJwYX2Vabov0"
	exUserToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJlN2YxNDI0OC03MWM3LTQ5MGQtOWYzOC0yNDdiMjRmNzI4YWEiLCJpc3MiOiJ1c2VyLXNlcnZpY2UiLCJhdWQiOiJtcGNzNTEyMDUiLCJlbWFpbCI6Im1hdHRfQG1wY3MuY29tIiwibmFtZSI6Im1hdHQiLCJhdXRob3JpdGllcyI6WyJST0xFX1VTRVIiXSwiaWF0IjoxNjY5NDA4MDk1LCJleHAiOjE2NzIwMDAwOTV9.j0_T0boKVL0MMpTmI_xSUfc3M25MoWqeo-Sdg9fVelQ"
	mattToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJjMDM4MmQyYy1iYTVlLTQxZDgtYTYzYi00YjMzNDlkN2Q4YTMiLCJhdWQiOiJtcGNzNTEyMDUiLCJpc3MiOiJ1c2VyLXNlcnZpY2UiLCJuYW1lIjoibWF0dCIsImV4cCI6MTY3MDA1Nzc4NSwiaWF0IjoxNjcwMDE0NTg1LCJlbWFpbCI6Im1hdHRAbXBjcy5jb20iLCJhdXRob3JpdGllcyI6WyJST0xFX1VTRVIiLCJST0xFX0FETUlOIl19.tPsUVZex3evg-ZE7e2vUbRt7XHYhABMoodsqy0FYohA"
	illToken := "asdfasdf"
	aBidder := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiI4MmRhMzUxOC0xMDQ0LTRlOWItODQ1Yy1kZGMyYmI4MmViNGQiLCJhdWQiOiJtcGNzNTEyMDUiLCJpc3MiOiJ1c2VyLXNlcnZpY2UiLCJuYW1lIjoiTGF5VXN5RTBPa1piQUtWIiwiZXhwIjoxNjcwMDYwMDYwLCJpYXQiOjE2NzAwMTY4NjAsImVtYWlsIjoiTGF5VXN5RTBPa1piQUtWQG1wY3MuY29tIiwiYXV0aG9yaXRpZXMiOlsiUk9MRV9VU0VSIl19.jyQf2_XJNWnLjCm8sQJ-H39eFIuX45U1FbC9JuFFTEw"

	jwtParser := NewJwtParser(secret)
	authenticator := NewAuthenticator(jwtParser)
	res1 := authenticator.IsAdmin(exAdminToken) // expect true
	res2 := authenticator.IsUser(exUserToken)   // expect true
	res3 := authenticator.IsAdmin(exUserToken)  // expect false
	res4 := authenticator.IsUser(exUserToken)   // expect true
	res5 := authenticator.IsAdmin(illToken)     // expect false
	res6 := authenticator.IsUser(illToken)      // expect false
	res7 := authenticator.IsAdmin(illToken)     // expect false
	res8 := authenticator.IsUser(illToken)      // expect false

	fmt.Println(res5)
	fmt.Println(res6)
	fmt.Println(res7)
	fmt.Println(res8)

	fmt.Println(authenticator.IsAdmin(aBidder))
	fmt.Println(authenticator.IsUser(aBidder))

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

	res9, gotErr := authenticator.ExtractUserId(mattToken) // expect true
	fmt.Println(res9, gotErr)

}
