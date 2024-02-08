package entity

import (
	"reflect"
	"testing"
	"time"
)

func TestUser_AddSignedTest(t *testing.T) {
	type fields struct {
		ID       int
		Username string
		Tests    []*Test
	}
	type args struct {
		questions string
		answers   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test AddSignedTest",
			fields: fields{
				ID:       1,
				Username: "test",
				Tests:    nil,
			},
			args: args{
				questions: "test",
				answers:   "test",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := &User{
				ID:       tt.fields.ID,
				Username: tt.fields.Username,
				Tests:    tt.fields.Tests,
			}
			u.AddSignedTest(tt.args.questions, tt.args.answers)
			if len(u.Tests) != 1 {
				t.Errorf("User.AddSignedTest(%v, %v) = %v, want %v", tt.args.questions, tt.args.answers, len(u.Tests), 1)
			}
		})
	}
}

func TestUser_VerifySignature(t *testing.T) {
	type fields struct {
		ID       int
		Username string
		Tests    []*Test
	}
	type args struct {
		signature string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
		want2  time.Time
	}{
		{
			name: "Test VerifySignature",
			fields: fields{
				ID:       1,
				Username: "test",
				Tests: []*Test{
					{
						Questions: "test",
						Answers:   "test",
						Signature: "test",
						SignedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			args: args{
				signature: "test",
			},
			want:  true,
			want1: "test",
			want2: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Test VerifySignature not found",
			fields: fields{
				ID:       1,
				Username: "test",
				Tests: []*Test{
					{
						Questions: "test",
						Answers:   "test",
						Signature: "test",
						SignedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			args: args{
				signature: "notfound",
			},
			want:  false,
			want1: "",
			want2: time.Time{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := &User{
				ID:       tt.fields.ID,
				Username: tt.fields.Username,
				Tests:    tt.fields.Tests,
			}
			got, got1, got2 := u.VerifySignature(tt.args.signature)
			if got != tt.want {
				t.Errorf("User.VerifySignature(%v) got = %v, want %v", tt.args.signature, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("User.VerifySignature(%v) got1 = %v, want %v", tt.args.signature, got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("User.VerifySignature(%v) got2 = %v, want %v", tt.args.signature, got2, tt.want2)
			}
		})
	}
}

func TestUser_MarshalBinary(t *testing.T) {
	type fields struct {
		ID       int
		Username string
		Tests    []*Test
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Test MarshalBinary",
			fields: fields{
				ID:       1,
				Username: "test",
				Tests:    []*Test{},
			},
			want:    []byte(`{"ID":1,"Username":"test","Tests":[]}`),
			wantErr: false,
		},
		{
			name: "Test MarshalBinary recursive",
			fields: fields{
				ID:       1,
				Username: "test",
				Tests: []*Test{
					{
						Questions: "test",
						Answers:   "test",
						Signature: "test",
						SignedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want: []byte(
				`{"ID":1,"Username":"test","Tests":[{"Questions":"test","Answers":"test","Signature":"test","SignedAt":"2021-01-01T00:00:00Z"}]}`,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := &User{
				ID:       tt.fields.ID,
				Username: tt.fields.Username,
				Tests:    tt.fields.Tests,
			}
			got, err := u.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.MarshalBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}
