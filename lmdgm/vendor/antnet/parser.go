package antnet

import (
	"reflect"
)

type IMsgParser interface {
	C2S() interface{}
	S2C() interface{}
	C2SData() []byte
	S2CData() []byte
	C2SString() string
	S2CString() string
}

type MsgParser struct {
	s2c     interface{}
	c2s     interface{}
	c2sFunc ParseFunc
	s2cFunc ParseFunc
	parser  IParser
}

func (r *MsgParser) C2S() interface{} {
	if r.c2s == nil && r.c2sFunc != nil {
		r.c2s = r.c2sFunc()
	}
	return r.c2s
}

func (r *MsgParser) S2C() interface{} {
	if r.s2c == nil && r.s2cFunc != nil {
		r.s2c = r.s2cFunc()
	}
	return r.s2c
}

func (r *MsgParser) C2SData() []byte {
	return r.parser.PackMsg(r.C2S())
}

func (r *MsgParser) S2CData() []byte {
	return r.parser.PackMsg(r.S2C())
}

func (r *MsgParser) C2SString() string {
	return string(r.C2SData())
}

func (r *MsgParser) S2CString() string {
	return string(r.S2CData())
}

type ParserType int

const (
	ParserTypePB   ParserType = iota //protobuf类型，用于和客户端交互
	ParserTypeJson                   //json类型，可以用于客户端或者服务器之间交互
	ParserTypeCmd                    //cmd类型，类似telnet指令，用于直接和程序交互
	ParserTypeRaw                    //不做任何解析
)

type ParseErrType int

const (
	ParseErrTypeSendRemind ParseErrType = iota //消息解析失败发送提醒消息
	ParseErrTypeContinue                       //消息解析失败则跳过本条消息
	ParseErrTypeAlways                         //消息解析失败依然处理
	ParseErrTypeClose                          //消息解析失败则关闭连接
)

type ParseFunc func() interface{}

type IParser interface {
	GetType() ParserType
	GetErrType() ParseErrType
	ParseC2S(msg *Message) (IMsgParser, error)
	PackMsg(v interface{}) []byte
	GetRemindMsg(err error, t MsgType) *Message
}

type Parser struct {
	Type    ParserType
	ErrType ParseErrType

	msgMap  map[int]MsgParser
	cmdRoot *cmdParseNode
	jsonMap map[string]*jsonParseNode
	parser  IParser
}

func (r *Parser) Get() IParser {
	switch r.Type {
	case ParserTypePB:
		if r.parser == nil {
			r.parser = &pBParser{r}
		}
	case ParserTypeJson:
		if r.parser == nil {
			r.parser = &jsonParser{r}
		}
	case ParserTypeCmd:
		return &cmdParser{factory: r}
	case ParserTypeRaw:
		return nil
	}

	return r.parser
}

func (r *Parser) RegisterFunc(cmd, act uint8, c2sFunc ParseFunc, s2cFunc ParseFunc) {
	if r.msgMap == nil {
		r.msgMap = map[int]MsgParser{}
	}

	r.msgMap[CmdAct(cmd, act)] = MsgParser{c2sFunc: c2sFunc, s2cFunc: s2cFunc}
}

func (r *Parser) Register(cmd, act uint8, c2s interface{}, s2c interface{}) {
	if r.msgMap == nil {
		r.msgMap = map[int]MsgParser{}
	}

	p := MsgParser{}
	if c2s != nil {
		c2sType := reflect.TypeOf(c2s).Elem()
		p.c2sFunc = func() interface{} {
			return reflect.New(c2sType).Interface()
		}
	}
	if s2c != nil {
		s2cType := reflect.TypeOf(s2c).Elem()
		p.s2cFunc = func() interface{} {
			return reflect.New(s2cType).Interface()
		}
	}
	if c2s != nil || s2c != nil {
		r.msgMap[CmdAct(cmd, act)] = p
	}
}

func (r *Parser) RegisterMsgFunc(c2sFunc ParseFunc, s2cFunc ParseFunc) {
	if r.Type == ParserTypeCmd {
		if r.cmdRoot == nil {
			r.cmdRoot = &cmdParseNode{}
		}
		registerCmdParser(r.cmdRoot, c2sFunc, s2cFunc)
	}

	if r.Type == ParserTypeJson {
		if r.jsonMap == nil {
			r.jsonMap = map[string]*jsonParseNode{}
		}
		registerJsonParser(r.jsonMap, c2sFunc, s2cFunc)
	}
}

func (r *Parser) RegisterMsg(c2s interface{}, s2c interface{}) {
	var c2sFunc ParseFunc = nil
	var s2cFunc ParseFunc = nil
	if c2s != nil {
		c2sType := reflect.TypeOf(c2s).Elem()
		c2sFunc = func() interface{} {
			return reflect.New(c2sType).Interface()
		}
	}
	if s2c != nil {
		s2cType := reflect.TypeOf(s2c).Elem()
		s2cFunc = func() interface{} {
			return reflect.New(s2cType).Interface()
		}
	}

	if r.Type == ParserTypeCmd {
		if r.cmdRoot == nil {
			r.cmdRoot = &cmdParseNode{}
		}
		registerCmdParser(r.cmdRoot, c2sFunc, s2cFunc)
	}

	if r.Type == ParserTypeJson {
		if r.jsonMap == nil {
			r.jsonMap = map[string]*jsonParseNode{}
		}
		registerJsonParser(r.jsonMap, c2sFunc, s2cFunc)
	}
}
