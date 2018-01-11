package main

import (
)

type lookupStrategy interface {
    increment (v Versions) (Versions)
}

type patchLookupStrategy struct {
}
func (s patchLookupStrategy) increment (v Versions) Versions {
    v.PATCH = v.PATCH + 1
    return v
}
type minorLookupStrategy struct {
}
func (s minorLookupStrategy) increment (v Versions) Versions {
    v.MINOR = v.MINOR + 1
    return v
}
type majorLookupStrategy struct {
}
func (s majorLookupStrategy) increment (v Versions) Versions {
    v.MAJOR = v.MAJOR + 1
    return v
}

type Versions struct {
    PATCH, MINOR, MAJOR int
}