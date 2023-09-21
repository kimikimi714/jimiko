package log

import "testing"

func Test_entry_String(t *testing.T) {
	type fields struct {
		Message  string
		Severity string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "info",
			fields: fields{Severity: "INFO", Message: "test"},
			want:   `{"message":"test","severity":"INFO"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := entry{
				Message:  tt.fields.Message,
				Severity: tt.fields.Severity,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("entry.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
