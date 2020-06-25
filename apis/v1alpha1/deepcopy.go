package v1alpha1

func (in *BucketList) DeepCopy() *BucketList {
	if in == nil {
		return nil
	}
	out := new(BucketList)
	in.DeepCopyInto(out)
	return out
}

func (in *BucketList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *BucketList) DeepCopyInto(out *BucketList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)

	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Bucket, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *Bucket) DeepCopy() *Bucket {
	if in == nil {
		return nil
	}
	out := new(Bucket)
	in.DeepCopyInto(out)
	return out
}

func (in *Bucket) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *Bucket) DeepCopyInto(out *Bucket) {
	*out = *in
	out.TypeMeta = in.TypeMeta

	in.Meta.DeepCopyInto(&out.Meta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *BucketSpec) DeepCopyInto(out *BucketSpec) {
	*out = *in

	out.BucketClassName = in.BucketClassName
	out.Provisioner = in.Provisioner
	in.AnonymousAccessMode.DeepCopyInto(&out.AnonymousAccessMode)

	if in.PermittedNamespaces != nil {
		in, out := &in.PermittedNamespaces, &out.PermittedNamespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.Parameters = make(map[string]string, len(in.Parameters))
	for key, val := range in.Parameters {
		out.Parameters[key] = val
	}
}

func (in *BucketStatus) DeepCopyInto(out *BucketStatus) {
	*out = *in
	out.Message = in.Message
	out.Phase = in.Phase
	in.BoundBucketRequests.DeepCopyInto(&out)
}

func (in *BucketRequestBinding) DeepCopyInto(out *BucketRequestBinding) {
	*out = *in
	out = make(BucketRequestBinding)
	out.Name = in.Name
	out.Namespace = in.Namespace
}

func (in *BucketRequestList) DeepCopy() *BucketRequestList {
	if in == nil {
		return nil
	}
	out := new(BucketRequestList)
	in.DeepCopyInto(out)
	return out
}

func (in *BucketRequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *BucketRequestList) DeepCopyInto(out *BucketRequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)

	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BucketRequest, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *BucketClassList) DeepCopy() *BucketClassList {
	if in == nil {
		return nil
	}
	out := new(BucketClassList)
	in.DeepCopyInto(out)
	return out
}

func (in *BucketClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *BucketClassList) DeepCopyInto(out *BucketClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)

	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BucketClass, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}
