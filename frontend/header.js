var header = `
<nav class="navbar navbar-default">
    <div class="container-fluid">
        <div class="navbar-header">
            <button type="button" class="navbar-toggle collapsed" data-toggle="collapse"
                data-target="#navbar-collapse-1" aria-expanded="false">
                <span class="sr-only">Toggle navigation</span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a class="navbar-brand" href="/">
                <img src="/assets/images/logo.svg" alt="" />
            </a>
        </div>
        <div class="collapse navbar-collapse" id="navbar-collapse-1">
            <ul class="nav navbar-nav">
                <li><a href="/router/list.html">路由</a></li>
                <li><a href="/service/list.html">服务</a></li>
            </ul>
        </div>
    </div>
</nav>
`

var styleFiles = [
    "/assets/css/bootstrap3.4.min.css",
    "/assets/css/my.css",
    "/assets/bytemd/bytemd.css",
    "/assets/bytemd/github-markdown.css",
    "/assets/bytemd/highlight.css",
    "/assets/css/jsoneditor9.8.css"
]
var jsFile1 = [
    "/assets/js/jquery.js",
    "/assets/js/vue.global.js",
    "/assets/js/util.js",
    "/assets/js/axios.min.js",
    "/assets/js/dayjs.min.js",
    "/assets/bytemd/bytemd.umd.js",
    "/assets/js/jsoneditor9.8.js"
]
var jsFile2 = [
    "/assets/js/bootstrap.min.js",
    "/assets/js/bootbox.min.js",
    "/assets/bytemd/bytemd-plugin-gfm.js",
    "/assets/bytemd/plugin-highlight.js",
]
document.addEventListener("DOMContentLoaded", async () => {
    loadStyles(styleFiles)
    await loadJs(jsFile1)
    loadNavigation()
    await loadJs(jsFile2)
    await sleep(200)
    startWork()

}, false);

function loadNavigation() {
    if(window.location.href.indexOf('/user/login.html') !== -1) {
        return
    }
    $('body').prepend(header)
    let parts = window.location.pathname.split('/')
    console.log(parts[1])
    $('a[id="' + parts[1] + '-a"]').parent().addClass('active')
}


async function loadStyles(urls) {
    var head = document.getElementsByTagName("head")[0];
    for (var i = 0; i < urls.length; i++) {
        head.appendChild(createStyleNode(urls[i]));
    }
}

function createStyleNode(url) {
    var link = document.createElement("link");
    link.type = "text/css";
    link.rel = "stylesheet";
    link.href = url;
    return link
}

async function loadJs(urls) {
    //var head = document.getElementsByTagName("head")[0];
    for (var i = 0; i < urls.length; i++) {
        await loadJsUrl(urls[i])
    }
}

function loadJsUrl(url) {
    return new Promise((resolve) => {
        let domScript = createJsNode(url)
        domScript.onload = domScript.onreadystatechange = function () {
            if (!this.readyState || 'loaded' === this.readyState || 'complete' === this.readyState) {
                resolve()
            }
        }
        document.getElementsByTagName('head')[0].appendChild(domScript);
    });
}

function createJsNode(url) {
    var scriptNode = document.createElement("script");
    scriptNode.src = url;
    return scriptNode
}

function sleep(time) {
    return new Promise((resolve) => {
        setTimeout(() => {
            resolve();
        }, time);
    });
}