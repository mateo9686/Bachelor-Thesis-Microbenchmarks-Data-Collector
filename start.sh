destination_folder="results"

rm -rf /tmp/go-build*
# Check if the source file exists
if [ -e ./log ]; then
    # Create the destination directory if it doesn't exist
    mkdir -p ./"$destination_folder"

    cat ./log >> ./"$destination_folder"/log
    echo "Content from ./log copied and appended to ./$destination_folder/log"
else
    echo "Source file ./log not found. Nothing to copy"
fi

script -c "GOMEMLIMIT=6GiB go run ." log
