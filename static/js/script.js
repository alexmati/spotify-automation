function createPlaylist() {
    const http = new XMLHttpRequest();
    http.onload = function() {
        document.getElementById("button").innerHTML = this.responseText;
    }
    http.open("GET", "index.html", true); //404 page not found
    http.send();
}