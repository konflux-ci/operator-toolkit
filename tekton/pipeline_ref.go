package tekton

import tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"

type PipelineRef struct {
	Resolver string  `json:"resolver"`
	Params   []Param `json:"params"`
}

type Param struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (pr PipelineRef) ToTektonPipelineRef() tektonv1.PipelineRef {
	params := tektonv1.Params{}

	for _, p := range pr.Params {
		params = append(params, tektonv1.Param{
			Name: p.Name,
			Value: tektonv1.ParamValue{
				Type:      tektonv1.ParamTypeString,
				StringVal: p.Value,
			},
		})
	}

	tektonPipelineRef := tektonv1.PipelineRef{
		ResolverRef: tektonv1.ResolverRef{
			Resolver: tektonv1.ResolverName(pr.Resolver),
			Params:   params,
		},
	}

	return tektonPipelineRef
}

func (pr PipelineRef) IsClusterScoped() bool {
	return pr.Resolver == "cluster"
}
