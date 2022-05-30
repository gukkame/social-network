# social-network

social-network made by #

## Installation
```
git clone https://01.kood.tech/git/nimi25820/social-network
```

## Usage
### There are 2 ways to run this project:
#### Manually in your own machine:
- Open 2 terminals
- In first terminal, navigate to ./server and run ```go run server.go```
- In the second terminal, navigate to ./client, run ```npm install```, then ```npm run build```(PS. you need to have npm installed to do this), then run ```node server.js```
- Server should be now running on ```http://localhost:8090/```
#### Using Docker
- Make sure your docker is running. By using ```sudo service docker run``` It should start it up.
- First, run ```chmod u+x dockerbuild.sh``` and then ```./dockerbuild.sh```. This will build two images: one for `server` and one for `client`, puts them both in a container.
- Server should be now running on ```http://localhost:8090/```
- After, you can remove the created images by first, running ```chmod u+x dockerremove.sh``` and then ```./dockerremove.sh```.

## Audit
#### [Audit questions](https://github.com/01-edu/public/tree/master/subjects/social-network/audit#functional)
## Used technologies
#### Main
- [JavaScript](https://www.javascript.com/)
- [Vue.js](https://vuejs.org/)
- [Golang](https://go.dev/)
- [Sqlite](https://www.sqlite.org/index.html)
- [SQL](https://en.wikipedia.org/wiki/SQL)
- [Docker](https://www.docker.com/)
- [Markup languages: HTML5 & CSS](https://en.wikipedia.org/wiki/Markup_language)
#### Libraries & other frameworks
- [Axios](https://axios-http.com/)
- [Bootstrap 5](https://getbootstrap.com/)
- [connect-history-api-fallback](https://www.npmjs.com/package/connect-history-api-fallback)
- [Express](https://expressjs.com/)
- [jQuery](https://jquery.com/)
- [lodash](https://lodash.com/)
- [VeeValidate](https://vee-validate.logaretm.com/v4/)
- [vue-cookies-reactive](https://www.npmjs.com/package/vue-cookies-reactive)
- [Vue Router](https://router.vuejs.org/)
- [yup](https://github.com/jquense/yup)

