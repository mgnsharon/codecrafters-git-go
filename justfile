set positional-arguments

@ls-tree hash:
    ./your_git.sh $0 $1

@write-tree:
    ./your_git.sh $0

@dhash hash:
    rm -rf .git/objects/$1
create-testfiles:
    set -e
    echo "hello world" > ./cmd/testdata/test_file_1.txt
    mkdir -p ./cmd/testdata/test_dir_1
    echo "hello world" > ./cmd/testdata/test_dir_1/test_file_2.txt
    mkdir -p ./cmd/testdata/test_dir_2
    echo "hello world" > ./cmd/testdata/test_dir_2/test_file_3.txt

clean-testfiles:
    rm -f ./cmd/testdata/test_file_1.txt
    rm -rf ./cmd/testdata/test_dir_1
    rm -rf ./cmd/testdata/test_dir_2