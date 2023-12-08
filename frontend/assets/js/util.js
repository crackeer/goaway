function getQuery(key, value) {
    let url = new URLSearchParams(window.location.search)
    return url.get(key) || value
}

function getAllQuery() {
    let params = new URLSearchParams(window.location.search)
    let retData = {}
    for (let param of params) {
        retData[param[0]] = param[1]
    }
    return retData
}

function pushStateWith(query) {
    let newURL = window.location.pathname + "?" + httpBuildQuery(query)
    window.history.pushState(null, "", newURL)
}

function nanoid(t) {
    return crypto.getRandomValues(new Uint8Array(t)).reduce(((t, e) => t += (e &= 63) < 36 ? e.toString(36) : e < 62 ? (e - 26).toString(36).toUpperCase() : e > 62 ? "-" : "_"), "")
}

function simpleReload() {
    window.location.reload()
}

function reloadWith(query) {
    window.location.href = window.location.pathname + '?' + httpBuildQuery(query)
}

function jump(path, query) {
    window.location.href = path + '?' + httpBuildQuery(query)
}

function httpBuildQuery(query) {
    let params = new URLSearchParams("")
    Object.keys(query).forEach(k => {
        params.append(k, query[k])
    })
    return params.toString()
}

function cloneObject(data) {
    let raws = JSON.stringify(data)
    return JSON.parse(raws)
}

function parseBookmark(content) {
    let parts = content.split("\n")
    let list = []
    console.log(parts)
    for (var i in parts) {
        if (parts[i].length > 0) {
            let temp = parts[i].split(">")
            if (temp.length == 1) {
                list.push({
                    title: temp[0],
                    href: temp[0]
                })
            } else if (temp.length > 1) {
                list.push({
                    title: temp[1],
                    href: temp[0]
                })
            }
        }
    }
    return list
}

async function uploadFile(files) {
    let retData = []
    for (var i in files) {
        let formData = new FormData();
        formData.append('file', files[i])
        let config = getImageUploadConfig(files[i])
        let dest = config.dir + "/" + config.name
        let data = await axios.put(dest, formData, {
            headers: {
                'proxy': 'upload2dir'
            }
        });
        if (data.status != 200) {
            return new Promise((_, reject) => {
                resolve(data.statusText)
            })
        }
        retData.push({
            url: dest,
            title: config.name,
        })
    }

    return new Promise((resolve, _) => {
        resolve(retData)
    })
}

function getImageUploadConfig(file) {
    let parts = file.type.split('/')
    let ext = ''
    if (parts.length > 1) {
        ext = parts[1]
    }
    let fileName = dayjs().format('HH-mm-ss@') + nanoid(3)
    if (ext.length > 0) {
        fileName = fileName + '.' + ext
    }
    return {
        dir: '/assets/upload/' + dayjs().format('YYYY-MM-DD'),
        name: fileName,
    }
}


function saveFile(data, name) {
    //Blob为js的一个对象，表示一个不可变的, 原始数据的类似文件对象，这是创建文件中不可缺少的！
    var urlObject = window.URL || window.webkitURL || window;
    var export_blob = new Blob([data]);
    var save_link = document.createElementNS("http://www.w3.org/1999/xhtml", "a")
    save_link.href = urlObject.createObjectURL(export_blob);
    save_link.download = name;
    save_link.click();
}

//js 读取文件
function readFiles(file) {
    var reader = new FileReader();//new一个FileReader实例
    if (/text+/.test(file.type)) {//判断文件类型，是不是text类型
        reader.onload = function (result) {
            console.log(result)
        }
        reader.readAsText(file);
    } else if (/image+/.test(file.type)) {//判断文件是不是imgage类型
        reader.onload = function (result) {
            console.log(result)
        }
        reader.readAsDataURL(file);
    }
}

function initByteMarkdownPreview(target, value) {
    let markdown = 'markdown-preview-' + target
    if (window[markdown] != undefined) {
        $("#" + target).html('')
    }
    window[markdown] = new bytemd.Viewer({
        target: document.getElementById(target),
        props: {
            value: value,
            plugins: [
                bytemdPluginGfm(), bytemdPluginHighlight()
            ]
        },
    });
}

function initJSONEditor(target, value) {
    let jsonEditor = 'jsonEditor-' + target
    if (window[jsonEditor] == undefined) {
        window[jsonEditor] = new JSONEditor(document.getElementById("jsoneditor"), {
            "mode": "code",
            "search": true,
            "indentation": 4
        })
    }
    try {
        let jsonValue = JSON.parse(value)
        
        window[jsonEditor].set(jsonValue)
    } catch (e) {
    }
}
