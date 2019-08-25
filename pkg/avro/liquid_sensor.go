// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     measurementpackage.avsc
 */

package avro

import (
	"github.com/actgardner/gogen-avro/compiler"
	"github.com/actgardner/gogen-avro/container"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/vm/types"
	"io"
)

type LiquidSensor struct {

	// Ordinal number of sensor
	SensorNumber int32

	// Value is set if error occurred during reading value from liquid sensor
	ErrorFlag string

	// This value is set if liquid sensor measures level of the liqid in a vessel. Otherwise value is 0
	ValueMillimeters int32

	// This value is set if liquid sensor measures volume of the liqid in a vessel. Otherwise value is 0
	ValueLitres int32
}

func NewLiquidSensorWriter(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error) {
	str := &LiquidSensor{}
	return container.NewWriter(writer, codec, recordsPerBlock, str.Schema())
}

func DeserializeLiquidSensor(r io.Reader) (*LiquidSensor, error) {
	t := NewLiquidSensor()

	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	return t, err
}

func NewLiquidSensor() *LiquidSensor {
	return &LiquidSensor{}
}

func (r *LiquidSensor) Schema() string {
	return "{\"fields\":[{\"doc\":\"Ordinal number of sensor\",\"name\":\"sensorNumber\",\"type\":\"int\"},{\"doc\":\"Value is set if error occurred during reading value from liquid sensor\",\"name\":\"errorFlag\",\"type\":\"string\"},{\"doc\":\"This value is set if liquid sensor measures level of the liqid in a vessel. Otherwise value is 0\",\"name\":\"valueMillimeters\",\"type\":\"int\"},{\"doc\":\"This value is set if liquid sensor measures volume of the liqid in a vessel. Otherwise value is 0\",\"name\":\"valueLitres\",\"type\":\"int\"}],\"name\":\"LiquidSensor\",\"namespace\":\"com.projavt.ecology.schema\",\"type\":\"record\"}"
}

func (r *LiquidSensor) SchemaName() string {
	return "com.projavt.ecology.schema.LiquidSensor"
}

func (r *LiquidSensor) Serialize(w io.Writer) error {
	return writeLiquidSensor(r, w)
}

func (_ *LiquidSensor) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *LiquidSensor) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *LiquidSensor) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *LiquidSensor) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *LiquidSensor) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *LiquidSensor) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *LiquidSensor) SetString(v string)   { panic("Unsupported operation") }
func (_ *LiquidSensor) SetUnionElem(v int64) { panic("Unsupported operation") }
func (r *LiquidSensor) Get(i int) types.Field {
	switch i {
	case 0:
		return (*types.Int)(&r.SensorNumber)
	case 1:
		return (*types.String)(&r.ErrorFlag)
	case 2:
		return (*types.Int)(&r.ValueMillimeters)
	case 3:
		return (*types.Int)(&r.ValueLitres)

	}
	panic("Unknown field index")
}
func (r *LiquidSensor) SetDefault(i int) {
	switch i {

	}
	panic("Unknown field index")
}
func (_ *LiquidSensor) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *LiquidSensor) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *LiquidSensor) Finalize()                        {}

type LiquidSensorReader struct {
	r io.Reader
	p *vm.Program
}

func NewLiquidSensorReader(r io.Reader) (*LiquidSensorReader, error) {
	containerReader, err := container.NewReader(r)
	if err != nil {
		return nil, err
	}

	t := NewLiquidSensor()
	deser, err := compiler.CompileSchemaBytes([]byte(containerReader.AvroContainerSchema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	return &LiquidSensorReader{
		r: containerReader,
		p: deser,
	}, nil
}

func (r *LiquidSensorReader) Read() (*LiquidSensor, error) {
	t := NewLiquidSensor()
	err := vm.Eval(r.r, r.p, t)
	return t, err
}
