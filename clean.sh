# /bin/bash
# Remove all unused images not just dangling ones in the vm (minikube)
# eval $(minikube docker-env)

# api auth
docker rmi $(docker images | grep registry.gitlab.com/isaiahwong/cluster/api/accounts) --force 2>/dev/null 
docker rmi $(docker images | grep auth_auth-test) --force 2>/dev/null 
docker rmi $( docker images | grep '<none>') --force 2>/dev/null 

# Deletes dangling Images
docker system prune -f --all
