package flightpath

import (
	"passenger-tracker/internal/models"
)

// Normalize will take the list of paths of airport and
// will return a path which is the final start and end to the travel.
func Normalize(interPaths []*models.Path) *models.Path {
	if len(interPaths) == 1 {
		return interPaths[0]
	}

	// Assume we have points: A, B, C, D, E
	// we need a route and path from A to E
	// we check for an airport code which is only at the start and another code which is only at the end.

	startPath := map[models.AirPortCode]struct{}{}
	endPath := map[models.AirPortCode]struct{}{}

	for _, path := range interPaths {
		if _, present := endPath[path.Start]; !present {
			startPath[path.Start] = struct{}{}
		} else {
			delete(startPath, path.Start)
		}

		if _, present := startPath[path.End]; !present {
			endPath[path.End] = struct{}{}
		} else {
			delete(endPath, path.End)
		}
	}

	if len(startPath) != 1 || len(endPath) != 1 {
		return nil
	}

	finalPath := &models.Path{}
	for start := range startPath {
		finalPath.Start = start
	}

	for end := range endPath {
		finalPath.End = end
	}

	return finalPath
}
