1. Start infra once:
docker compose up -d elasticsearch kafka kibana kafbat-ui

2. For each Go API code change, rebuild only API:
docker compose up -d --build --no-deps api

check log 
docker compose logs -f api

3. Stop all services:
docker compose down