// Code generated by "stringer -type=ExchangeName -linecomment"; DO NOT EDIT.

package common

import "strconv"

const _ExchangeName_name = "binancehuobistable_exchange"

var _ExchangeName_index = [...]uint8{0, 7, 12, 27}

func (i ExchangeName) String() string {
	if i < 0 || i >= ExchangeName(len(_ExchangeName_index)-1) {
		return "ExchangeName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ExchangeName_name[_ExchangeName_index[i]:_ExchangeName_index[i+1]]
}