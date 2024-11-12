commit_hash=$(git rev-parse HEAD)
output_file="wakapi.$commit_hash"

go build -ldflags "-X github.com/muety/wakapi/utils.CommitHash=$commit_hash" -o $output_file