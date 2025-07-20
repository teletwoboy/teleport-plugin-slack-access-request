package accessrequest

import "github.com/gravitational/teleport/api/types"

type FilterBuilder interface {
	Build() types.AccessRequestFilter
}

type filterBuilder struct {
	accessRequestName string
}

func NewFilterBuilder(a string) FilterBuilder {
	return &filterBuilder{
		accessRequestName: a,
	}
}

func (f *filterBuilder) Build() types.AccessRequestFilter {
	return types.AccessRequestFilter{
		ID: f.accessRequestName,
	}
}
