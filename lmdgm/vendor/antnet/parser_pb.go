package antnet

import (
	"github.com/golang/protobuf/proto"
)

type pBParser struct {
	factory *Parser
}

func (r *pBParser) ParseC2S(msg *Message) (IMsgParser, error) {
	if msg == nil {
		return nil, ErrPBUnPack
	}
	p, ok := r.factory.msgMap[msg.Head.CmdAct()]
	if ok {
		if p.C2S() != nil {
			err := PBUnPack(msg.Data, p.C2S())
			if err != nil {
				return nil, err
			}
			p.parser = r
			return &p, nil
		}
	}

	return nil, ErrPBUnPack
}

func (r *pBParser) PackMsg(v interface{}) []byte {
	data, _ := PBPack(v)
	return data
}
func (r *pBParser) GetRemindMsg(err error, t MsgType) *Message {
	if t == MsgTypeMsg {
		return NewErrMsg(err)
	} else {
		return NewStrMsg(err.Error() + "\n")
	}
}
func (r *pBParser) GetType() ParserType {
	return r.factory.Type
}
func (r *pBParser) GetErrType() ParseErrType {
	return r.factory.ErrType
}
func PBUnPack(data []byte, msg interface{}) error {
	if data == nil || msg == nil {
		return ErrPBUnPack
	}

	err := proto.Unmarshal(data, msg.(proto.Message))
	if err != nil {
		return ErrPBUnPack
	}
	return nil
}

func PBPack(msg interface{}) ([]byte, error) {
	if msg == nil {
		return nil, ErrPBPack
	}

	data, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		LogInfo("")
	}

	return data, nil
}
