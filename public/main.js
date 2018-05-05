const $ = document.querySelector.bind(document);
const $all = document.querySelectorAll.bind(document);

(function () {
    var sshareKey = localStorage.getItem("sshareKey");

    if (sshareKey == null) {
        showAuthDialog();
    }
})()

function showAuthDialog() {
    $("#auth-dialog").style["display"] = "block";
}
