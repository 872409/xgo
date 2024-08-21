package xerr

import (
	"testing"
)

func TestWrapCodeErrMsgF(t *testing.T) {

	//err := errors.Join(MerchantStatusInvalid.WithStack(), fmt.Errorf("aaaa"))
	//fmt.Printf("%+v\n", AsCodeError(err))
	//fmt.Println(MerchantStatusInvalid.WithStack())

	//err := WrapCodeErrMsgF(ChannelAPIError, "OpenC@Reply[%s=%s]", "reply.Code", "reply.Msg")
	//
	//AsCodeErrorCase(err, func(e *CodeError) {
	//	log.Printf("e %s:%s", e.Code(), e.Msg())
	//
	//	AsCodeErrorCase(e, func(ee *CodeError) {
	//		log.Printf("ee %s:%s", ee.Code(), ee.Msg())
	//	})
	//
	//	//if ee := xerr.UnwrapCodeError(e); ee != nil {
	//	//	log.Printf("ee %s:%s", ee.Code(), ee.Msg())
	//	//}
	//})

}
