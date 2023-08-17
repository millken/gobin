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

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *BinPackage) MarshalBinary() (data []byte, err error) {
	sz := o.Type.Size() + len(o.Data) + len(o.Signature) + 20
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
		i, n int
		err  error
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
	Gyroscope     int32
	Accelerometer int32
	Random        string
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *SensorData) MarshalBinary() (data []byte, err error) {
	sz := len(o.Random) + 56
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
	if n, err = o.MarshalInt32(o.Gyroscope, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalInt32(o.Accelerometer, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalString(o.Random, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if offset != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, offset)
	}
	return data, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *SensorData) UnmarshalBinary(data []byte) error {
	var (
		i, n int
		err  error
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
	if o.Gyroscope, i, err = o.UnmarshalInt32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Accelerometer, i, err = o.UnmarshalInt32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Random, i, err = o.UnmarshalString(data[n:]); err != nil {
		return err
	}
	n += i

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

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *SensorConfig) MarshalBinary() (data []byte, err error) {
	sz := len(o.Firmware) + 33
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
		i, n int
		err  error
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

	return nil
}

type SensorState struct {
	gobin.Safe
	State uint32
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *SensorState) MarshalBinary() (data []byte, err error) {
	sz := 4
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
		i, n int
		err  error
	)
	if o.State, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i

	return nil
}

type SensorConfirm struct {
	gobin.Safe
	Owner string
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *SensorConfirm) MarshalBinary() (data []byte, err error) {
	sz := len(o.Owner) + 8
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
		i, n int
		err  error
	)
	if o.Owner, i, err = o.UnmarshalString(data[n:]); err != nil {
		return err
	}
	n += i

	return nil
}
