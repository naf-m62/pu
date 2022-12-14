// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package rmqhandlers

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson8b02993dDecodeFandoverseUserCmdHandlersRmqHandlers(in *jlexer.Lexer, out *AddPointsMessage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "points":
			out.Points = int32(in.Int32())
		case "userID":
			out.UserID = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson8b02993dEncodeFandoverseUserCmdHandlersRmqHandlers(out *jwriter.Writer, in AddPointsMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"points\":"
		out.RawString(prefix)
		out.Int32(int32(in.Points))
	}
	{
		const prefix string = ",\"userID\":"
		out.RawString(prefix)
		out.Int64(int64(in.UserID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AddPointsMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8b02993dEncodeFandoverseUserCmdHandlersRmqHandlers(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AddPointsMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson8b02993dEncodeFandoverseUserCmdHandlersRmqHandlers(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AddPointsMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson8b02993dDecodeFandoverseUserCmdHandlersRmqHandlers(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AddPointsMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson8b02993dDecodeFandoverseUserCmdHandlersRmqHandlers(l, v)
}
