package sb

import (
	"os"

	supabase "github.com/nedpals/supabase-go"
)

const BaseAuthURL = "https://bvutljexdthpferjwsgz.supabase.co/auth/v1/recover"

// https://<project_ref>.supabase.co/rest/v1/
const ResetPasswordEndpoint = "auth/v1/recover"

var Client *supabase.Client

func Init() error {
	sbHost := os.Getenv("SUPABASE_URL")
	sbSecret := os.Getenv("SUPABASE_SECRET")
	Client = supabase.CreateClient(sbHost, sbSecret)
	return nil
}
