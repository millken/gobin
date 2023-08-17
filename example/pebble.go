package example

import (
	"fmt"

	"github.com/millken/gobin"
)

// PackageType is the type of package
type PackageType uint16

const (
	// DATA is a data package
	PackageType_DATA   PackageType = 0
	PackageType_CONFIG PackageType = 1
	PackageType_STATE  PackageType = 2
)

func (o *PackageType) Size() int {
	return 2
}

// MarshalTo writes a wire-format message to w.
func (o *PackageType) MarshalTo(w []byte) (int, error) {
	w[0] = byte(*o)
	w[1] = byte(*o >> 8)
	return 2, nil
}

// UnmarshalTo reads a wire-format message from data.
func (o *PackageType) UnmarshalTo(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, fmt.Errorf("invalid data size %d", len(data))
	}
	*o = PackageType(uint16(data[0]) | uint16(data[1])<<8)
	return 2, nil
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *PackageType) MarshalBinary() ([]byte, error) {
	data := make([]byte, 2)
	_, err := o.MarshalTo(data)
	return data, err
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *PackageType) UnmarshalBinary(data []byte) error {
	_, err := o.UnmarshalTo(data)
	return err
}

type BinPackage struct {
	gobin.Safe
	Type      *PackageType
	Data      []byte
	Timestamp uint32
	Signature []byte
}

func (o *BinPackage) Size() int {
	var sz int

	sz += o.Type.Size()
	sz += len(o.Data)

	sz += len(o.Signature)
	sz += 20
	return sz
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *BinPackage) MarshalBinary() (data []byte, err error) {
	sz := o.Size()
	data = make([]byte, sz)
	var offset, n int
	if n, err = o.Type.MarshalTo(data); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalBytes(o.Data, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.Timestamp, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalBytes(o.Signature, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if offset != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, offset)
	}
	return data, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *BinPackage) UnmarshalBinary(data []byte) error {
	var (
		i, n, l int
		err     error
	)
	if i, err = o.Type.UnmarshalTo(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Data, i, err = o.UnmarshalBytes(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Timestamp, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Signature, i, err = o.UnmarshalBytes(data[n:]); err != nil {
		return err
	}
	n += i

	_ = l
	return nil
}

type SensorData struct {
	gobin.Safe
	Snr           uint32
	Vbat          uint32
	Latitude      int32
	Longitude     int32
	GasResistance uint32
	Temperature   int32
	Pressure      uint32
	Humidity      uint32
	Light         uint32
	Temperature2  uint32
	Gyroscope     []int32
	Accelerometer []int32
	Random        []string
}

func (o *SensorData) Size() int {
	var sz int

	sz += len(o.Gyroscope) * 4
	sz += len(o.Accelerometer) * 4
	for _, v := range o.Random {
		sz = sz + len(v) + 8
	}
	sz += 64
	return sz
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *SensorData) MarshalBinary() (data []byte, err error) {
	sz := o.Size()
	data = make([]byte, sz)
	var offset, n int
	if n, err = o.MarshalUint32(o.Snr, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.Vbat, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalInt32(o.Latitude, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalInt32(o.Longitude, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.GasResistance, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalInt32(o.Temperature, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.Pressure, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.Humidity, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.Light, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.Temperature2, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalInt(len(o.Gyroscope), data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	for _, v := range o.Gyroscope {
		if n, err = o.MarshalInt32(v, data[offset:]); err != nil {
			return nil, err
		}
		offset += n
	}
	if n, err = o.MarshalInt(len(o.Accelerometer), data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	for _, v := range o.Accelerometer {
		if n, err = o.MarshalInt32(v, data[offset:]); err != nil {
			return nil, err
		}
		offset += n
	}
	if n, err = o.MarshalInt(len(o.Random), data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	for _, v := range o.Random {
		if n, err = o.MarshalString(v, data[offset:]); err != nil {
			return nil, err
		}
		offset += n
	}
	if offset != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, offset)
	}
	return data, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *SensorData) UnmarshalBinary(data []byte) error {
	var (
		i, n, l int
		err     error
	)
	if o.Snr, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Vbat, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Latitude, i, err = o.UnmarshalInt32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Longitude, i, err = o.UnmarshalInt32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.GasResistance, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Temperature, i, err = o.UnmarshalInt32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Pressure, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Humidity, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Light, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Temperature2, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if l, i, err = o.UnmarshalInt(data[n:]); err != nil {
		return err
	}
	n += i
	o.Gyroscope = make([]int32, l)
	for j := range o.Gyroscope {
		if v, m, err := o.UnmarshalInt32(data[n:]); err != nil {
			return err
		} else {
			i = m
			o.Gyroscope[j] = v
		}
		n += i
	}
	if l, i, err = o.UnmarshalInt(data[n:]); err != nil {
		return err
	}
	n += i
	o.Accelerometer = make([]int32, l)
	for j := range o.Accelerometer {
		if v, m, err := o.UnmarshalInt32(data[n:]); err != nil {
			return err
		} else {
			i = m
			o.Accelerometer[j] = v
		}
		n += i
	}
	if l, i, err = o.UnmarshalInt(data[n:]); err != nil {
		return err
	}
	n += i
	o.Random = make([]string, l)
	for j := range o.Random {
		if v, m, err := o.UnmarshalString(data[n:]); err != nil {
			return err
		} else {
			i = m
			o.Random[j] = v
		}
		n += i
	}

	_ = l
	return nil
}

type SensorConfig struct {
	gobin.Safe
	BulkUpload             uint32
	DataChannel            uint32
	UploadPeriod           uint32
	BulkUploadSamplingCnt  uint32
	BulkUploadSamplingFreq uint32
	Beep                   uint32
	Firmware               string
	DeviceConfigurable     bool
}

func (o *SensorConfig) Size() int {
	var sz int

	sz += len(o.Firmware)
	sz += 33
	return sz
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *SensorConfig) MarshalBinary() (data []byte, err error) {
	sz := o.Size()
	data = make([]byte, sz)
	var offset, n int
	if n, err = o.MarshalUint32(o.BulkUpload, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.DataChannel, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.UploadPeriod, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.BulkUploadSamplingCnt, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.BulkUploadSamplingFreq, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.Beep, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalString(o.Firmware, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalBool(o.DeviceConfigurable, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if offset != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, offset)
	}
	return data, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *SensorConfig) UnmarshalBinary(data []byte) error {
	var (
		i, n, l int
		err     error
	)
	if o.BulkUpload, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.DataChannel, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.UploadPeriod, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.BulkUploadSamplingCnt, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.BulkUploadSamplingFreq, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Beep, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Firmware, i, err = o.UnmarshalString(data[n:]); err != nil {
		return err
	}
	n += i
	if o.DeviceConfigurable, i, err = o.UnmarshalBool(data[n:]); err != nil {
		return err
	}
	n += i

	_ = l
	return nil
}

type SensorState struct {
	gobin.Safe
	State uint32
}

func (o *SensorState) Size() int {
	var sz int
	sz += 4
	return sz
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *SensorState) MarshalBinary() (data []byte, err error) {
	sz := o.Size()
	data = make([]byte, sz)
	var offset, n int
	if n, err = o.MarshalUint32(o.State, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if offset != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, offset)
	}
	return data, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *SensorState) UnmarshalBinary(data []byte) error {
	var (
		i, n, l int
		err     error
	)
	if o.State, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i

	_ = l
	return nil
}

type SensorConfirm struct {
	gobin.Safe
	Owner string
}

func (o *SensorConfirm) Size() int {
	var sz int

	sz += len(o.Owner)
	sz += 8
	return sz
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *SensorConfirm) MarshalBinary() (data []byte, err error) {
	sz := o.Size()
	data = make([]byte, sz)
	var offset, n int
	if n, err = o.MarshalString(o.Owner, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if offset != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, offset)
	}
	return data, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *SensorConfirm) UnmarshalBinary(data []byte) error {
	var (
		i, n, l int
		err     error
	)
	if o.Owner, i, err = o.UnmarshalString(data[n:]); err != nil {
		return err
	}
	n += i

	_ = l
	return nil
}
