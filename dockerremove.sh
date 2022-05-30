#stops our applications running containers
docker stop clientcontainer
docker stop servercontainer

#Deletes all the containers not running
docker container prune

#Deletes all the images not used by a container
docker image prune -a

echo -e "Removal successful\n" 