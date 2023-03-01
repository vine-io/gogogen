//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Copyright 2023 lack
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Code generated by ___go_build_github_com_vine_io_gogogen_cmd_deepcopy_gen. Do NOT EDIT.

package meta

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *Meta) DeepCopyInto(out *Meta) {
	*out = *in
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an auto-generated deepcopy function, copying the receiver, creating a new Meta.
func (in *Meta) DeepCopy() *Meta {
	if in == nil {
		return nil
	}
	out := new(Meta)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *Resource) DeepCopyInto(out *Resource) {
	*out = *in
	in.Meta.DeepCopyInto(&out.Meta)
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Ann != nil {
		in, out := &in.Ann, &out.Ann
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Subs != nil {
		in, out := &in.Subs, &out.Subs
		*out = make([]*Sub, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*out)[i], &(*out)[i]
				*out = new(Sub)
				**out = **in
			}
		}
	}
	if in.SubMap != nil {
		in, out := &in.SubMap, &out.SubMap
		*out = make(map[string]*Sub, len(*in))
		for key, val := range *in {
			var outVal *Sub
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(Sub)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
	if in.Enable != nil {
		in, out := &in.Enable, &out.Enable
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an auto-generated deepcopy function, copying the receiver, creating a new Resource.
func (in *Resource) DeepCopy() *Resource {
	if in == nil {
		return nil
	}
	out := new(Resource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *Sub) DeepCopyInto(out *Sub) {
	*out = *in
	return
}

// DeepCopy is an auto-generated deepcopy function, copying the receiver, creating a new Sub.
func (in *Sub) DeepCopy() *Sub {
	if in == nil {
		return nil
	}
	out := new(Sub)
	in.DeepCopyInto(out)
	return out
}
