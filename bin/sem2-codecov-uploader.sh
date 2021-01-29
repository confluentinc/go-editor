#!/bin/bash
set -e

# This script sets environment variables for code coverage
# uploader script that match semaphore2 environment.
#
# Semaphore docs for env vars:
# https://docs.semaphoreci.com/ci-cd-environment/environment-variables/
#
# The codecov bash uploader can be found at:
# https://github.com/codecov/codecov-bash

echo "### Starting sem2-codecov-uploader.sh script"

case "${SEMAPHORE_GIT_REF_TYPE}" in

"pull-request")
    echo "Detected pull request."
    export VCS_COMMIT_ID="${SEMAPHORE_GIT_PR_SHA}"
    export VCS_BRANCH_NAME="${SEMAPHORE_GIT_PR_BRANCH}"
    export VCS_PULL_REQUEST="${SEMAPHORE_GIT_PR_NUMBER}"
    export VCS_SLUG="${SEMAPHORE_GIT_PR_SLUG}"
    unset VCS_TAG
    export CI_BUILD_ID="${SEMAPHORE_BUILD_NUMBER}"
    export CI_JOB_ID="${SEMAPHORE_JOB_ID}"
    ;;

"branch")
    echo "Detected branch commit."
    export VCS_COMMIT_ID="${SEMAPHORE_GIT_SHA}"
    export VCS_BRANCH_NAME="${SEMAPHORE_GIT_BRANCH}"
    unset VCS_PULL_REQUEST
    export VCS_SLUG="${SEMAPHORE_GIT_REPO_SLUG}"
    unset VCS_TAG
    export CI_BUILD_ID="${SEMAPHORE_BUILD_NUMBER}"
    export CI_JOB_ID="${SEMAPHORE_JOB_ID}"
    ;;

## TODO: this doesn't take care of CI builds for tags

*)
    echo "Unexpected SEMAPHORE_GIT_REF_TYPE value"
    ;;
esac


## For complex workflows, we may test using multiple jobs / OSes.
export CODECOV_ENV=SEMAPHORE_AGENT_MACHINE_OS_IMAGE,SEMAPHORE_JOB_NAME

echo "### Default env variables used by codecov bash uploader"
echo VCS_COMMIT_ID="${VCS_COMMIT_ID}"
echo VCS_BRANCH_NAME="${VCS_BRANCH_NAME}"
echo VCS_PULL_REQUEST="${VCS_PULL_REQUEST}"
echo VCS_SLUG="${VCS_SLUG}"
echo VCS_TAG="${VCS_TAG}"
echo CI_BUILD_ID="${CI_BUILD_ID}"
echo CI_JOB_ID="${CI_JOB_ID}"
echo CODECOV_ENV="${CODECOV_ENV}"
echo SEMAPHORE_AGENT_MACHINE_OS_IMAGE="${SEMAPHORE_AGENT_MACHINE_OS_IMAGE}"
echo SEMAPHORE_JOB_NAME="${SEMAPHORE_JOB_NAME}"

## The environment vars are named differently in semaphore classic vs sem2
## but the detection uses the same SEMAPHORE=true check.
unset SEMAPHORE

echo "### Calling https://codecov.io/bash script"
## NOTE: this the recommended way to upload coverage reports.
bash <(curl -s https://codecov.io/bash) "$@"
