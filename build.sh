commit_hash=$(git rev-parse --short=8 HEAD)
branch_name=$(git rev-parse --abbrev-ref HEAD | sed 's/[^a-zA-Z0-9]/_/g')
output_file="wakapi.$branch_name.$commit_hash"

go build -ldflags "-X github.com/muety/wakapi/utils.CommitHash=$commit_hash" -o $output_file