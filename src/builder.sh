CWD=$(pwd)

build_project() {
  local dir=$1
  local output=$2
  cd "$dir" || exit
  go build -o "$output"
  cd "$CWD" || exit
}

build_project client client
build_project server server
build_project torrentCLI cli