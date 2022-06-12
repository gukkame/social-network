//express - Node.js server library
const express = require('express');

//connect-history-api-fallback - Needed for Vue to understand to always go back to index.html, as it is a SPA
const history = require('connect-history-api-fallback');
const app = express();

//This disables a dot rule in the URL. `.` is now allowed in the url (for category name - dont ask me why.)
app.use(history({
  disableDotRule: true,
  htmlAcceptHeaders: ['text/html']
}));

//Serves the /dist folder to the server
app.use(express.static(__dirname + "/dist"));
app.get('/', (req, res) => {
  res.sendFile('src/index.html');
});

//Express server 
app.listen(8090, () => console.log('Frontend Server running at port http://localhost:8090/'));