// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"strings"

	"github.com/juju/cmd"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/juju/charm.v4"

	"github.com/juju/juju/cmd/envcmd"
	jujutesting "github.com/juju/juju/juju/testing"
	"github.com/juju/juju/testcharms"
	"github.com/juju/juju/testing"
)

type UnexposeSuite struct {
	jujutesting.RepoSuite
}

var _ = gc.Suite(&UnexposeSuite{})

func runUnexpose(c *gc.C, args ...string) error {
	_, err := testing.RunCommand(c, envcmd.Wrap(&UnexposeCommand{}), args...)
	return err
}

func (s *UnexposeSuite) assertExposed(c *gc.C, service string, expected bool) {
	svc, err := s.State.Service(service)
	c.Assert(err, jc.ErrorIsNil)
	actual := svc.IsExposed()
	c.Assert(actual, gc.Equals, expected)
}

func (s *UnexposeSuite) TestUnexpose(c *gc.C) {
	testcharms.Repo.CharmArchivePath(s.SeriesPath, "dummy")
	err := runDeploy(c, "local:dummy", "some-service-name")
	c.Assert(err, jc.ErrorIsNil)
	curl := charm.MustParseURL("local:trusty/dummy-1")
	s.AssertService(c, "some-service-name", curl, 1, 0)

	err = runExpose(c, "some-service-name")
	c.Assert(err, jc.ErrorIsNil)
	s.assertExposed(c, "some-service-name", true)

	err = runUnexpose(c, "some-service-name")
	c.Assert(err, jc.ErrorIsNil)
	s.assertExposed(c, "some-service-name", false)

	err = runUnexpose(c, "nonexistent-service")
	c.Assert(err, gc.ErrorMatches, `service "nonexistent-service" not found`)
}

func (s *UnexposeSuite) TestBlockUnexpose(c *gc.C) {
	testcharms.Repo.CharmArchivePath(s.SeriesPath, "dummy")
	err := runDeploy(c, "local:dummy", "some-service-name")
	c.Assert(err, jc.ErrorIsNil)
	curl := charm.MustParseURL("local:trusty/dummy-1")
	s.AssertService(c, "some-service-name", curl, 1, 0)

	// Block operation
	s.AssertConfigParameterUpdated(c, "block-all-changes", true)
	err = runExpose(c, "some-service-name")
	c.Assert(err, gc.ErrorMatches, cmd.ErrSilent.Error())
	// msg is logged
	stripped := strings.Replace(c.GetTestLog(), "\n", "", -1)
	c.Check(stripped, gc.Matches, ".*To unblock changes.*")
}
