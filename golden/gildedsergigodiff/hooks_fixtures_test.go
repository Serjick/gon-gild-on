package gildedsergigodiff_test

const FieldLodgePrev = `{
    "key1": {
        "a": "{{ .Actual.item1.key }}",
        "b": {},
        "z": true
    },
    "key2": {
        "a": "{{ .Actual.item2.key }}",
        "b": {},
        "z": true
    },
    "key3": {
        "a": "{{ .Actual.item3.key }}",
        "b": {},
        "z": true
    }
}`

const FieldLodgeNext = `{
    "key1": {
        "a": "key1",
        "b": {},
        "c: 1,
        "z": true
    },
    "key2": {
        "a": "key2",
        "b": {},
        "c: 2,
        "z": true
    },
    "key3": {
        "a": "key3",
        "b": {},
        "c: 3,
        "z": true
    }
}`

const FieldLodgeWant = `{
    "key1": {
        "a": "{{ .Actual.item1.key }}",
        "b": {},
        "c: 1,
        "z": true
    },
    "key2": {
        "a": "{{ .Actual.item2.key }}",
        "b": {},
        "c: 2,
        "z": true
    },
    "key3": {
        "a": "{{ .Actual.item3.key }}",
        "b": {},
        "c: 3,
        "z": true
    }
}`

const FuncMapKeyPrev = `{
    "4a66538b-7507-4374-9c83-571904fdce98": {
        "url": "{{ urlquery (index .Actual "4a66538b-7507-4374-9c83-571904fdce98" "url") }}",
        "revision": 1
    }
`

const FuncMapKeyNext = `{
    "4a66538b-7507-4374-9c83-571904fdce98": {
        "url": "http://example.com?true&foo=bar",
        "revision": 2
    }
`

const FuncMapKeyWant = `{
    "4a66538b-7507-4374-9c83-571904fdce98": {
        "url": "{{ urlquery (index .Actual "4a66538b-7507-4374-9c83-571904fdce98" "url") }}",
        "revision": 2
    }
`

const TimeInMapPrev = `{
    "4a66538b-7507-4374-9c83-571904fdce98": {
        "time": "{{ index .Actual "4a66538b-7507-4374-9c83-571904fdce98" "time" }}",
        "revision": 1
    },
	"a31e9427-7273-4777-88a5-1e9714d67b2a": {
        "time": "{{ index .Actual "a31e9427-7273-4777-88a5-1e9714d67b2a" "time" }}",
        "revision": 1
    }
`

const TimeInMapNext = `{
    "4a66538b-7507-4374-9c83-571904fdce98": {
        "time": "2025-03-31T08:27:37.052418693Z",
        "revision": 2
    },
	"a31e9427-7273-4777-88a5-1e9714d67b2a": {
        "time": "2025-03-31T10:47:37.711121877Z",
        "revision": 2
    }
`

const TimeInMapWant = `{
    "4a66538b-7507-4374-9c83-571904fdce98": {
        "time": "{{ index .Actual "4a66538b-7507-4374-9c83-571904fdce98" "time" }}",
        "revision": 2
    },
	"a31e9427-7273-4777-88a5-1e9714d67b2a": {
        "time": "{{ index .Actual "a31e9427-7273-4777-88a5-1e9714d67b2a" "time" }}",
        "revision": 2
    }
`

const TimeFuzzPrev = `{
    "4a66538b-7507-4374-9c83-571904fdce98": {
        "time": "{{ index .Actual "4a66538b-7507-4374-9c83-571904fdce98" "time" }}",
        "revision": 1
    },
`

const TimeFuzzNext = `{
    "4a66538b-7507-4374-9c83-571904fdce98": {
        "time": "%s",
        "revision": 2
    },
`

const TimeFuzzWant = `{
    "4a66538b-7507-4374-9c83-571904fdce98": {
        "time": "{{ index .Actual "4a66538b-7507-4374-9c83-571904fdce98" "time" }}",
        "revision": 2
    },
`
