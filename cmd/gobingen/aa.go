package main

import "fmt"

func (o *AAA) Size() int {
	sz := 0
	sz += len(o.A1)
	sz += len(o.A3) * 4
	sz += len(o.Sub.Name)
	sz += len(o.Sub.Sub3.B3)
	sz += 57
	return sz
}

func (o *AAA) MarshalTo(data []byte) (int, error) {
	var offset, n int
	var err error
	if n, err = o.MarshalString(o.A1, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	if n, err = o.MarshalInt(o.A2, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	if n, err = o.MarshalInt(len(o.A3), data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	for _, v := range o.A3 {
		if n, err = o.MarshalFloat32(v, data[offset:]); err != nil {
			return 0, err
		}
		offset += n
	}
	if n, err = o.MarshalFloat64(o.A4, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	if n, err = o.MarshalString(o.Sub.Name, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	if n, err = o.MarshalInt(o.Sub.Sub3.Id, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	if n, err = o.MarshalBytes(o.Sub.Sub3.B3, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	return offset, nil
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *AAA) MarshalBinary() (data []byte, err error) {
	sz := o.Size()
	data = make([]byte, sz)
	n, err := o.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	if n != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, n)
	}
	return data, nil
}

func (o *AAA) UnmarshalTo(data []byte) (int, error) {
	var (
		i, n, l int
		err     error
	)
	if o.A1, i, err = o.UnmarshalString(data[n:]); err != nil {
		return 0, err
	}
	n += i
	if o.A2, i, err = o.UnmarshalInt(data[n:]); err != nil {
		return 0, err
	}
	n += i
	if l, i, err = o.UnmarshalInt(data[n:]); err != nil {
		return 0, err
	}
	n += i
	o.A3 = make([]float32, l)
	for j := range o.A3 {
		if o.A3[j], i, err = o.UnmarshalFloat32(data[n:]); err != nil {
			return 0, err
		}
		n += i
	}
	if o.A4, i, err = o.UnmarshalFloat64(data[n:]); err != nil {
		return 0, err
	}
	n += i
	if o.Sub.Name, i, err = o.UnmarshalString(data[n:]); err != nil {
		return 0, err
	}
	n += i
	if o.Sub.Sub3.Id, i, err = o.UnmarshalInt(data[n:]); err != nil {
		return 0, err
	}
	n += i
	if o.Sub.Sub3.B3, i, err = o.UnmarshalBytes(data[n:]); err != nil {
		return 0, err
	}
	n += i
	_ = l
	return n, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *AAA) UnmarshalBinary(data []byte) error {
	_, err := o.UnmarshalTo(data)
	return err
}
