//go:build !ignore_autogenerated

/*
Copyright 2024.

Licensed under the BSD 3-Clause License
you may see the license in the LICENSE.md file
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Challenge) DeepCopyInto(out *Challenge) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Challenge.
func (in *Challenge) DeepCopy() *Challenge {
	if in == nil {
		return nil
	}
	out := new(Challenge)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Challenge) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChallengeFile) DeepCopyInto(out *ChallengeFile) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChallengeFile.
func (in *ChallengeFile) DeepCopy() *ChallengeFile {
	if in == nil {
		return nil
	}
	out := new(ChallengeFile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChallengeHint) DeepCopyInto(out *ChallengeHint) {
	*out = *in
	if in.Requirements != nil {
		in, out := &in.Requirements, &out.Requirements
		*out = new(HintRequirements)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChallengeHint.
func (in *ChallengeHint) DeepCopy() *ChallengeHint {
	if in == nil {
		return nil
	}
	out := new(ChallengeHint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChallengeList) DeepCopyInto(out *ChallengeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Challenge, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChallengeList.
func (in *ChallengeList) DeepCopy() *ChallengeList {
	if in == nil {
		return nil
	}
	out := new(ChallengeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChallengeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChallengeRequirements) DeepCopyInto(out *ChallengeRequirements) {
	*out = *in
	if in.Anonymize != nil {
		in, out := &in.Anonymize, &out.Anonymize
		*out = new(bool)
		**out = **in
	}
	if in.Prerequisites != nil {
		in, out := &in.Prerequisites, &out.Prerequisites
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChallengeRequirements.
func (in *ChallengeRequirements) DeepCopy() *ChallengeRequirements {
	if in == nil {
		return nil
	}
	out := new(ChallengeRequirements)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChallengeSpec) DeepCopyInto(out *ChallengeSpec) {
	*out = *in
	if in.Initial != nil {
		in, out := &in.Initial, &out.Initial
		*out = new(int)
		**out = **in
	}
	if in.Decay != nil {
		in, out := &in.Decay, &out.Decay
		*out = new(int)
		**out = **in
	}
	if in.Minimum != nil {
		in, out := &in.Minimum, &out.Minimum
		*out = new(int)
		**out = **in
	}
	if in.MaxAttempts != nil {
		in, out := &in.MaxAttempts, &out.MaxAttempts
		*out = new(int)
		**out = **in
	}
	if in.Hints != nil {
		in, out := &in.Hints, &out.Hints
		*out = make([]ChallengeHint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Files != nil {
		in, out := &in.Files, &out.Files
		*out = make([]ChallengeFile, len(*in))
		copy(*out, *in)
	}
	in.Requirements.DeepCopyInto(&out.Requirements)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChallengeSpec.
func (in *ChallengeSpec) DeepCopy() *ChallengeSpec {
	if in == nil {
		return nil
	}
	out := new(ChallengeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HintRequirements) DeepCopyInto(out *HintRequirements) {
	*out = *in
	if in.Anonymize != nil {
		in, out := &in.Anonymize, &out.Anonymize
		*out = new(bool)
		**out = **in
	}
	if in.Prerequisites != nil {
		in, out := &in.Prerequisites, &out.Prerequisites
		*out = make([]int, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HintRequirements.
func (in *HintRequirements) DeepCopy() *HintRequirements {
	if in == nil {
		return nil
	}
	out := new(HintRequirements)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstancedChallenge) DeepCopyInto(out *InstancedChallenge) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstancedChallenge.
func (in *InstancedChallenge) DeepCopy() *InstancedChallenge {
	if in == nil {
		return nil
	}
	out := new(InstancedChallenge)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstancedChallenge) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstancedChallengeExposedPort) DeepCopyInto(out *InstancedChallengeExposedPort) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstancedChallengeExposedPort.
func (in *InstancedChallengeExposedPort) DeepCopy() *InstancedChallengeExposedPort {
	if in == nil {
		return nil
	}
	out := new(InstancedChallengeExposedPort)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstancedChallengeList) DeepCopyInto(out *InstancedChallengeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]InstancedChallenge, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstancedChallengeList.
func (in *InstancedChallengeList) DeepCopy() *InstancedChallengeList {
	if in == nil {
		return nil
	}
	out := new(InstancedChallengeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstancedChallengeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstancedChallengeOraclePodPort) DeepCopyInto(out *InstancedChallengeOraclePodPort) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstancedChallengeOraclePodPort.
func (in *InstancedChallengeOraclePodPort) DeepCopy() *InstancedChallengeOraclePodPort {
	if in == nil {
		return nil
	}
	out := new(InstancedChallengeOraclePodPort)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstancedChallengePod) DeepCopyInto(out *InstancedChallengePod) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]InstancedChallengePodPort, len(*in))
		copy(*out, *in)
	}
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstancedChallengePod.
func (in *InstancedChallengePod) DeepCopy() *InstancedChallengePod {
	if in == nil {
		return nil
	}
	out := new(InstancedChallengePod)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstancedChallengePodPort) DeepCopyInto(out *InstancedChallengePodPort) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstancedChallengePodPort.
func (in *InstancedChallengePodPort) DeepCopy() *InstancedChallengePodPort {
	if in == nil {
		return nil
	}
	out := new(InstancedChallengePodPort)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstancedChallengeRegistrySecret) DeepCopyInto(out *InstancedChallengeRegistrySecret) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstancedChallengeRegistrySecret.
func (in *InstancedChallengeRegistrySecret) DeepCopy() *InstancedChallengeRegistrySecret {
	if in == nil {
		return nil
	}
	out := new(InstancedChallengeRegistrySecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstancedChallengeSpec) DeepCopyInto(out *InstancedChallengeSpec) {
	*out = *in
	in.ChallengeSpec.DeepCopyInto(&out.ChallengeSpec)
	if in.ExposedPorts != nil {
		in, out := &in.ExposedPorts, &out.ExposedPorts
		*out = make([]InstancedChallengeExposedPort, len(*in))
		copy(*out, *in)
	}
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		*out = make([]InstancedChallengePod, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.RegistrySecret != nil {
		in, out := &in.RegistrySecret, &out.RegistrySecret
		*out = new(InstancedChallengeRegistrySecret)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstancedChallengeSpec.
func (in *InstancedChallengeSpec) DeepCopy() *InstancedChallengeSpec {
	if in == nil {
		return nil
	}
	out := new(InstancedChallengeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OracleInstancedChallenge) DeepCopyInto(out *OracleInstancedChallenge) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OracleInstancedChallenge.
func (in *OracleInstancedChallenge) DeepCopy() *OracleInstancedChallenge {
	if in == nil {
		return nil
	}
	out := new(OracleInstancedChallenge)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OracleInstancedChallenge) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OracleInstancedChallengeList) DeepCopyInto(out *OracleInstancedChallengeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]OracleInstancedChallenge, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OracleInstancedChallengeList.
func (in *OracleInstancedChallengeList) DeepCopy() *OracleInstancedChallengeList {
	if in == nil {
		return nil
	}
	out := new(OracleInstancedChallengeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OracleInstancedChallengeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OracleInstancedChallengeSpec) DeepCopyInto(out *OracleInstancedChallengeSpec) {
	*out = *in
	in.InstancedChallengeSpec.DeepCopyInto(&out.InstancedChallengeSpec)
	out.OraclePort = in.OraclePort
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OracleInstancedChallengeSpec.
func (in *OracleInstancedChallengeSpec) DeepCopy() *OracleInstancedChallengeSpec {
	if in == nil {
		return nil
	}
	out := new(OracleInstancedChallengeSpec)
	in.DeepCopyInto(out)
	return out
}
