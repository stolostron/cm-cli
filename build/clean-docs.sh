#! /bin/bash
# Copyright Contributors to the Open Cluster Management project

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Fix sed issues on Mac by using gsed
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
SED="sed"
if [ "${OS}" == "darwin" ]; then
    SED="gsed"
    # Add /usr/local/bin to PATH if it's missing
    if [[ ":${PATH}:" != *":/usr/local/bin:"* ]]; then
      export PATH=${PATH}:/usr/local/bin
    fi
    if [ ! -x "$(command -v ${SED})"  ]; then
       echo "ERROR: ${SED} required, but not found."
       echo "Perform \"brew install gnu-sed\" and try again."
       exit 1
    fi
fi

echo "Cleaning doc files:"
echo "* Replacing '${HOME}' with '\${HOME}'"
for FILE in $(grep -rl "${HOME}" ${DIR}/../docs); do 
  ${SED} -i 's%'${HOME}'%${HOME}%g' ${FILE}
done
