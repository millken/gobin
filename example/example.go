package example

import "github.com/millken/gobin"

type SearchRequest struct {
	query           string
	page_number     int32
	result_per_page int32
}

func (o *SearchRequest) MarshalBinary() (data []byte, err error) {
	sz := len(o.query) + 16
	data = make([]byte, sz)
	var i int
	i += gobin.MarshalString(o.query, data[i:])
	i += gobin.MarshalInt32(o.page_number, data[i:])
	i += gobin.MarshalInt32(o.result_per_page, data[i:])

	return data, nil
}

func (o *SearchRequest) UnmarshalBinary(data []byte) error {
	var (
		i, n int
		err  error
	)
	if o.query, i, err = gobin.UnmarshalString(data[n:]); err != nil {
		return err
	}
	n += i
	if o.page_number, i, err = gobin.UnmarshalInt32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.result_per_page, i, err = gobin.UnmarshalInt32(data[n:]); err != nil {
		return err
	}
	n += i

	return nil
}

type SearchResponse struct {
	results string
}

func (o *SearchResponse) MarshalBinary() (data []byte, err error) {
	sz := len(o.results) + 8
	data = make([]byte, sz)
	var i int
	i += gobin.MarshalString(o.results, data[i:])

	return data, nil
}

func (o *SearchResponse) UnmarshalBinary(data []byte) error {
	var (
		i, n int
		err  error
	)
	if o.results, i, err = gobin.UnmarshalString(data[n:]); err != nil {
		return err
	}
	n += i

	return nil
}
