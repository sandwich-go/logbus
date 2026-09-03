package thinkingdata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalAsJsonV2_AllFieldsNonEmpty(t *testing.T) {
	data := initData(10)
	data.Appid = "appid789"

	jsonV1, err := data.MarshalAsJson()
	require.NoError(t, err)

	jsonV2, err := data.MarshalAsJsonV2()
	require.NoError(t, err)

	require.JSONEq(t, string(jsonV1), string(jsonV2))
	require.Contains(t, string(jsonV2), `"#app_id":"appid789"`)
	require.Contains(t, string(jsonV2), `"properties"`)

	callbackCalls := 0
	err = data.WithJSONV2(func(jsonV2Borrowed []byte) {
		callbackCalls++
		require.JSONEq(t, string(jsonV2), string(jsonV2Borrowed))
	})
	require.NoError(t, err)
	require.Equal(t, 1, callbackCalls)
}

func BenchmarkMarshalAsJsonSmallData(b *testing.B) {
	data := initData(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := data.MarshalAsJson()
		require.NoError(b, err)
	}
}

func BenchmarkMarshalAsJsonMediumData(b *testing.B) {
	data := initData(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := data.MarshalAsJson()
		require.NoError(b, err)
	}
}

func BenchmarkMarshalAsJsonLargeData(b *testing.B) {
	data := initData(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := data.MarshalAsJson()
		require.NoError(b, err)
	}
}

func BenchmarkMarshalAsJsonV2SmallData(b *testing.B) {
	data := initData(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := data.MarshalAsJsonV2()
		require.NoError(b, err)
	}
}

func BenchmarkMarshalAsJsonV2MediumData(b *testing.B) {
	data := initData(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := data.MarshalAsJsonV2()
		require.NoError(b, err)
	}
}

func BenchmarkMarshalAsJsonV2LargeData(b *testing.B) {
	data := initData(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := data.MarshalAsJsonV2()
		require.NoError(b, err)
	}
}

func initData(numProperties int) Data {
	return Data{
		AccountId:    "account123",
		DistinctId:   "distinct456",
		Type:         "type1",
		Time:         "2023-10-20 10:00:00",
		EventName:    "event1",
		EventId:      "event123",
		FirstCheckId: "check123",
		Ip:           "127.0.0.1",
		UUID:         "uuid123",
		Properties:   initProperties(numProperties),
	}
}

func initProperties(numElements int) map[string]interface{} {
	properties := make(map[string]interface{})
	for i := 1; i <= numElements; i++ {
		key := fmt.Sprintf("property%d", i)
		value := fmt.Sprintf("value%d", i)
		properties[key] = value
	}
	return properties
}
