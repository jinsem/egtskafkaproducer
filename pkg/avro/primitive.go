// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     measurementpackage.avsc
 */

package avro

import (
	"fmt"
	"github.com/actgardner/gogen-avro/vm/types"
	"io"
	"math"
)

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type StringWriter interface {
	WriteString(string) (int, error)
}

func encodeFloat(w io.Writer, byteCount int, bits uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}
	for i := 0; i < byteCount; i++ {
		if bw != nil {
			err = bw.WriteByte(byte(bits & 255))
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(bits&255))
		}
		bits = bits >> 8
	}
	if bw == nil {
		_, err = w.Write(bb)
		return err
	}
	return nil
}

func encodeInt(w io.Writer, byteCount int, encoded uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	// To avoid reallocations, grow capacity to the largest possible size
	// for this integer
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}

	if encoded == 0 {
		if bw != nil {
			err = bw.WriteByte(0)
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(0))
		}
	} else {
		for encoded > 0 {
			b := byte(encoded & 127)
			encoded = encoded >> 7
			if !(encoded == 0) {
				b |= 128
			}
			if bw != nil {
				err = bw.WriteByte(b)
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, b)
			}
		}
	}
	if bw == nil {
		_, err := w.Write(bb)
		return err
	}
	return nil

}

func writeAnalogSensor(r *AnalogSensor, w io.Writer) error {
	var err error
	err = writeInt(r.SensorNumber, w)
	if err != nil {
		return err
	}
	err = writeInt(r.Value, w)
	if err != nil {
		return err
	}

	return nil
}

func writeArrayAnalogSensor(r []*AnalogSensor, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil || len(r) == 0 {
		return err
	}
	for _, e := range r {
		err = writeAnalogSensor(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeArrayLiquidSensor(r []*LiquidSensor, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil || len(r) == 0 {
		return err
	}
	for _, e := range r {
		err = writeLiquidSensor(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
}

func writeInt(r int32, w io.Writer) error {
	downShift := uint32(31)
	encoded := uint64((uint32(r) << 1) ^ uint32(r>>downShift))
	const maxByteSize = 5
	return encodeInt(w, maxByteSize, encoded)
}

func writeLiquidSensor(r *LiquidSensor, w io.Writer) error {
	var err error
	err = writeInt(r.SensorNumber, w)
	if err != nil {
		return err
	}
	err = writeString(r.ErrorFlag, w)
	if err != nil {
		return err
	}
	err = writeInt(r.ValueMillimeters, w)
	if err != nil {
		return err
	}
	err = writeInt(r.ValueLitres, w)
	if err != nil {
		return err
	}

	return nil
}

func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := uint64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}

func writeMeasurementPackage(r *MeasurementPackage, w io.Writer) error {
	var err error
	err = writeUnionNullString(r.Imei, w)
	if err != nil {
		return err
	}
	err = writeString(r.Guid, w)
	if err != nil {
		return err
	}
	err = writeLong(r.ClinetId, w)
	if err != nil {
		return err
	}
	err = writeLong(r.PacketID, w)
	if err != nil {
		return err
	}
	err = writeLong(r.MeasurementTimestamp, w)
	if err != nil {
		return err
	}
	err = writeLong(r.ReceivedTimestamp, w)
	if err != nil {
		return err
	}
	err = writeUnionNullDouble(r.Latitude, w)
	if err != nil {
		return err
	}
	err = writeUnionNullDouble(r.Longitude, w)
	if err != nil {
		return err
	}
	err = writeUnionNullInt(r.Speed, w)
	if err != nil {
		return err
	}
	err = writeUnionNullInt(r.Pdop, w)
	if err != nil {
		return err
	}
	err = writeUnionNullInt(r.Hdop, w)
	if err != nil {
		return err
	}
	err = writeUnionNullInt(r.Vdop, w)
	if err != nil {
		return err
	}
	err = writeUnionNullInt(r.NumOfSatelites, w)
	if err != nil {
		return err
	}
	err = writeNavigationSystem(r.NavigationSystem, w)
	if err != nil {
		return err
	}
	err = writeUnionNullInt(r.Direction, w)
	if err != nil {
		return err
	}
	err = writeUnionArrayAnalogSensorNull(r.AnalogSensors, w)
	if err != nil {
		return err
	}
	err = writeUnionArrayLiquidSensorNull(r.LiquidSensors, w)
	if err != nil {
		return err
	}

	return nil
}

func writeNavigationSystem(r NavigationSystem, w io.Writer) error {
	return writeInt(int32(r), w)
}

func writeNull(_ interface{}, _ io.Writer) error {
	return nil
}

func writeString(r string, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	if sw, ok := w.(StringWriter); ok {
		_, err = sw.WriteString(r)
	} else {
		_, err = w.Write([]byte(r))
	}
	return err
}

func writeUnionArrayAnalogSensorNull(r *UnionArrayAnalogSensorNull, w io.Writer) error {
	err := writeLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionArrayAnalogSensorNullTypeEnumArrayAnalogSensor:
		return writeArrayAnalogSensor(r.ArrayAnalogSensor, w)
	case UnionArrayAnalogSensorNullTypeEnumNull:
		return writeNull(r.Null, w)

	}
	return fmt.Errorf("invalid value for *UnionArrayAnalogSensorNull")
}

func writeUnionArrayLiquidSensorNull(r *UnionArrayLiquidSensorNull, w io.Writer) error {
	err := writeLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionArrayLiquidSensorNullTypeEnumArrayLiquidSensor:
		return writeArrayLiquidSensor(r.ArrayLiquidSensor, w)
	case UnionArrayLiquidSensorNullTypeEnumNull:
		return writeNull(r.Null, w)

	}
	return fmt.Errorf("invalid value for *UnionArrayLiquidSensorNull")
}

func writeUnionNullDouble(r *UnionNullDouble, w io.Writer) error {
	err := writeLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionNullDoubleTypeEnumNull:
		return writeNull(r.Null, w)
	case UnionNullDoubleTypeEnumDouble:
		return writeDouble(r.Double, w)

	}
	return fmt.Errorf("invalid value for *UnionNullDouble")
}

func writeUnionNullInt(r *UnionNullInt, w io.Writer) error {
	err := writeLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionNullIntTypeEnumNull:
		return writeNull(r.Null, w)
	case UnionNullIntTypeEnumInt:
		return writeInt(r.Int, w)

	}
	return fmt.Errorf("invalid value for *UnionNullInt")
}

func writeUnionNullString(r *UnionNullString, w io.Writer) error {
	err := writeLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionNullStringTypeEnumNull:
		return writeNull(r.Null, w)
	case UnionNullStringTypeEnumString:
		return writeString(r.String, w)

	}
	return fmt.Errorf("invalid value for *UnionNullString")
}

type ArrayLiquidSensorWrapper []*LiquidSensor

func (_ *ArrayLiquidSensorWrapper) SetBoolean(v bool)                { panic("Unsupported operation") }
func (_ *ArrayLiquidSensorWrapper) SetInt(v int32)                   { panic("Unsupported operation") }
func (_ *ArrayLiquidSensorWrapper) SetLong(v int64)                  { panic("Unsupported operation") }
func (_ *ArrayLiquidSensorWrapper) SetFloat(v float32)               { panic("Unsupported operation") }
func (_ *ArrayLiquidSensorWrapper) SetDouble(v float64)              { panic("Unsupported operation") }
func (_ *ArrayLiquidSensorWrapper) SetBytes(v []byte)                { panic("Unsupported operation") }
func (_ *ArrayLiquidSensorWrapper) SetString(v string)               { panic("Unsupported operation") }
func (_ *ArrayLiquidSensorWrapper) SetUnionElem(v int64)             { panic("Unsupported operation") }
func (_ *ArrayLiquidSensorWrapper) Get(i int) types.Field            { panic("Unsupported operation") }
func (_ *ArrayLiquidSensorWrapper) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *ArrayLiquidSensorWrapper) Finalize()                        {}
func (_ *ArrayLiquidSensorWrapper) SetDefault(i int)                 { panic("Unsupported operation") }
func (r *ArrayLiquidSensorWrapper) AppendArray() types.Field {
	var v *LiquidSensor
	v = NewLiquidSensor()

	*r = append(*r, v)
	return (*r)[len(*r)-1]
}

type ArrayAnalogSensorWrapper []*AnalogSensor

func (_ *ArrayAnalogSensorWrapper) SetBoolean(v bool)                { panic("Unsupported operation") }
func (_ *ArrayAnalogSensorWrapper) SetInt(v int32)                   { panic("Unsupported operation") }
func (_ *ArrayAnalogSensorWrapper) SetLong(v int64)                  { panic("Unsupported operation") }
func (_ *ArrayAnalogSensorWrapper) SetFloat(v float32)               { panic("Unsupported operation") }
func (_ *ArrayAnalogSensorWrapper) SetDouble(v float64)              { panic("Unsupported operation") }
func (_ *ArrayAnalogSensorWrapper) SetBytes(v []byte)                { panic("Unsupported operation") }
func (_ *ArrayAnalogSensorWrapper) SetString(v string)               { panic("Unsupported operation") }
func (_ *ArrayAnalogSensorWrapper) SetUnionElem(v int64)             { panic("Unsupported operation") }
func (_ *ArrayAnalogSensorWrapper) Get(i int) types.Field            { panic("Unsupported operation") }
func (_ *ArrayAnalogSensorWrapper) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *ArrayAnalogSensorWrapper) Finalize()                        {}
func (_ *ArrayAnalogSensorWrapper) SetDefault(i int)                 { panic("Unsupported operation") }
func (r *ArrayAnalogSensorWrapper) AppendArray() types.Field {
	var v *AnalogSensor
	v = NewAnalogSensor()

	*r = append(*r, v)
	return (*r)[len(*r)-1]
}
