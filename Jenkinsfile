pipeline {
    agent any

    environment {
        // Set timezone for logs if desired
        TZ = 'Asia/Tokyo'
        DEFAULT_FRONTEND_BACKEND_URL = 'http://manage-yukikaki-backend:1323'
        DOCKER_NETWORK_NAME = 'manage-yukikaki-net'
    }

    stages {
        stage('Prepare Network') {
            steps {
                script {
                    def networkExists = sh(
                        script: "docker network ls --filter name=^${env.DOCKER_NETWORK_NAME}\\$ --format=\"{{ .Name }}\"",
                        returnStdout: true
                    ).trim()

                    if (networkExists == '') {
                        echo "Creating external network ${env.DOCKER_NETWORK_NAME}..."
                        sh "docker network create ${env.DOCKER_NETWORK_NAME}"
                    } else {
                        echo "External network ${env.DOCKER_NETWORK_NAME} already exists."
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
                    script {
                        def targetBackendUrl = env.FRONTEND_BACKEND_URL?.trim()
                        if (!targetBackendUrl) {
                            targetBackendUrl = env.DEFAULT_FRONTEND_BACKEND_URL
                        }

                        sh """
                            export FRONTEND_BACKEND_URL='${targetBackendUrl}'
                            docker compose up -d --remove-orphans
                        """
                    }
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
