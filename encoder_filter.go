package encoder_filter


import (
        "fmt"
        "github.com/mozilla-services/heka/message"
        "github.com/pborman/uuid"
        . "github.com/mozilla-services/heka/pipeline"
)

type EncoderFilter struct {
        *EncoderFilterConfig
}

type EncoderFilterConfig struct {
        EncoderTag       string `toml:"encoder_filter_tag"`
        EncoderName   string `toml:"encoder"`
}

func (f *EncoderFilter) ConfigStruct() interface{} {
        return &EncoderFilterConfig{
        EncoderTag:       "encoder_filtered",
        }
}

func (f *EncoderFilter) Init(config interface{}) (err error) {
        f.EncoderFilterConfig = config.(*EncoderFilterConfig)

        if f.EncoderTag == "" {
            return fmt.Errorf(`An encoder_filtered value must be specified for the EncoderTag Field`)
        }

        if f.EncoderName == "" {
            return fmt.Errorf(`An encoder must be specified`)
        }

        return
}

func (f *EncoderFilter) Run(fr FilterRunner, h PluginHelper) (err error) {
        base_name := f.EncoderName
        full_name := fr.Name() + "-" + f.EncoderName
        encoder, ok := h.Encoder(base_name, full_name)
        if !ok {
            return fmt.Errorf("Encoder not found: %s", full_name)
        }


        var (
            tag string
            pack *PipelinePack
            e        error
            outBytes []byte
        )
        tag = f.EncoderTag

        for pack = range fr.InChan() {
            pack2, _ := h.PipelinePack(pack.MsgLoopCount)
            if pack2 == nil {
                fr.LogError(fmt.Errorf("exceeded MaxMsgLoops = %d",
                        h.PipelineConfig().Globals.MaxMsgLoops))
                break
            }

            if outBytes, e = encoder.Encode(pack); e != nil {
                    fr.LogError(fmt.Errorf("Error encoding message: %s", e))
            } else {
                if len(outBytes) > 0 {
                    tagField, _ := message.NewField("EncoderTag", tag, "")
                    pack2.Message.AddField(tagField)
                    pack2.Message.SetUuid(uuid.NewRandom())
                    pack2.Message.SetPayload(string(outBytes))
                    fr.Inject(pack2)
                }
            }
            fr.UpdateCursor(pack.QueueCursor)
            pack.Recycle(nil)
        }

    return
}

func init() {
    RegisterPlugin("EncoderFilter", func() interface{} {
        return new(EncoderFilter)
    })
}
