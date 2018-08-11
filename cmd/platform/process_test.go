package platform

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
	"testing"
)



func TestParse(t *testing.T) {

	//tests := []struct {
	//	name string
	//	args string
	//	want string
	//}{
	//
	//}

	data := []byte(`
platforms:
 cf:
    modules:
    - native-type: html5
      platform-type: "javascript.nodejs"
    - native-type: html5
      platform-type: "javascript.nodejs"
 xsa:
    modules:
    - native-type: html5
      platform-type: "javascript.nodejs"
    - native-type: html5
      platform-type: "javascript.nodejs"
`)
	y := Platforms{}

	err := yaml.Unmarshal([]byte(data), &y)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("%+v\n", y)

	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		if got := ; got != tt.want {
	//			t.Errorf("SetMtaProp() = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}

func main() {

}
