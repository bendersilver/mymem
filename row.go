package mymem

func (r *Rows) Values() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	return append([]string{r.Key}, r.values...), nil
}

func (r *Rows) Map() (map[string]string, error) {
	if !r.containers {
		return nil, ErrContainers
	}
	if r.err != nil {
		return nil, r.err
	}
	return append([]string{r.Key}, r.values...), nil
}
