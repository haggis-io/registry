package model

import (
	"database/sql/driver"
	"github.com/masterminds/semver"
)

type Version string

func (v Version) Validate() error {
	_, err := semver.NewVersion(string(v))

	if err != nil {
		return err
	}

	return nil
}

func (v Version) ToSemanticVersion() *semver.Version {
	ver, err := semver.NewVersion(string(v))

	if err != nil {
		return nil
	}

	return ver
}

func (v *Version) Value() (driver.Value, error) {
	return string(*v), nil
}

func (v *Version) Scan(src interface{}) error {
	*v = Version(src.(string))
	return nil
}
