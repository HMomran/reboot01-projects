# 🐳 Docker Compose Orchestration Guide (`run.sh`)

This project uses Docker Compose to run the Go backend and Vue 3 frontend in a containerized environment. A helper script is provided to simplify container management.

> [!IMPORTANT]
> **CRITICAL: You MUST execute `run.sh` with root privileges (`sudo`).**
> On this environment, the Docker daemon socket (`/var/run/docker.sock`) requires root/sudo access. Running commands like `./run.sh up` without `sudo` will fail with a permission error:
> ```bash
> unable to get image: permission denied while trying to connect to the docker API
> ```
> Always prefix your commands with `sudo` (e.g., `sudo ./run.sh up`).

---

## 📋 Available Commands

### 1. Start all containers
Builds the docker images and starts the frontend and backend services in the background.
```bash
sudo ./run.sh up
```
- **Frontend**: Available at [http://localhost/](http://localhost/) (Port 80)
- **Backend API**: Running at [http://localhost:8080](http://localhost:8080) (Port 8080)

### 2. Stop all containers
Stops and removes all running containers for the application.
```bash
sudo ./run.sh down
```

### 3. Restart all containers
Restarts the active containers without rebuilding or recreating them.
```bash
sudo ./run.sh restart
```

### 4. Rebuild all containers
Forces Docker to rebuild the container images from the dockerfiles.
```bash
sudo ./run.sh build
```

### 5. Check container status
Lists the running containers and shows their current health/ports.
```bash
sudo ./run.sh status
```

### 6. Clean Database & Volumes (Reset)
Stops containers and deletes the sqlite database and all media upload volumes. Use this to restore a clean database state populated with initial seeds.
```bash
sudo ./run.sh clean
```
*Note: This command will ask for confirmation before erasing data.*

---

## 🛠 Troubleshooting

### "obsolete version attribute" Warning
You might see the following warning when running compose commands:
```text
WARN[0000] .../docker-compose.yml: the attribute `version` is obsolete, it will be ignored
```
This is expected behavior in modern Docker Compose versions. The version key is obsolete, but it is kept for backwards compatibility and can be safely ignored.
