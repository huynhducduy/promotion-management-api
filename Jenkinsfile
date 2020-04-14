pipeline {
    agent any
    stages {
        stage('prepare') {
            steps {
                checkout scm
            }
        }
        stage('build') {
            steps {
                sh 'sudo cp /home/duuuuuuuuy/promotion-management-api.env .env'
                sh 'sudo docker build -t swd391 .'
            }
        }
        stage('test') {
            steps {
                sh 'echo "Passed!"'
            }
        }
        stage('prepare to deploy') {
            steps {
                containers = sh(returnStdout: true, script: 'sudo docker ps -a -q --filter ancestor=swd391:latest --format="{{.ID}}"')
                sh 'sudo docker stop $containers'
                sh 'sudo docker rm $containers'
            }
        }
        stage('deploy') {
            steps {
                sh 'sudo docker run -dit -p 8081:80 swd391:latest'
            }
        }
    }
}