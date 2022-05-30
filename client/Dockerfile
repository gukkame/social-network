# Base image
FROM node:17.6

# Selects everything in current directory and copies to /clientcontainer
COPY ./ /clientcontainer

# Move to working directory /clientcontainer
WORKDIR /clientcontainer

# Install the necessary dependencies 
RUN npm install

# Compiles the applicaton into executable /dist
RUN npm run build

# Expose necessary port
EXPOSE 8090

# Command to run when starting the container
CMD [ "node", "server.js"]
