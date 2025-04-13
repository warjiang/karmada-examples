#!/usr/bin/env bash

set -euo pipefail
#set -x
cd "$(dirname "${BASH_SOURCE[0]}")"
SHELL_FOLDER=$(pwd)
REPO_ROOT=$(cd ../ && pwd)
WORK_DIR="${REPO_ROOT}/kubernetes-dashboard"

#echo "clone kubernetes dashboard"
#BRANCH_NAME=release/7.10.1
#git clone --depth=1 --branch ${BRANCH_NAME} git@github.com:kubernetes/dashboard.git ${REPO_ROOT}/tmp
#echo "clone finished"

rm -rf ${WORK_DIR}
mkdir -p ${WORK_DIR}
cd ${WORK_DIR}
go mod init github.com/karmada-io/dashboard

mkdir -p api
cp -R ${REPO_ROOT}/tmp/modules/api/ ./api
cp -R ${REPO_ROOT}/tmp/modules/common/ ./common
rm -rf ./api/go.{mod,sum}
for file in $(find ${WORK_DIR}/api/* -type f -name "*.go"|| sed 's|//*|/|g' ); do
  sed -i "" "s|k8s.io/dashboard/api/pkg/|kubernetes-dashboard/api/pkg/|g" $file
done

# add the following lines in go.mod manually
#replace (
#	k8s.io/dashboard/certificates => ./common/certificates
#	k8s.io/dashboard/client => ./common/client
#	k8s.io/dashboard/csrf => ./common/csrf
#	k8s.io/dashboard/errors => ./common/errors
#	k8s.io/dashboard/helpers => ./common/helpers
#	k8s.io/dashboard/types => ./common/types
#)