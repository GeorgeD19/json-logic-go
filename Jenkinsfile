pipeline {
    agent {
        dockerfile {
            label 'shared'
            filename 'build/Dockerfile'
            additionalBuildArgs '--target tool-env'
        }
    }
    options {
        ansiColor('xterm')
        timeout(time: 20, unit: 'MINUTES')
        timestamps()
        buildDiscarder(logRotator(
            daysToKeepStr: '60',
            numToKeepStr: '300',
            artifactDaysToKeepStr: '5',
            artifactNumToKeepStr: '20'))
    }
    triggers {
        issueCommentTrigger('.*test this please.*')
        pollSCM('')
    }
    environment {
        MISE_TRUSTED_CONFIG_PATHS = "${env.WORKSPACE}/.mise.toml"
    }
    stages {
        stage('Setup') {
            steps {
                echo 'Setup'
                setBuildNameToGitCommit()
                stopOngoingPrBuilds()
            }
        }

        stage('Push Go Module') {
          when {
              buildingTag()
          }
          steps {
            releasePrivateGoModule(
                repository: "go-hosted-private",
                apiKeyId: "cloudsmith-go-hosted-private",
                version: vintedExtractModuleVersion(getVersionFromTag:true),
                skipPublishGitTag: true
            )
          }
        }
    }
}
