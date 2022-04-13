package models

import pb "passenger-tracker/api/pb/v1alpha1/flightpath"

// Path defines the object which has the beginning and ending airport codes.
type Path struct {
	Start AirPortCode
	End   AirPortCode
}

func (path *Path) ConvertToAPIObject() *pb.Path {
	if path == nil {
		return nil
	}

	return &pb.Path{
		Start: string(path.Start),
		End:   string(path.End),
	}
}

func (path *Path) ConvertFromAPIObject(apiObj *pb.Path) {
	if apiObj == nil {
		return
	}

	if path == nil {
		path = &Path{}
	}

	path.Start = AirPortCode(apiObj.Start)
	path.End = AirPortCode(apiObj.End)
}

// AirPortCode represents the code of airport stations.
type AirPortCode string

// Validate return true if the airport code is valid and returns false if code is invalid.
func (code AirPortCode) Validate() bool {
	if len(code) != 3 {
		return false
	}

	return true
}
