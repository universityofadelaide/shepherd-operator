#!/bin/bash

SEARCH_FORMAT="@expire %s"
DATE_FORMAT="+%b %Y"
DIRS="./internal/ ./pkg/"
SEARCH_LAST_N_MONTHS=4

# Cross-platform date formatting with a month offset.
case `uname` in
  Darwin)
    function date_offset_month() {
      date -v $1m "$DATE_FORMAT";
    }
    ;;
  Linux)
    function date_offset_month() {
      date --date="$1 month" "$DATE_FORMAT"
    }
    ;;
  *)
esac

for i in $(seq 0 $SEARCH_LAST_N_MONTHS); do
    FORMATTED_DATE=$(date_offset_month -$i)
    SEARCH_STRING=$(printf "$SEARCH_FORMAT" "$FORMATTED_DATE")
    echo "Searching codebase for \"$SEARCH_STRING\"."
    grep -rni "$SEARCH_STRING" $DIRS && exit 1
done

exit 0
