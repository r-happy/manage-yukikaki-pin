pipeline {
    agent any

    environment {
        // Set timezone for logs if desired
        TZ = 'Asia/Tokyo'
    }

    stages {
        stage('Prepare Network') {
            steps {
                script {
                    def networkExists = sh(
                        script: 'docker network ls --filter name=^gm-yatai-network$ --format="{{ .Name }}"',
                        returnStdout: true
                    ).trim()

                    if (networkExists == '') {
                        echo 'Creating external network gm-yatai-network...'
                        sh 'docker network create gm-yatai-network'
                    } else {
                        echo 'External network gm-yatai-network already exists.'
                    }
                }
            }
        }

        stage('Build') {
            steps {
                echo 'Building Docker images...'
                sh 'docker compose build frontend'
                sh 'docker compose build backend'
            }
        }

        stage('Deploy') {
            steps {
                echo 'Deploying application with provided credentials...'

                withCredentials([
                    string(credentialsId: 'BACKEND_SECRET_KEY', variable: 'SECRET_KEY')
                ]) {
                    sh 'docker compose up -d --remove-orphans'
                }
            }
        }

        stage('Cleanup') {
            steps {
                echo 'Pruning dangling Docker images...'
                sh 'docker image prune -f'
            }
        }
    }

    post {
        always {
            echo 'Pipeline finished.'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}
