package handlers

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/chamilto/dummy/internal/config"
	"github.com/chamilto/dummy/internal/db"
)

var loadenv = flag.Bool("loadenv", false, "load the .env.test file in the project root")

func loadEnv() {
	if !*loadenv {
		logrus.Info("no -loadenv option found. reading env vars...")
		return
	}

	if err := godotenv.Load("../../.env.test"); err != nil {
		logrus.Fatal("Error loading .env file")
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	loadEnv()
	v := m.Run()
	os.Exit(v)
}

func newTestingHandlerController() *HandlerController {
	c := config.NewConfig()
	db := db.NewDB(c)
	ctlr := NewHandlerController(c, db)
	_ = NewRouter(ctlr)

	return ctlr
}

func assertEquals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func getResponseBytes(t *testing.T, rr *httptest.ResponseRecorder) []byte {
	b, err := ioutil.ReadAll(rr.Body)

	if err != nil {
		t.Fatal(err)
	}

	return b
}

func assertStatus(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	if status := rr.Code; status != expected {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expected)
	}
}

func newRequest(t *testing.T, method, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, path, body)

	if err != nil {
		t.Fatal(err)
	}

	return req
}

func getTestData(t *testing.T, relpath string) []byte {
	bytes, err := ioutil.ReadFile(relpath)

	if err != nil {
		t.Fatal(err)
	}

	return bytes
}

type testPayload struct {
	name string
	data []byte
}

func getAllPayloadsForDir(t *testing.T, dir string) *[]testPayload {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		t.Fatal(err)
	}

	var tp []testPayload

	for _, f := range files {
		fn := f.Name()
		fmt.Println(fn)
		data := getTestData(t, filepath.Join(dir, fn))
		tp = append(tp, testPayload{data: data, name: fn})
	}

	return &tp
}

type mockRedisCallData struct {
	hsetCallCount    int
	hexistsCallCount int
}

// REDIS_CALL_TRACKER should be assigned to a new instance of
// mockRedisCallData before each test run that uses the MockRedisClient.
var REDIS_CALL_TRACKER = mockRedisCallData{}

type MockRedisClient struct {
	db.RedisClient
	failOnCallCount int
}

func (MockRedisClient) BuildKey(key string) string {
	return fmt.Sprintf("test:dummy:%s", key)
}

func (c MockRedisClient) HSet(key string, values ...interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd()

	if REDIS_CALL_TRACKER.hsetCallCount == c.failOnCallCount {
		cmd.SetErr(fmt.Errorf("Failure on HSet call count %d", REDIS_CALL_TRACKER.hsetCallCount))
	}

	REDIS_CALL_TRACKER.hsetCallCount++

	return cmd
}

func (c MockRedisClient) HExists(key, field string) *redis.BoolCmd {
	cmd := redis.NewBoolCmd("hexists", key, field)
	REDIS_CALL_TRACKER.hexistsCallCount++
	return cmd
}
