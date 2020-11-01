package main

type options map[string]string

func (o *options) getValues() (r []string) {
	for k := range *o {
		r = append(r, k)
	}
	return
}
func (o *options) getNames() (r []string) {
	for _, v := range *o {
		r = append(r, v)
	}
	return
}

func (o *options) getValueFromName(value string) (r string) {
	for k, v := range *o {
		if v == value {
			r = k
		}
	}
	return
}
