#!/bin/bash
# ═══════════════════════════════════════════════════
# Build & Run Helper Script
# ═══════════════════════════════════════════════════

set -e

# Change to the script's directory so docker compose works from any path
cd "$(dirname "$0")"

# Detect if docker requires sudo
DOCKER_COMPOSE="docker compose"
if ! docker ps >/dev/null 2>&1; then
    DOCKER_COMPOSE="sudo docker compose"
fi

function print_usage() {
    echo "Usage: ./run.sh [up|down|restart|clean|build|status]"
    echo "  up      - Start all containers in the background"
    echo "  down    - Stop all running containers"
    echo "  restart - Restart all containers"
    echo "  clean   - Delete the database and all Docker volumes for a clean run"
    echo "  build   - Rebuild all Docker containers"
    echo "  status  - Check status of running containers"
}

if [ -z "$1" ]; then
    print_usage
    exit 1
fi

case "$1" in
    up)
        echo "Starting containers..."
        $DOCKER_COMPOSE up --build -d
        echo "Containers started successfully. Access frontend at http://localhost/"
        ;;
    down)
        echo "Stopping containers..."
        $DOCKER_COMPOSE down
        echo "Containers stopped."
        ;;
    restart)
        echo "Restarting containers..."
        $DOCKER_COMPOSE restart
        echo "Containers restarted."
        ;;
    clean)
        echo "WARNING: This will delete the SQLite database and all uploaded files."
        read -p "Are you sure you want to clean? (y/N) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "Stopping containers and cleaning volumes..."
            $DOCKER_COMPOSE down -v
            echo "Volumes cleaned. Run './run.sh up' to start a clean database populated with seeds."
        else
            echo "Clean aborted."
        fi
        ;;
    build)
        echo "Building containers..."
        $DOCKER_COMPOSE build
        echo "Build complete."
        ;;
    status)
        $DOCKER_COMPOSE ps
        ;;
    *)
        print_usage
        exit 1
        ;;
esac
