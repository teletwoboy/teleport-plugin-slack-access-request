package message

import "github.com/slack-go/slack"

/*
Builder 는 빌더 패턴을 따름
Message의 종류는 여러개이기 때문에 각각 Message를 보내기 위해
service.go 에서 모든 것을 구현하는 것은 매우 복잡한 일임

service 에선 하나의 Message 보내는 메서드( PostMessage() )만 만들고,
여러 가지 Message 종류를 받아서 Build 후 사용하기 위함
*/
type Builder interface {
	Build() slack.MsgOption
}
