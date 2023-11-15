package cappedbuffer

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CappedBufferTestSuite struct {
	suite.Suite
}

func (s *CappedBufferTestSuite) TestCappedBufferOverflow() {
	var b []byte
	buff := New(b, 2)
	_, err := buff.Write([]byte("hello"))
	s.Assert().EqualError(err, "buffer overflow")
}

func (s *CappedBufferTestSuite) TestCappedBufferSuccessWithExactSize() {
	var b []byte
	buff := New(b, 2)
	bytesWritten, err := buff.Write([]byte("he"))
	s.Assert().Equal(err, nil)
	s.Assert().Equal(bytesWritten, 2)
}

func (s *CappedBufferTestSuite) TestCappedBufferSuccessWithSmallerSize() {
	var b []byte
	buff := New(b, 2)
	bytesWritten, err := buff.Write([]byte("h"))
	s.Assert().Equal(err, nil)
	s.Assert().Equal(bytesWritten, 1)
}

func TestCappedBuffer(t *testing.T) {
	suite.Run(t, new(CappedBufferTestSuite))
}
