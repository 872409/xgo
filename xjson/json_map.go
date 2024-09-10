package xjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var NULL = ""

type Map map[string]interface{}

//
//type xmlMapMarshal struct {
//	XMLName xml.Name
//	Value   interface{} `xml:",cdata"`
//}
//
//type xmlMapUnmarshal struct {
//	XMLName xml.Name
//	Value   string `xml:",cdata"`
//}

var mu = new(sync.RWMutex)

// 设置参数
func (bm Map) Set(key string, value interface{}) Map {
	mu.Lock()
	bm[key] = value
	mu.Unlock()
	return bm
}

// 设置参数
func (bm Map) Keys() []string {
	mu.Lock()
	keys := make([]string, 0, len(bm))
	for k, _ := range bm {
		keys = append(keys, k)
	}
	return keys
}

func (bm Map) SetBodyMap(key string, value func(bm Map)) Map {
	_bm := make(Map)
	value(_bm)

	mu.Lock()
	bm[key] = _bm
	mu.Unlock()
	return bm
}

// Get 获取参数，同 GetString()
func (bm Map) Get(key string) string {
	return bm.GetString(key)
}

// GetString 获取参数转换string
func (bm Map) GetString(key string) string {
	if bm == nil {
		return NULL
	}
	mu.RLock()
	defer mu.RUnlock()
	value, ok := bm[key]
	if !ok {
		return NULL
	}
	v, ok := value.(string)
	if !ok {
		return convertToString(value)
	}
	return v
}

// GetInterface 获取原始参数
func (bm Map) GetInterface(key string) interface{} {
	if bm == nil {
		return nil
	}
	mu.RLock()
	defer mu.RUnlock()
	return bm[key]
}

func (bm Map) GetInt64(key string) int64 {
	val := bm.GetInterface(key)
	if val == nil {
		return 0
	}

	valInt, err := strconv.ParseInt(fmt.Sprintf("%s", val), 10, 64)
	if err != nil {
		return 0
	}
	return valInt
}

// 删除参数
func (bm Map) Remove(key string) {
	mu.Lock()
	delete(bm, key)
	mu.Unlock()
}

// 置空BodyMap
func (bm Map) Reset() {
	mu.Lock()
	for k := range bm {
		delete(bm, k)
	}
	mu.Unlock()
}

func (bm Map) Json() (jb string) {
	mu.Lock()
	defer mu.Unlock()
	bs, err := json.Marshal(bm)
	if err != nil {
		return ""
	}
	jb = string(bs)
	return jb
}

//
//func (bm Map) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
//	if len(bm) == 0 {
//		return nil
//	}
//	start.Name = xml.Name{NULL, "xml"}
//	if err = e.EncodeToken(start); err != nil {
//		return
//	}
//	for k := range bm {
//		if v := bm.GetString(k); v != NULL {
//			e.Encode(xmlMapMarshal{XMLName: xml.Name{Local: k}, Value: v})
//		}
//	}
//	return e.EncodeToken(start.End())
//}
//
//func (bm *Map) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
//	for {
//		var e xmlMapUnmarshal
//		err = d.Decode(&e)
//		if err != nil {
//			if err == io.EOF {
//				return nil
//			}
//			return err
//		}
//		bm.Set(e.XMLName.Local, e.Value)
//	}
//}

// ("bar=baz&foo=quux")
func (bm Map) EncodeURLParams() string {
	var (
		buf strings.Builder
	)
	for k, _ := range bm {
		if v := bm.GetString(k); v != NULL {
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v)
			buf.WriteByte('&')
		}
	}
	if buf.Len() <= 0 {
		return NULL
	}
	return buf.String()[:buf.Len()-1]
}

func (bm Map) CheckEmptyError(keys ...string) error {
	var emptyKeys []string
	for _, k := range keys {
		if v := bm.GetString(k); v == NULL {
			emptyKeys = append(emptyKeys, k)
		}
	}
	if len(emptyKeys) > 0 {
		return errors.New(strings.Join(emptyKeys, ", ") + " : cannot be empty")
	}
	return nil
}

func convertToString(v interface{}) (str string) {
	if v == nil {
		return NULL
	}
	var (
		bs  []byte
		err error
	)
	if bs, err = json.Marshal(v); err != nil {
		return NULL
	}
	str = string(bs)
	return
}
