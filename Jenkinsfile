// The build image is where we build all our go binaries.
// In order to allow the release to work appropriately, we build this before all the others.
uacf_generic_pipeline {
    publish_branches = []
    images = [
            build: [
                    DOCKER_IMAGE_NAME: "infra/libra/build",
                    DOCKER_CONTEXT_PATH: ".",
                    DOCKER_FILE: "./build/Dockerfile.build",
            ]
    ]
    tests = []
    last_pipeline=false
}

// Put this after the precursor build above, so that we dont' get alerted when it is built.
env.SLACK_ROOM_FINISH = "infra-team-events"

uacf_generic_pipeline {
    publish_branches = ['master', 'feature/']
    images = [
            // libra uses the PREBUILD_SCRIPT function
            // to run a script to prepare the workspace -before- running a docker build.
            libra: [
                    DOCKER_IMAGE_NAME: "infra/libra",
                    DOCKER_CONTEXT_PATH: ".",
                    DOCKER_FILE: "./Dockerfile",
                    PREBUILD_SCRIPT: "./build/extract-binaries.sh"
            ]
    ]
    tests = []
    first_pipeline=false
}