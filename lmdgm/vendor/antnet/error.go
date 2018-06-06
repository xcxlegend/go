package antnet

type Error struct {
	Id  uint16
	Str string
}

func (r *Error) Error() string {
	return r.Str
}

var idErrMap = map[uint16]error{}
var errIdMap = map[error]uint16{}

func NewError(str string, id uint16) *Error {
	err := &Error{id, str}
	idErrMap[id] = err
	errIdMap[err] = id
	return err
}

var (
	ErrOk            = NewError("正确", 0)
	ErrDBErr         = NewError("数据库错误", 1)
	ErrProtoPack     = NewError("协议解析错误", 2)
	ErrProtoUnPack   = NewError("协议打包错误", 3)
	ErrMsgPackPack   = NewError("msgpack打包错误", 4)
	ErrMsgPackUnPack = NewError("msgpack解析错误", 5)
	ErrPBPack        = NewError("pb打包错误", 6)
	ErrPBUnPack      = NewError("pb解析错误", 7)
	ErrJsonPack      = NewError("json打包错误", 8)
	ErrJsonUnPack    = NewError("json解析错误", 9)
	ErrCmdUnPack     = NewError("cmd解析错误", 10)
	ErrFileRead      = NewError("文件读取错误", 100)
	ErrDBDataType    = NewError("数据库数据类型错误", 101)
	ErrNetTimeout    = NewError("网络超时", 200)

	ErrErrIdNotFound = NewError("错误没有对应的错误码", 255)
)

var MinUserError = 256

func GetError(id uint16) error {
	if e, ok := idErrMap[id]; ok {
		return e
	}
	return ErrErrIdNotFound
}

func GetErrId(err error) uint16 {
	if id, ok := errIdMap[err]; ok {
		return id
	}
	return errIdMap[ErrErrIdNotFound]
}

type ErrJsonStr struct {
	Error    int    `json:"error"`
	ErrorStr string `json:"errstr"`
}

func GetErrJsonStr(err error) string {
	return string(GetErrJsonData(err))
}
func GetErrJsonData(err error) []byte {
	data, _ := JsonPack(&ErrJsonStr{Error: int(GetErrId(err)), ErrorStr: err.Error()})
	return data
}
