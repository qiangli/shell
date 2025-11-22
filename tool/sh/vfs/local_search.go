package vfs

// Recursively search for files and directories matching a pattern
// Parameters: path (required): Starting path for the search,
// pattern (required): Search pattern to match against file names
func (s *LocalFS) SearchFiles(path string, options *SearchOptions) (string, error) {
	if options == nil {
		options = &SearchOptions{}
	}
	validPath, err := s.validatePath(path)
	if err != nil {
		return "", err
	}

	return Search(validPath, options)
}
