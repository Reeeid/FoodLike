package usecase

import "errors"

// ErrNotMember は操作対象のグループに呼び出し主が所属していない場合のエラー。
// ハンドラでは404に変換する(グループの存在を漏らさないため)。
var ErrNotMember = errors.New("member does not belong to the group")
