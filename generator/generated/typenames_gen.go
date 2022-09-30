// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

func SetTypenameRecursively[T any](x T) T {
	switch val := any(x).(type) {
	case nil:
		return x

	case AudioClip:
		val.Typename = "AudioClip"
		return any(val).(T)
	case *AudioClip:
		val.Typename = "AudioClip"
		return any(val).(T)

	case Library:
		if val.Books != nil {
			for i, e := range val.Books {
				e = SetTypenameRecursively(e)
				val.Books[i] = e
			}
		}
		return any(val).(T)
	case *Library:
		if val.Books != nil {
			for i, e := range val.Books {
				e = SetTypenameRecursively(e)
				val.Books[i] = e
			}
		}
		return any(val).(T)

	case TextBook:
		if val.SupplementaryMaterial != nil {
			for i, e := range val.SupplementaryMaterial {
				e = SetTypenameRecursively(e)
				val.SupplementaryMaterial[i] = e
			}
		}
		val.Typename = "TextBook"
		return any(val).(T)
	case *TextBook:
		if val.SupplementaryMaterial != nil {
			for i, e := range val.SupplementaryMaterial {
				e = SetTypenameRecursively(e)
				val.SupplementaryMaterial[i] = e
			}
		}
		val.Typename = "TextBook"
		return any(val).(T)

	case VideoClip:
		val.Typename = "VideoClip"
		return any(val).(T)
	case *VideoClip:
		val.Typename = "VideoClip"
		return any(val).(T)

	case []Book:
		for i, e := range val {
			e = SetTypenameRecursively(e)
			val[i] = e
		}
		return any(val).(T)

	case []MediaItem:
		for i, e := range val {
			e = SetTypenameRecursively(e)
			val[i] = e
		}
		return any(val).(T)

	default:
		return any(val).(T)
	}
}
