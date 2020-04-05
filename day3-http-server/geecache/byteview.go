package geecache

//A ByteView holds an immutable view of bytes
type ByteView struct {
	b []byte
}

//Len returns the view's length
func (v ByteView) Len() int{
	return  len(v.b)
}

//ByteSlice returns a copy of the data as a byte slice
func (v ByteView) ByteSlice() []byte{
	return cloneByte(v.b)
}

//String returns the data as a string
func (v ByteView) String() string{
	return string(v.b)
}

func cloneByte(b []byte) []byte{
	c := make([]byte,len(b))
	copy(c,b)
	return c
}