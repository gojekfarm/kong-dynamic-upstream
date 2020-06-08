package dynamicupstream_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gojekfarm/kong-dyamic-upstream/dynamicupstream"
)

func TestInitializeWithPrefixWhenPrefixAndUpstream(t *testing.T) {
	du, _ := dynamicupstream.New("/foo/$upstream/bar", "1234")
	assert.Equal(t, "/foo/", du.Prefix)
}

func TestInitializeErrorsWhenPrefixAndNoUpstream(t *testing.T) {
	du, err := dynamicupstream.New("/foo/bar", "1234")
	assert.Error(t, err)
	assert.Equal(t, "Upstream doesnt have a prefix", err.Error())
	assert.Nil(t, du)

	du, err = dynamicupstream.New("/foo/$somethingelse/", "1234")
	assert.Error(t, err)
	assert.Equal(t, "Upstream doesnt have a prefix", err.Error())
	assert.Nil(t, du)
}

func TestInitializeErrorsWhenPortInvalid(t *testing.T) {
	du, err := dynamicupstream.New("/foo/$upstream", "foo")
	assert.Error(t, err)
	assert.Equal(t, "strconv.Atoi: parsing \"foo\": invalid syntax", err.Error())
	assert.Nil(t, du)

	du, err = dynamicupstream.New("/foo/$upstream/bar", "bar")
	assert.Error(t, err)
	assert.Equal(t, "strconv.Atoi: parsing \"bar\": invalid syntax", err.Error())

	assert.Nil(t, du)
}

func TestMatchPrefix(t *testing.T) {
	du, _ := dynamicupstream.New("/foo/$upstream/bar", "1234")

	assert.True(t, du.MatchPrefix("/foo/a.b.c/bar/baz"))
	assert.False(t, du.MatchPrefix("/bar/a.b.c/bar/baz"))
	assert.False(t, du.MatchPrefix("/foooooo/a.b.c/bar/baz"))
}

type MockKongPDK struct {
	mock.Mock
}

func (r *MockKongPDK) RequestPath() (string, error) {
	args := r.Called()
	return args.String(0), args.Error(1)
}

func (r *MockKongPDK) SetUpstreamTarget(host string, port int) error {
	args := r.Called(host, port)
	return args.Error(0)
}

func (r *MockKongPDK) SetUpstreamTargetRequestPath(path string) error {
	args := r.Called(path)
	return args.Error(0)
}

func TestTargetforValidRequestPath(t *testing.T) {
	du, _ := dynamicupstream.New("/foo/$upstream/bar", "1234")
	mockKongPDK := &MockKongPDK{}

	mockKongPDK.On("RequestPath").Return("/foo/1.1.1.1/bar/baz", nil)
	mockKongPDK.On("SetUpstreamTarget", "1.1.1.1", 1234).Return(nil)
	mockKongPDK.On("SetUpstreamTargetRequestPath", "/bar/baz").Return(nil)
	err := du.Target(mockKongPDK)
	assert.Nil(t, err)
	assert.NoError(t, err)
	mockKongPDK.AssertExpectations(t)
}

func TestTargetforRequestPathReturnsError(t *testing.T) {
	du, _ := dynamicupstream.New("/foo/$upstream/bar", "1234")
	mockKongPDK := &MockKongPDK{}

	mockKongPDK.On("RequestPath").Return("", fmt.Errorf("some error"))

	err := du.Target(mockKongPDK)
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, "some error", err.Error())

	mockKongPDK.AssertExpectations(t)
}

func TestTargetforRequestPathPrefixNotMatchingReturnsError(t *testing.T) {
	du, _ := dynamicupstream.New("/foo/$upstream/bar", "1234")
	mockKongPDK := &MockKongPDK{}
	mockKongPDK.On("RequestPath").Return("/bar/1.1.1.1/bar/baz", nil)

	err := du.Target(mockKongPDK)
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, "Path /bar/1.1.1.1/bar/baz, didnt match prefix: /foo/", err.Error())

	mockKongPDK.AssertExpectations(t)
}

func TestWhenTargetSetFailsReturnsError(t *testing.T) {
	du, _ := dynamicupstream.New("/foo/$upstream/bar", "1234")
	mockKongPDK := &MockKongPDK{}
	mockKongPDK.On("RequestPath").Return("/foo/1.1.1.1/bar/baz", nil)
	mockKongPDK.On("SetUpstreamTarget", "1.1.1.1", 1234).Return(fmt.Errorf("some target set error"))

	err := du.Target(mockKongPDK)
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, "some target set error", err.Error())

	mockKongPDK.AssertExpectations(t)
}

func TestWhenTargetPathSetFailsReturnsError(t *testing.T) {
	du, _ := dynamicupstream.New("/foo/$upstream/bar", "1234")
	mockKongPDK := &MockKongPDK{}
	mockKongPDK.On("RequestPath").Return("/foo/1.1.1.1/bar/baz", nil)
	mockKongPDK.On("SetUpstreamTarget", "1.1.1.1", 1234).Return(nil)
	mockKongPDK.On("SetUpstreamTargetRequestPath", "/bar/baz").Return(fmt.Errorf("some target request path set error"))

	err := du.Target(mockKongPDK)
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, "some target request path set error", err.Error())

	mockKongPDK.AssertExpectations(t)
}
