#!/usr/bin/groovy

// Build constants
def VERSION_STRING_DEFAULT = "0.3.0" // Note: change version string here after branching a stable branch
def PUSH_TO_PRODUCTION_DEFAULT = false
def PAUSE_AFTER_BUILD_DEFAULT = false
def BUILD_CONFIG_DEFAULT = "Debug"

properties([[$class: 'ParametersDefinitionProperty',
  parameterDefinitions: [
    [$class: 'BooleanParameterDefinition',
      defaultValue: PUSH_TO_PRODUCTION_DEFAULT,
      description: 'This will publish the code on Azure (solarmonitornz.azurewebsites.net/api).\nNote: This step can only be performed on a stable branch. \nChanges will go live when selecting this!', 
      name: 'PUSH_TO_PRODUCTION'],
    [$class: 'StringParameterDefinition',
      defaultValue: VERSION_STRING_DEFAULT,
            description: 'String used to version all binaries on the master branch (along with the build number).', 
      name: 'VERSION_STRING'],
    [$class: 'ChoiceParameterDefinition',
      choices: 'Debug', // ['Debug', 'Release'], 
      description: 'Specifies which configuration to build.',
      name: 'BUILD_CONFIG'],
    [$class: 'BooleanParameterDefinition',
      defaultValue: PAUSE_AFTER_BUILD_DEFAULT,
      description: 'Pause at the end of the Build stage so that the build container used by the build is still available for inspection.',
      name: 'PAUSE_AFTER_BUILD']]]])
     

def newVersion = "${VERSION_STRING}.${env.BUILD_NUMBER}"
def rootDirectory = ""
def outputDirectory = "build"
def releaseBranchPrefix = "rel-"
def productionGitTag = "GA_RELEASE" // General Availability
//def checkErrorResultAndLog = "res=\$? ; if [[ \$res -ne 0 ]] ; then echo \"Error: script failed! Status returned: \$res\" ;  fi"
//def checkErrorResultAndLog = 'res=$? ; if [[ $res -ne 0 ]] ; then echo "Error: script failed! Status returned: $res" ;  fi'

def testSystemIpAddress = "10.0.0.99"
def testSystemSshUser = "sol"

def stagingSystemIpAddress = "10.0.0.99"
def stagingSystemSshUser = "sol"

def productionSystemIpAddress = "pi-solar"
def productionSystemSshUser = "sol"


if (env.BRANCH_NAME == "master") {
    newVersion = "${VERSION_STRING}.${env.BUILD_NUMBER}"
} else if (env.BRANCH_NAME.startsWith(releaseBranchPrefix)) {
    newVersion = env.BRANCH_NAME.replace(releaseBranchPrefix, "") + "." + env.BUILD_NUMBER
} else {
    // other development branches
    newVersion = "0.0.0.${env.BUILD_NUMBER}"
}

def gitScm = null

// ----------------------------------------------------------------------------
// Test variables
// ----------------------------------------------------------------------------
echo "New version: ${newVersion}"
def minorVersion = "0";
if(newVersion.split('\\.').size() > 1) {
    minorVersion = newVersion.split('\\.')[1];
}
// unlikely that we reached such a high number - but just in case 
if (minorVersion.toInteger() >= 255) {
    minorVersion = (minorVersion.toInteger() % 255).toString()
} 
//def testSubnet = "172.18.${minorVersion}.0"
//def testSubnetName = "testnet${minorVersion}"
// ----------------------------------------------------------------------------


def printBuildInfo() {
    echo "\n\n--------------------------------------------------------------------------------------------------------------"
    echo "   Build constants"
    echo "--------------------------------------------------------------------------------------------------------------"
    echo "BRANCH_NAME: ${env.BRANCH_NAME}"
    echo "BUILD_NUMBER: ${env.BUILD_NUMBER}"
    // parameters
    echo "PUSH_TO_PRODUCTION: ${PUSH_TO_PRODUCTION}";
    echo "VERSION_STRING: ${VERSION_STRING}";
    echo "BUILD_CONFIG: ${BUILD_CONFIG}";
    echo "PAUSE_AFTER_BUILD: ${PAUSE_AFTER_BUILD}";
    echo "--------------------------------------------------------------------------------------------------------------\n\n"
}

try {

    printBuildInfo() 

    echo "\n\n========================   STAGE: INIT   =================================================================\n"
    
    // Purpose: sets up global variables and Jenkins UI - this stage only performs steps
    //          on the master node and is not slave specific
    
    stage('Init')
    
        node('master') {
            wrap([$class: 'TimestamperBuildWrapper']) {
                def newBuildName = newVersion
                try {
                    if (newBuildName) {
                        manager.build.displayName = newBuildName
                    }
                    println "Build display name is set to ${newBuildName}"
                } catch (MissingPropertyException e) {
                    echo "Unable to change build display name."
                }
            }
        }
//    }

    echo "\n\n========================   STAGE: CHECKOUT   =============================================================\n"

    // Purpose: checks out source code and stashes it
    //          This is necessary to avoid multiple git checkouts (and to prevent Jenkins
    //          from showing the same commits multiple times in the web interface)

    stage("Checkout")
    
        def sourceCodeStashName = "source_code_stash"
        
        node {
            wrap([$class: 'TimestamperBuildWrapper']) {
                dir ("_checkout") {
                    deleteAllFiles()
                    
                    gitScm = git url: 'gitolite3@ciprian-desktop:SolarMonitorController.git', branch: env.BRANCH_NAME
                    
                    sh "git tag -a ${newVersion} -m '' "
                    sh "git push --follow-tags origin ${newVersion}"
                    
                    // save the version in a local file so that we can print it when extracing
                    // the stash (in the later stages)
                    sh "echo ${newVersion} > ./VERSION"

                    // stash source code as it's needed by the later stages ("build", staging" and "production")
                    stash name: sourceCodeStashName, useDefaultExcludes: false // include everything since this is a clean checkout
                }
            }
        }
//    }
    
    echo "\n\n========================   STAGE: BUILD   ================================================================\n"

    // Purpose: runs the build and related steps.
    stage('Build')
        
        // def unit_tests_stash_prefix = "unit_tests_${BUILD_CONFIG}_${newVersion}"
        // def integration_tests_stash_prefix = "integration_tests_${BUILD_CONFIG}_${newVersion}"
        // def tools_stash_prefix = "tools_${BUILD_CONFIG}_${newVersion}"
        
        parallel 'Linux amd64 build': {
            node("master") {
                wrap([$class: 'TimestamperBuildWrapper']) {
                    deleteAllFiles()
                    unstash sourceCodeStashName

                    def platform = "linux_amd64"
                    
                    // this is not really necessary as it's a clean git repo
                    // leaving it in as confirmation that the dir is indeed clean of binaries
                    sh "scripts/clean_build.sh ."

                    sh "mkdir -p ${outputDirectory}"

                    sh "scripts/update_version.sh ${newVersion}"  
                    sh "scripts/build_all.sh x64 ${BUILD_CONFIG}"
                    sh "scripts/check_version.sh ${newVersion} ${platform}"

                    sh "scripts/verify_source_code.sh"

                    stash name: "solarcmd-x64", 
                        includes: "${outputDirectory}/${platform}/solarcmd, src/cmd/solarcmd/*.config.json"
                    
                    sh "scripts/create_packages.sh ${BUILD_CONFIG} ${newVersion} ${platform}"
                    def archives = "${outputDirectory}/*.tar.gz"
                    step([$class: 'ArtifactArchiver', artifacts: archives, fingerprint: true])
                        
                    // echo "Generating docs."
                    // TODO: generate docs with go doc
                    // publishHTML(target: [
                    //     reportName: "API Documentation",
                    //     reportDir: "output/docfx_site",
                    //     reportFiles: "index.html",
                    //     allowMissing: false,
                    //     alwaysLinkToLaskBuild: false,
                    //     keepAll: true])
                }
            }
        }, 'Linux arm build': {
            node("master") {
                wrap([$class: 'TimestamperBuildWrapper']) {
                    deleteAllFiles()
                    unstash sourceCodeStashName

                    def platform = "linux_arm"
                    
                    // this is not really necessary as it's a clean git repo
                    // leaving it in as confirmation that the dir is indeed clean of binaries
                    sh "scripts/clean_build.sh ."

                    sh "mkdir -p ${outputDirectory}"

                    sh "scripts/update_version.sh ${newVersion}"  
                    sh "scripts/build_all.sh arm ${BUILD_CONFIG}"
                    sh "scripts/check_version.sh ${newVersion} ${platform} ${testSystemIpAddress}"

                    stash name: "solarcmd-arm", 
                        includes: "${outputDirectory}/${platform}/solarcmd, src/cmd/solarcmd/*.config.json"
                    
                    stash name: "scripts_stash"
                        include: "scripts/*"
                        
                    stash name: "tools_stash"
                        include: "tools/*"

                    sh "scripts/create_packages.sh ${BUILD_CONFIG} ${newVersion} ${platform}"
                    def archives = "${outputDirectory}/*.tar.gz"
                    step([$class: 'ArtifactArchiver', artifacts: archives, fingerprint: true])

                    // build tests and stash
                    // bin_dir_expr = "test/server/integration/**/bin/${BUILD_CONFIG}/**"
                    // stash name: "${integration_tests_stash_prefix}_linux",
                    //     includes: "${bin_dir_expr}/*.dll, ${bin_dir_expr}/*.json, test/server/integration/**/project*.json"
                }
            }
        }
//    }
    
    
    echo "\n\n========================   STAGE: PUBLISH   ==============================================================\n"
    
    // Purpose: pushes the binary and other tools to the local test machine.

    stage('Publish') 
        node("master") {
            wrap([$class: 'TimestamperBuildWrapper']) {
                dir("_publish") {
                    deleteAllFiles()
                    unstash "solarcmd-arm"

                    echo "Setting up remote dir."
                    sh "ssh ${testSystemSshUser}@${testSystemIpAddress} mkdir -p /home/sol/test/solarcmd-${newVersion}"

                    echo "Upload executable."
                    sh "scp build/linux_arm/solarcmd ${testSystemSshUser}@${testSystemIpAddress}:test/solarcmd-${newVersion}/"

                    echo "Upload config files."
                    sh "scp src/cmd/solarcmd/*.config.json ${testSystemSshUser}@${testSystemIpAddress}:test/solarcmd-${newVersion}/"

                    echo "Check version of uploaded executable."
                    sh "ssh ${testSystemSshUser}@${testSystemIpAddress} /home/sol/test/solarcmd-${newVersion}/solarcmd -q -v"
                }
            }
        }
//    }
    
    echo "\n\n========================   STAGE: TEST   =================================================================\n"
    
    // Purpose: runs various tests on the local test machine (TBD)
    stage('Test')
        node("master") {
            wrap([$class: 'TimestamperBuildWrapper']) {
                dir("_tests") {
                    deleteAllFiles()
                    unstash sourceCodeStashName
                    sh "mkdir -p ${outputDirectory}"
                    sh "scripts/run_unit_tests.sh"

                    def testResultsPattern = "${outputDirectory}/*test_report.xml"
                    step([$class: 'ArtifactArchiver', artifacts: testResultsPattern, fingerprint: false])
                    
                    stash name: "unit_test_results_stash"
                        include: testResultsPattern
                }
            }
        }
            
    
    echo "\n\n========================   STAGE: STAGING   ==============================================================\n"
    
    // Purpose: deploy application to a long-term test server.
    stage (name: 'Staging', concurrency: 1)
        
        echo "Checking branch name (${env.BRANCH_NAME}) to see if we are allowed to push to staging."
        if (!env.BRANCH_NAME.startsWith(releaseBranchPrefix)) {
            echo "Skipping Staging stage."
        } else {
            node {
                dir ("_staging") {
                    deleteAllFiles()
                    unstash "solarcmd-arm"
                    echo "Upload executable."
                    sh "scp build/linux_arm/solarcmd ${stagingSystemSshUser}@${stagingSystemIpAddress}:bin"
                    echo "Check version of uploaded executable."
                    sh "ssh ${stagingSystemSshUser}@${stagingSystemIpAddress} /home/sol/bin/solarcmd -q -v"
                }
            }
        }
//    }
    
    
    echo "\n\n========================   STAGE: PRODUCTION   ===========================================================\n"
    
    // Purpose: deploy application to the production server (pi-solar machine)
    //          (this stage is almost identical to the staging, except that it's only performed manually)

    stage (name: 'Production', concurrency: 1)

        echo "Checking branch name (${env.BRANCH_NAME}) to see if we are allowed to push to production."
        if (!env.BRANCH_NAME.startsWith(releaseBranchPrefix)) {
            echo "Skipping Production stage."
        } else {
            if(PUSH_TO_PRODUCTION == "false") {
            echo "Production stage SKIPPED."
            } else {
                input 'Confirm deploy to production?'
                node {
                    wrap([$class: 'TimestamperBuildWrapper']) {
                        echo "Pushing to production."
                        dir ("_production") {
                            deleteAllFiles()
                            unstash sourceCodeStashName
                            sh "cat ./VERSION"

                            unstash "solarcmd-arm"
                            echo "Upload executable."
                            sh "scp build/linux_arm/solarcmd ${productionSystemSshUser}@${productionSystemIpAddress}:bin"

                            echo "Check version of uploaded executable."
                            sh "ssh ${productionSystemSshUser}@${productionSystemIpAddress} /home/sol/bin/solarcmd -q -v"

                            echo "Deployment successful. Tagging source code."
                            sh "git tag ${productionGitTag}"
                            sh "git push origin ${productionGitTag}"
                            
                            echo "Deployment successful. Updating status badge."
                            manager.addBadge("green.gif", "GA Release. Deployed to production server.")
                        }
                    }
                }
            }
        }
//    }

} finally {

    echo "\n\n========================   STAGE: FINALISE   =============================================================\n"
    
    // Purpose: cleanup stage at the end of the pipeline
    //          (this stage is always executed, whether the pipeline succeeds or fails)
    
    stage('Finalise')
        node('master') {
            wrap([$class: 'TimestamperBuildWrapper']) {
                deleteAllFiles()
                try {
                    unstash "unit_test_results_stash"
                    // unstash "integration_test_results_stash"
        
                    def testResultsPattern = "${outputDirectory}/*test_report.xml"
                    junit testResultsPattern

                } catch(ex) {
                    echo "[Finalise] Error: failed to retrieve stashes with the test results or to publish them. Was the test stage skipped?"
                    echo "Exception caught: ${ex}"
                }
            
                step([$class: 'LogParserPublisher', 
                    parsingRulesPath: '/var/lib/jenkins/jenkins-rule-logparser', 
                    failBuildOnError: true, 
                    unstableOnWarning: true,
                    showGraphs: true,
                    useProjectRule: false])

                // Note: disabled for now since JiraIssueUpdater doesn't seem to support multi-branch pipeline builds yet
                // See: http://ciprian-desktop:8080/job/SolarMonitor/job/master/34/console
                // exception: java.lang.IllegalArgumentException: Unsupported run type org.jenkinsci.plugins.workflow.job.WorkflowRun
                // Note2: this may not be necessary if we get the Jenkins plugin in JIRA to work (which is also not supporting pipeline builds yet)

                // step([$class: 'hudson.plugins.jira.JiraIssueUpdater', 
                //     issueSelector: [$class: 'hudson.plugins.jira.selector.DefaultIssueSelector'], 
                //     scm: gitScm,
                //     labels: [ "$newVersion", "jenkins" ]])
                //

                gitScm = null
            }
        }
//    }
}    


def deleteAllFiles() {
    sh "ls -laF"
    sh 'for dir in `ls -Ab` ; do rm -rf $dir ; done'
    sh "ls -laF"
    checkDirEmpty()
}

def checkDirEmpty(def pattern = "") {
    sh 'if [ "`ls -A -I \"' + pattern + '\"`" ]; then echo "Error: unable to delete all files in dir: $PWD"; else echo "Empty dir: $PWD" ; fi'
}

def runScript(def script, def label = "", def echoScript = true, def checkResult = true) {
    echo "[${label}]"
    if(echoScript) {
        echo "Running bash script: \"${script}\""
    }
    if(checkResult) {
        sh script + " ; res=\$? ; set +x ; if [ \$res -ne 0 ] ; then echo \"Error: script failed! Status returned: \$res\" ; fi ; exit \$res"
    } else {
        sh script
    }
}
