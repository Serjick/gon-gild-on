package gildedsergigodiff_test

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
