package models

import (
	"reflect"
	"testing"
)

func getJSONTag(t *testing.T, typ reflect.Type, field string) string {
	t.Helper()
	f, ok := typ.FieldByName(field)
	if !ok {
		t.Fatalf("field %s not found on type %s", field, typ.Name())
	}
	return f.Tag.Get("json")
}

func TestJSONTags(t *testing.T) {
	tests := []struct {
		typ   reflect.Type
		field string
		want  string
	}{
		{reflect.TypeOf(User{}), "Password", "-"},
		{reflect.TypeOf(User{}), "CreatedAt", "created_at"},
		{reflect.TypeOf(Post{}), "UpdatedAt", "updated_at"},
		{reflect.TypeOf(Comment{}), "UpdatedAt", "updated_at"},
		{reflect.TypeOf(Post{}), "Categories", "categories,omitempty"},
		{reflect.TypeOf(PostReaction{}), "PostID", "post_id,omitempty"},
		{reflect.TypeOf(CommentReaction{}), "CommentID", "comment_id,omitempty"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.typ.Name()+"."+tc.field, func(t *testing.T) {
			t.Parallel()
			got := getJSONTag(t, tc.typ, tc.field)
			if got != tc.want {
				t.Fatalf("unexpected json tag for %s.%s: got %q want %q", tc.typ.Name(), tc.field, got, tc.want)
			}
		})
	}
}
