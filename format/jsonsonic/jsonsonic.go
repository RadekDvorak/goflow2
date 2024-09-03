package jsonsonic

import (
	"github.com/bytedance/sonic"
	"github.com/netsampler/goflow2/v2/format"
	protoproducer "github.com/netsampler/goflow2/v2/producer/proto"
)

type JsonSonicDriver struct {
	Api sonic.API
}

func (d *JsonSonicDriver) Prepare() error {
	return nil
}

func (d *JsonSonicDriver) Init() error {
	return nil
}

func (d *JsonSonicDriver) Format(data interface{}) ([]byte, []byte, error) {
	var key []byte
	if dataIf, ok := data.(interface{ Key() []byte }); ok {
		key = dataIf.Key()
	}

	//fmt.Printf("Integer: %+v", data)
	//os.Exit(22)
	var output []byte
	var err error
	if m, ok := data.(*protoproducer.ProtoProducerMessage); ok {
		x := Wrapper{m}
		output, err = sonic.Marshal(x)
	} else {
		output, err = sonic.Marshal(data)
	}

	return key, output, err
}

func init() {
	d := &JsonSonicDriver{
		Api: sonic.ConfigFastest,
	}

	format.RegisterFormatDriver("jsonsonic", d)
}

// Wrapper is a wrapper for ProtoProducerMessage that allows custom JSON rendering
type Wrapper struct {
	*protoproducer.ProtoProducerMessage
}

// WrappedRendered contains hardcoded JSON properties for performant marshalling
type WrappedRendered struct {
	SrcAddr any `json:"SrcAddr"`
	DstAddr any `json:"DstAddr"`
}

func (d Wrapper) MarshalJSON() ([]byte, error) {
	temp := WrappedRendered{
		SrcAddr: d.renderField("SrcAddr", d.GetSrcAddr()),
		DstAddr: d.renderField("DstAddr", d.GetDstAddr()),
	}

	return sonic.Marshal(&temp)
}

func (d Wrapper) renderField(name string, data interface{}) interface{} {
	if r, ok := d.ProtoProducerMessage.Formatter().Render(name); ok {
		return r(d.ProtoProducerMessage, name, data)
	}

	return protoproducer.NilRenderer(nil, name, data)
}
