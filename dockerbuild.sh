#!/bin/bash
echo -e "\nBuilding docker image serverimage\n"
docker image build -f Dockerfile -t serverimage .
echo -e "\nServerimage built successfully"

cd client/

echo -e "\nBuilding docker image clientimage\n"
docker image build -f Dockerfile -t clientimage .
echo -e "\nClientimage built successfully"

echo -e "\nStarting container servercontainer with serverimage at port :8080\n"
docker container run -p 8080:8080 --detach --name servercontainer serverimage
echo -e "\nContainer servercontainer running at port :8080"


echo -e "\nStarting container clientcontainer with clientimage at port :8090\n"
docker container run -p 8090:8090 --detach --name clientcontainer clientimage
echo -e "\nContainer clientcontainer running at port :8090" 

echo -e "\nOpening http://localhost:8090/\n"
xdg-open http://localhost:8090/
