package gofuzz

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"testing"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

type GoFuzzSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&GoFuzzSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *GoFuzzSuite) TestParsing(c *C) {
	raw := "2021/09/21 21:46:20 workers: 8, corpus: 205 (1m45s ago), crashers: 2, restarts: 1/9710, execs: 4078324 (38839/sec), cover: 225, uptime: 1m45s"
	line, err := Parse(raw)

	c.Assert(err, IsNil)
	c.Assert(line.DateTime.String(), Equals, "2021-09-21 21:46:20 +0000 UTC")
	c.Assert(line.Workers, Equals, 8)
	c.Assert(line.Corpus, Equals, 205)
	c.Assert(line.Crashers, Equals, 2)
	c.Assert(line.Restarts, Equals, 9710)
	c.Assert(line.Execs, Equals, 4078324)
	c.Assert(line.ExecsPerSec, Equals, 38839)
	c.Assert(line.Cover, Equals, 225)
}

func (s *GoFuzzSuite) TestError(c *C) {
	// raw := "2021/09/21 21:44:38 workers: 8, corpus: 205 (3s ago), crashers: 0, restarts: 1/0, execs: 0 (0/sec), cover: 0, uptime: 3s"
	raw := "2021/09/21 21:44:38"
	_, err := Parse(raw)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Output line must contain at least 7 fields")

	raw = "2021/09/21 AA:AA:AA workers: 8, corpus: 205 (3s ago), crashers: 0, restarts: 1/0, execs: 0 (0/sec), cover: 0, uptime: 3s"
	_, err = Parse(raw)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Can't parse date and time field: parsing time \"2021/09/21 AA:AA:AA\" as \"2006/01/02 15:04:05\": cannot parse \"AA:AA:AA\" as \"15\"")

	raw = "2021/09/21 21:44:38 workers: A, corpus: 205 (3s ago), crashers: 0, restarts: 1/0, execs: 0 (0/sec), cover: 0, uptime: 3s"
	_, err = Parse(raw)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Can't parse workers field: strconv.Atoi: parsing \"A\": invalid syntax")

	raw = "2021/09/21 21:44:38 workers: 8, corpus: BBB (3s ago), crashers: 0, restarts: 1/0, execs: 0 (0/sec), cover: 0, uptime: 3s"
	_, err = Parse(raw)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Can't parse corpus field: strconv.Atoi: parsing \"BBB\": invalid syntax")

	raw = "2021/09/21 21:44:38 workers: 8, corpus: 205 (3s ago), crashers: V, restarts: 1/0, execs: 0 (0/sec), cover: 0, uptime: 3s"
	_, err = Parse(raw)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Can't parse crashers field: strconv.Atoi: parsing \"V\": invalid syntax")

	raw = "2021/09/21 21:44:38 workers: 8, corpus: 205 (3s ago), crashers: XX, restarts: 1/0, execs: 0 (0/sec), cover: 0, uptime: 3s"
	_, err = Parse(raw)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Can't parse crashers field: strconv.Atoi: parsing \"XX\": invalid syntax")

	raw = "2021/09/21 21:44:38 workers: 8, corpus: 205 (3s ago), crashers: 0, restarts: 1/G, execs: 0 (0/sec), cover: 0, uptime: 3s"
	_, err = Parse(raw)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Can't parse restarts field: strconv.Atoi: parsing \"G\": invalid syntax")

	raw = "2021/09/21 21:44:38 workers: 8, corpus: 205 (3s ago), crashers: 0, restarts: 1/0, execs: K (0/sec), cover: 0, uptime: 3s"
	_, err = Parse(raw)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Can't parse execs field: strconv.Atoi: parsing \"K\": invalid syntax")

	raw = "2021/09/21 21:44:38 workers: 8, corpus: 205 (3s ago), crashers: 0, restarts: 1/0, execs: 0 (H/sec), cover: 0, uptime: 3s"
	_, err = Parse(raw)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Can't parse execs per sec field: strconv.Atoi: parsing \"H\": invalid syntax")

	raw = "2021/09/21 21:44:38 workers: 8, corpus: 205 (3s ago), crashers: 0, restarts: 1/0, execs: 0 (0/sec), cover: T, uptime: 3s"
	_, err = Parse(raw)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Can't parse cover field: strconv.Atoi: parsing \"T\": invalid syntax")
}
