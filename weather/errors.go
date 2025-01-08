package weather

import "fmt"

func newLocationError(location string) error {
	return fmt.Errorf("location '%s' is outside US coverage area.\nThis API only supports US locations like:\n"+
		"- San Francisco, CA\n"+
		"- New York, NY\n"+
		"- Miami, FL\n"+
		"- Chicago, IL", location)
}
