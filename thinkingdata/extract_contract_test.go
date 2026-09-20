package thinkingdata

import (
	"reflect"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Fixed expectations keep these tests independent of the formatting code shared
// by the scalar and MapObjectEncoder paths.
func TestExtractContract(t *testing.T) {
	stamp := time.Date(2026, 9, 20, 12, 34, 56, 789000000, time.FixedZone("UTC+8", 8*60*60))
	const formatted = "2026-09-20 12:34:56.789"
	base := []zap.Field{
		zap.String(ACCOUNT, "player"),
		zap.String(EVENT, "battle_finish"),
		zap.Time(TIME, stamp),
	}
	data := func(properties map[string]interface{}) Data {
		return Data{AccountId: "player", Type: TRACK, EventName: "battle_finish", Time: formatted, Properties: properties}
	}

	tests := []struct {
		name   string
		fields []zap.Field
		want   Data
		err    string
	}{
		{
			name: "invalid keys are silently filtered",
			fields: []zap.Field{
				zap.String("invalid-key", "ignored"), zap.String("1invalid", "ignored"),
				zap.String("", "ignored"), zap.Int("score", 7),
			},
			want: data(map[string]interface{}{"score": int64(7)}),
		},
		{
			name: "invalid key and reserved field collision",
			fields: []zap.Field{
				zap.String("#account-id", "ignored"),
				zap.Int(ACCOUNT, 123), zap.String(ACCOUNT, "last_player"),
				zap.String(EVENT, "last_event"), zap.String(TYPE, TRACK_UPDATE),
				zap.String(APPID, "app"), zap.Int("score", 7),
			},
			want: Data{AccountId: "last_player", Type: TRACK_UPDATE, EventName: "last_event", Appid: "app", Time: formatted, Properties: map[string]interface{}{"score": int64(7)}},
		},
		{
			name: "duplicate properties use last value",
			fields: []zap.Field{
				zap.Int("score", 1), zap.String("score", "last"),
				zap.String(EVENT_ID, "first"), zap.String(EVENT_ID, "9-last"),
			},
			want: Data{AccountId: "player", Type: TRACK, EventName: "battle_finish", EventId: "9-last", Time: formatted, Properties: map[string]interface{}{"score": "last"}},
		},
		{
			name:   "event ID wrong type",
			fields: []zap.Field{zap.Int(EVENT_ID, 1)},
			err:    "#event_id must be string",
		},
		{
			name:   "event ID invalid first character",
			fields: []zap.Field{zap.String(EVENT_ID, "-event")},
			err:    "the event name must start with a letter or number",
		},
		{
			name:   "nil event ID is empty",
			fields: []zap.Field{zap.Reflect(EVENT_ID, nil)},
			want:   data(map[string]interface{}{}),
		},
		{
			name:   "invalid event name remains rejected",
			fields: []zap.Field{zap.String(EVENT, "invalid-event")},
			err:    "Invalid event name: invalid-event",
		},
		{
			name:   "empty event name remains rejected",
			fields: []zap.Field{zap.String(EVENT, "")},
			err:    "the event name must be provided",
		},
		{
			name: "protocol time IP UUID and property time",
			fields: []zap.Field{
				zap.String(DISTINCT, "device"), zap.String(IP, "127.0.0.1"),
				zap.String(UUID, "51c7b43d-f8ca-447f-ae05-dfc5fcc184c9"),
				zap.Time("created_at", stamp),
			},
			want: Data{AccountId: "player", DistinctId: "device", Type: TRACK, EventName: "battle_finish", Time: formatted, Ip: "127.0.0.1", UUID: "51c7b43d-f8ca-447f-ae05-dfc5fcc184c9", Properties: map[string]interface{}{"created_at": formatted}},
		},
		{
			name:   "string time passes through",
			fields: []zap.Field{zap.String(TIME, "custom-time")},
			want:   Data{AccountId: "player", Type: TRACK, EventName: "battle_finish", Time: "custom-time", Properties: map[string]interface{}{}},
		},
		{
			name:   "non string UUID is removed",
			fields: []zap.Field{zap.Int(UUID, 1)},
			want:   data(map[string]interface{}{}),
		},
		{
			name: "object array and stringer fallback",
			fields: []zap.Field{
				zap.Object("object", testObject{}), zap.Array("array", testArray{}),
				zap.Stringer("text", testStringer("value")), zap.Time("created_at", stamp),
			},
			want: data(map[string]interface{}{"object": map[string]interface{}{"name": "object"}, "array": []interface{}{"array"}, "text": "value", "created_at": formatted}),
		},
		{
			name: "inline generated keys are filtered after encoding",
			fields: []zap.Field{zap.Inline(zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
				enc.AddString(ACCOUNT, "inline_player")
				enc.AddString("invalid-key", "ignored")
				enc.AddInt64("score", 3)
				return nil
			}))},
			want: Data{AccountId: "inline_player", Type: TRACK, EventName: "battle_finish", Time: formatted, Properties: map[string]interface{}{"score": int64(3)}},
		},
		{
			name:   "nested keys retain original behavior",
			fields: []zap.Field{zap.Namespace("nested"), zap.String("invalid-key", "retained")},
			want:   data(map[string]interface{}{"nested": map[string]interface{}{"invalid-key": "retained"}}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := append(append([]zap.Field{}, base...), tt.fields...)
			checkExtraction(t, fields, tt.want, tt.err)
		})
	}
}

func TestExtractUserContract(t *testing.T) {
	for _, tt := range []struct {
		name   string
		kind   string
		fields []zap.Field
		want   map[string]interface{}
		err    string
	}{
		{name: "add numbers", kind: USER_ADD, fields: []zap.Field{zap.Int("score", 2)}, want: map[string]interface{}{"score": int64(2)}},
		{name: "add string rejected", kind: USER_ADD, fields: []zap.Field{zap.String("score", "2")}, err: "Invalid property value: only numbers is supported by UserAdd"},
		{name: "add time rejected before formatting", kind: USER_ADD, fields: []zap.Field{zap.Time("score", time.Unix(0, 0))}, err: "Invalid property value: only numbers is supported by UserAdd"},
		{name: "add object rejected in fallback", kind: USER_ADD, fields: []zap.Field{zap.Object("score", testObject{})}, err: "Invalid property value: only numbers is supported by UserAdd"},
		{name: "invalid key filtered before add validation", kind: USER_ADD, fields: []zap.Field{zap.String("invalid-key", "not a number"), zap.Int("score", 2)}, want: map[string]interface{}{"score": int64(2)}},
		{name: "set time formatted", kind: USER_SET, fields: []zap.Field{zap.Time("created_at", time.Date(2026, 9, 20, 0, 0, 0, 0, time.UTC))}, want: map[string]interface{}{"created_at": "2026-09-20 00:00:00.000"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fields := append([]zap.Field{zap.String(ACCOUNT, "player"), zap.String(TYPE, tt.kind), zap.String(TIME, "fixed")}, tt.fields...)
			want := Data{AccountId: "player", Type: tt.kind, Time: "fixed", Properties: tt.want}
			checkExtraction(t, fields, want, tt.err)
		})
	}
}

func checkExtraction(t *testing.T, fields []zap.Field, want Data, wantErr string) {
	t.Helper()
	for _, path := range []struct {
		name string
		run  func([]zap.Field) (Data, error)
	}{
		{"fields", ExtractFields},
		{"encoder", extractWithMapObjectEncoder},
	} {
		t.Run(path.name, func(t *testing.T) {
			got, err := path.run(fields)
			if wantErr != "" {
				if err == nil || err.Error() != wantErr {
					t.Fatalf("error = %v, want %q", err, wantErr)
				}
				if !reflect.DeepEqual(got, emptyData) {
					t.Fatalf("failed extraction returned %#v", got)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("data = %#v, want %#v", got, want)
			}
		})
	}
}

func TestPublicAPIsValidatePropertyKeys(t *testing.T) {
	for _, tt := range []struct {
		name string
		run  func(map[string]interface{}) (Data, error)
	}{
		{"track", func(p map[string]interface{}) (Data, error) { return Track("player", "", "event", "", "", p) }},
		{"track with type", func(p map[string]interface{}) (Data, error) {
			return TrackWithType(TRACK_UPDATE, "player", "", "event", "", "", p)
		}},
		{"user", func(p map[string]interface{}) (Data, error) { return User("player", "", USER_SET, "", p) }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.run(map[string]interface{}{"invalid-key": 1})
			if err == nil || err.Error() != "Invalid property key: invalid-key" {
				t.Fatalf("error = %v, want invalid property key", err)
			}
		})
	}
}
