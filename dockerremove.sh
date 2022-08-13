#stops our applications running containers
docker stop clientcontainer
docker stop servercontainer

#Deletes both the containers
docker rm clientcontainer
docker rm servercontainer

#Deletes both the images
docker rmi serverimage
docker rmi clientimage
docker rmi 36fad710e29d
docker rmi a09e5e1306fb

echo -e "Removal successful\n" 