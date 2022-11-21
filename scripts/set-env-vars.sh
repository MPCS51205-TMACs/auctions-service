SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)

export PROJECT_DIR_PATH="$(dirname "$SCRIPT_DIR")"

result="${PROJECT_DIR_PATH%"${PROJECT_DIR_PATH##*[!/]}"}" # extglob-free multi-trailing-/ trim
result="${result##*/}"                                    # remove everything before the last /
result=${result:-/}                             # correct for dirname=/ case

export PROJECT_DIR_NAME=$result

echo
echo "SETTING ENVIRONMENT VARIABLES (auctions-service)..."
echo PROJECT_DIR_PATH=$PROJECT_DIR_PATH
echo PROJECT_DIR_NAME=$result
echo

