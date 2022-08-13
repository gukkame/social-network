export { createCookie }

function createCookie(Id, username) {
    let d = new Date();
    d.setTime(d.getTime() + 1000 * 3600 * 24 * 365);
    let expires = "expires=" + d.toUTCString();
    document.cookie = "Token=" + Id +":"+ username + ";" + expires + ";path=/"
}