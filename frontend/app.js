var PATH = []
var TABLE = null
var log = new EventSource(APP_URL+"/events?stream=log");
const d = new Date()
var today = d.getFullYear().toString()+"-"+(d.getMonth()+1).toString().padStart(2,'0')+"-"+d.getDate().toString().padStart(2,'0');
var CAN_RECIEVE_SSE = false
var DELAY_BETWEEN_EVENT = 500
var AUTO_SCROLL_DOWN = true
function changeDelayBetweenEvents(delay){
    DELAY_BETWEEN_EVENT = delay
}
log.onmessage = function(event) {
    if(today != currentdate) return
    mes = JSON.parse(event.data)
    if(CAN_RECIEVE_SSE){
        if(JSON.stringify(mes[0])==JSON.stringify(PATH)){
            handleEvent(mes[1])
        }
        var tableStatics = DatabaseTablesEventCount.get(mes[0][2])
        tableStatics.Total++
        tableStatics.Type[mes[0][1]]++
        //logs[mes[0][2]].push(mes[1])
        DatabaseTablesEventCount.set(mes[0][2],tableStatics)
        if(PATH.length == 2 && PATH[1] == mes[0][1]){
            updateEventCounter(mes[0][2])
            var cTime = new Date;
            var seconds = cTime.getSeconds();
            var minutes = cTime.getMinutes();
            var hour = cTime.getHours();
            document.getElementById("currentEventTime").innerHTML = hour+":"+minutes+":"+seconds
        }
    }
}
log.onerror = function(error) {
    console.error('SSE error:', error);
};

var PRIMARY_LIST = new Map()


function isConnected(state){
    if(state) return `
    <div class="bg-clip-padding backdrop-filter backdrop-blur-sm bg-opacity-20 bg-lime-600 text-white px-2 py-1 text-xs rounded font-semibold flex space-x-1 items-center">
        <svg viewBox="0 0 24 24" fill="currentColor" class="w-3 h-3 text-lime-700">
            <path fill-rule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm13.36-1.814a.75.75 0 10-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 00-1.06 1.06l2.25 2.25a.75.75 0 001.14-.094l3.75-5.25z" clip-rule="evenodd" />
        </svg>                  
        <span>Connected</span>
    </div>`
    return `
    <div class="bg-clip-padding backdrop-filter backdrop-blur-sm bg-opacity-20 bg-red-600 text-white px-2 py-1 text-xs rounded font-semibold flex space-x-1 items-center">
        <svg viewBox="0 0 24 24" fill="currentColor" class="w-3 h-3 text-red-800">
            <path fill-rule="evenodd" d="M12 2.25c-5.385 0-9.75 4.365-9.75 9.75s4.365 9.75 9.75 9.75 9.75-4.365 9.75-9.75S17.385 2.25 12 2.25zm-1.72 6.97a.75.75 0 10-1.06 1.06L10.94 12l-1.72 1.72a.75.75 0 101.06 1.06L12 13.06l1.72 1.72a.75.75 0 101.06-1.06L13.06 12l1.72-1.72a.75.75 0 10-1.06-1.06L12 10.94l-1.72-1.72z" clip-rule="evenodd" />
        </svg>
        <span>Disconnect</span>
    </div>`
}

function createHeader(server){
    return `
        <div>
            <button id="backbtn" onclick="getServers()" class="py-12 pl-1 mr-7 w-16 duration-200 bg-clip-padding backdrop-filter backdrop-blur-xl bg-black bg-opacity-20 hover:bg-opacity-40 rotate-45">
                <svg class="-rotate-45 w-6 h-6 m-auto" viewBox="0 0 24 24" fill="currentColor">
                    <path fill-rule="evenodd" d="M7.72 12.53a.75.75 0 010-1.06l7.5-7.5a.75.75 0 111.06 1.06L9.31 12l6.97 6.97a.75.75 0 11-1.06 1.06l-7.5-7.5z" clip-rule="evenodd" />
                </svg>
            </button>
        </div>
        <div>
            <button class="py-12 pr-2 w-48 duration-200 bg-clip-padding backdrop-filter backdrop-blur-xl bg-black bg-opacity-20 hover:bg-opacity-40 -rotate-45" onclick="getServers()">
                <div class="rotate-45">`+server.server.Name+`</div>
            </button>
        </div>
        <div>
            <button class="py-12 pr-2 w-48 duration-200 bg-clip-padding backdrop-filter backdrop-blur-xl bg-black bg-opacity-20 hover:bg-opacity-40 -rotate-45" onclick="getDatabase('`+server.server.Name+`','`+server.database.Name+`')">
                <div class="rotate-45">`+server.database.Name+`</div>
            </button>
        </div>
    `
}

function createServer(server){
    state = isConnected(server.IsConnected)
    var dbs = ''
    server.Dbs.forEach(
    (val)=>dbs += `
        <button onclick="getDatabase('`+server.Name+`','`+val+`')" class="text-slate-700 font-semibold bg-slate-100 rounded px-2 py-1">
            `+val+`
        </button>
    `)

    return `
        <div class="bg-clip-padding animate-fade backdrop-filter backdrop-blur-xl bg-opacity-20 duration-200 hover:bg-opacity-40 rounded-lg bg-black border-black overflow-hidden">
            <div class="flex px-4 pt-4 items-center">
                <h1 class="grow font-semibold text-xl">
                    `+server.Name+`
                </h1>
                `+state+`
            </div>
            <div class="p-4 space-y-2 text-sm">
                <div class="flex justify-between">
                    <div class="space-x-2 flex items-center opacity-80">
                        <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M21.721 12.752a9.711 9.711 0 00-.945-5.003 12.754 12.754 0 01-4.339 2.708 18.991 18.991 0 01-.214 4.772 17.165 17.165 0 005.498-2.477zM14.634 15.55a17.324 17.324 0 00.332-4.647c-.952.227-1.945.347-2.966.347-1.021 0-2.014-.12-2.966-.347a17.515 17.515 0 00.332 4.647 17.385 17.385 0 005.268 0zM9.772 17.119a18.963 18.963 0 004.456 0A17.182 17.182 0 0112 21.724a17.18 17.18 0 01-2.228-4.605zM7.777 15.23a18.87 18.87 0 01-.214-4.774 12.753 12.753 0 01-4.34-2.708 9.711 9.711 0 00-.944 5.004 17.165 17.165 0 005.498 2.477zM21.356 14.752a9.765 9.765 0 01-7.478 6.817 18.64 18.64 0 001.988-4.718 18.627 18.627 0 005.49-2.098zM2.644 14.752c1.682.971 3.53 1.688 5.49 2.099a18.64 18.64 0 001.988 4.718 9.765 9.765 0 01-7.478-6.816zM13.878 2.43a9.755 9.755 0 016.116 3.986 11.267 11.267 0 01-3.746 2.504 18.63 18.63 0 00-2.37-6.49zM12 2.276a17.152 17.152 0 012.805 7.121c-.897.23-1.837.353-2.805.353-.968 0-1.908-.122-2.805-.353A17.151 17.151 0 0112 2.276zM10.122 2.43a18.629 18.629 0 00-2.37 6.49 11.266 11.266 0 01-3.746-2.504 9.754 9.754 0 016.116-3.985z" />
                        </svg>
                        <h2>Server IP</h2>
                    </div>
                    <h2 class="space-x-1">
                        <span class="font-semibold">`+server.Host+`</span>
                        <span class="font-bold opacity-50">:</span>
                        <span class="font-semibold">`+server.Port+`</span>
                    </h2>
                </div>
                <div class="h-px bg-white opacity-10"></div>
                <div class="flex justify-between">
                    <div class="space-x-2 flex items-center opacity-80">
                        <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M5.507 4.048A3 3 0 017.785 3h8.43a3 3 0 012.278 1.048l1.722 2.008A4.533 4.533 0 0019.5 6h-15c-.243 0-.482.02-.715.056l1.722-2.008z" />
                            <path fill-rule="evenodd" d="M1.5 10.5a3 3 0 013-3h15a3 3 0 110 6h-15a3 3 0 01-3-3zm15 0a.75.75 0 11-1.5 0 .75.75 0 011.5 0zm2.25.75a.75.75 0 100-1.5.75.75 0 000 1.5zM4.5 15a3 3 0 100 6h15a3 3 0 100-6h-15zm11.25 3.75a.75.75 0 100-1.5.75.75 0 000 1.5zM19.5 18a.75.75 0 11-1.5 0 .75.75 0 011.5 0z" clip-rule="evenodd" />
                        </svg>                                  
                        <h1>
                            Databases
                        </h1>
                    </div>
                    <div class="text-xs flex space-x-1">`+dbs+`</div>
                </div>
                <div class="h-px bg-white opacity-10"></div>
                <div class="flex justify-between">
                    <div class="space-x-2 flex items-center opacity-80">
                        <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M11.25 4.533A9.707 9.707 0 006 3a9.735 9.735 0 00-3.25.555.75.75 0 00-.5.707v14.25a.75.75 0 001 .707A8.237 8.237 0 016 18.75c1.995 0 3.823.707 5.25 1.886V4.533zM12.75 20.636A8.214 8.214 0 0118 18.75c.966 0 1.89.166 2.75.47a.75.75 0 001-.708V4.262a.75.75 0 00-.5-.707A9.735 9.735 0 0018 3a9.707 9.707 0 00-5.25 1.533v16.103z" />
                        </svg>                                                                  
                        <h1>
                            Source
                        </h1>
                    </div>
                    <div class="flex space-x-2 text-xs">
                        <div>
                            <small class="leading-none text-slate-400">binlog:</small>
                            <h3 class="leading-none font-semibold">`+server.BinLogFile+`</h3>
                        </div>
                        <div>
                            <small class="leading-none text-slate-400">pos:</small>
                            <h3 class="leading-none font-semibold">`+server.BinLogPos+`</h3>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `
}

function createTable(server,db,table){
    var count = 0
    if(logs[table] !== undefined) count = logs[table].length
    return `
        <button onclick="getTable('`+server+`','`+db+`','`+table.name+`')" id="t`+table.name+`" class="bg-clip-padding backdrop-filter backdrop-blur-xl bg-opacity-20 duration-200 hover:bg-opacity-40 rounded-lg bg-black border-black overflow-hidden p-4">
            <h2 class="font-semibold">`+table.name+`</h2>
            <small>
                <span class="opacity-50">events:</span>
                <span class="font-semibold opacity-50" id="`+table.name+`.eventcounter">`+count+`</span>
            </small>
        </button>
    `
}

function createTableLog(table){
    getPrimarys(TABLE.table.columns)
    var template = '<table class="w-full border-collapse text-xs"><thead><tr>'
    var i = 0
    Object.values(table.columns).forEach(val => {
        var icon = ""
        if(PRIMARY_LIST.get(""+i) !== undefined) icon = `
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-3 h-3">
            <path fill-rule="evenodd" d="M15.75 1.5a6.75 6.75 0 00-6.651 7.906c.067.39-.032.717-.221.906l-6.5 6.499a3 3 0 00-.878 2.121v2.818c0 .414.336.75.75.75H6a.75.75 0 00.75-.75v-1.5h1.5A.75.75 0 009 19.5V18h1.5a.75.75 0 00.53-.22l2.658-2.658c.19-.189.517-.288.906-.22A6.75 6.75 0 1015.75 1.5zm0 3a.75.75 0 000 1.5A2.25 2.25 0 0118 8.25a.75.75 0 001.5 0 3.75 3.75 0 00-3.75-3.75z" clip-rule="evenodd" />
        </svg>
        <span>&nbsp;</span>
        `
        template += `
            <th class="text-left border-b px-2 border-white border-opacity-5 py-2">
                <div class="flex items-center">
                `+icon+`
                    `+val.name+`
                </div>
            </th>
        `
        i++
    });
    template+= '</tr></thead>'

    template += '<tbody id="rows"></tbody>'

    template += '</table>'
    return template
}

function updateEventCounter(table){
    var tableTag = document.getElementById('t'+table)
    var counter = document.getElementById(table+'.eventcounter')
    if(DatabaseTablesEventCount.get(table).Total > 0){
        counter.classList.remove("opacity-50")
        counter.classList.add("font-bold")
    }
    update("t"+table,tableTag.innerHTML)
    update(table+'.eventcounter',DatabaseTablesEventCount.get(table).Total)
    //hotit("t"+table,tableTag.innerHTML)
    //document.getElementById(table+'.eventcounter').innerHTML = logs[table].length
}

function createTableRow(columns,rows){
    Object.values(rows).forEach(r => {
        id = generateRowUniqueId(r)
        var tr = document.createElement('tr')
        tr.classList.add("animate-fade")
        tr.setAttribute('id',"row"+id)
        Object.keys(columns).forEach(key => {
            var td = document.createElement('td')
            td.classList.add("py-1","px-2","truncate","line","select-all")
            td.setAttribute('id',"row"+id+"col"+key)
            td.setAttribute('title',r[key])
            td.setAttribute('style',"max-width:1px")
            td.textContent = r[key]
            
            tr.appendChild(td)
        })
        insert('rows',tr)
    })
}

function updateTableRow(columns,event){
    var toUpdate = {}
    Object.keys(event[3]).forEach(i => {
        var primary = generateRowUniqueId(event[3][i])
        var newPrimary = primary
        // check for primary changes
        var primaryIndexs = [...PRIMARY_LIST.keys()]
        Object.keys(event[2][i]).forEach(col => {
            if(primaryIndexs.includes(col)){
                var oldP = "["+PRIMARY_LIST.get(col)+":"+event[3][i][col]+"]"
                var newP = "["+PRIMARY_LIST.get(col)+":"+event[2][i][col]+"]"
                newPrimary = newPrimary.replace(oldP,newP)
                console.log('oldP: '+oldP+" newP:"+newP)
            }
        })
        let rowElement = document.getElementById('row'+primary)
        if(rowElement !== null){
            document.getElementById('row'+primary).setAttribute('id','row'+newPrimary)
            Object.keys(event[2][i]).forEach(col => {
                var id = rowElement.querySelector('td:nth-child('+(parseInt(col)+1)+')').id
                update(id,event[2][i][col])
            })
        }else{
            primary = newPrimary
            var tr = document.createElement('tr')
            tr.classList.add("animate-fade")
            tr.setAttribute('id','row'+primary)
            r= event[2][i]
            var colId = 0
            Object.entries(columns).forEach((column,key) => {
                var td = document.createElement('td')
                td.classList.add("py-1","px-2","font-semibold","line","select-all")
                td.setAttribute('id','row'+primary+'col'+key)
                td.setAttribute('title',r[key])
                column = column[1]
                if(r[key] === undefined) r[key] = ''
                if(column.key == 'PRI') td.textContent = event[3][i][colId]
                else td.textContent = ""
                tr.appendChild(td)
                if(r[key] != "")
                toUpdate['row'+primary+'col'+key] = r[key]
                colId++
            })
            document.getElementById('rows').appendChild(tr)
        }
        i++
    })
    Object.keys(toUpdate).forEach(i => {
        update(i,toUpdate[i])
    })
}

function deleteTableRow(event){
    Object.values(event[2]).forEach(r => {
        var primary = generateRowUniqueId(r)
        let rowElement = document.getElementById('row'+primary)
        if(rowElement !== null){
            deletes('row'+primary)
        }
    })
}

function changeToolsButtons(elements){
    document.getElementById('buttons').replaceChildren(elements)
}

async function update(id,value){
    var element = document.getElementById(id)
    if(element === null){
        return
    }
    element.classList.add("duration-300","bg-white","bg-opacity-20")
    await sleep(500)
    var element = document.getElementById(id)
    if(element === null){
        return
    }
    element.innerHTML = value
    await sleep(1000)
    //await sleep(1000)
    element.classList.remove("bg-white")
}
async function insert(to,element){
    var target = document.getElementById(to)
    if(target === null){
        return
    }
    element.classList.add("duration-300","bg-white","bg-opacity-20")
    target.appendChild(element)
    await sleep(1500)
    document.getElementById(element.getAttribute('id')).classList.remove("bg-white")
}
async function deletes(id){
    var element = document.getElementById(id)
    //element.classList.add("duration-300","outline")
    element.classList.add("duration-300","bg-red-800","bg-opacity-20")
    await sleep(1500)
    var element = document.getElementById(element.getAttribute('id'))
    element.classList.add("opacity-20","line-through")
    element.classList.remove("bg-red-800")
}
// TODO Improve this hot color function in future
async function hotit(id){
    const element = document.getElementById(id);
    var size = await getShadowSize(element)
    console.log("befor: "+size)

    size+=10;
    element.style.boxShadow = `0px 0px `+size+`px `+(size/2)+`px red`;
    await sleep(1500)

    var size = await getShadowSize(element)
    console.log("after: "+size)

    size-=10;
    element.style.boxShadow = `0px 0px `+size+`px `+(size/2)+`px red`;
    
}
async function getShadowSize(element) {
    const styles = window.getComputedStyle(element);
    const boxShadow = styles.boxShadow;
    
    const regx = /(?:\d+px\s*){2}(\d+)px/;
    const matches = regx.exec(boxShadow);

    if (matches) {
      return parseInt(matches[1]);
    }
    
    return 0; // Default opacity is 1 (fully opaque) if there's no alpha value found
}

const sleep = ms => new Promise(r => setTimeout(r, ms));
function fileSizeSI(a,b,c,d,e){
    return (b=Math,c=b.log,d=1e3,e=c(a)/c(d)|0,a/b.pow(d,e)).toFixed(2)
    +' '+(e?'kMGTPEZY'[--e]+'B':'Bytes')
}
// routes


var logs = {}
var primary = -1
var currentdate = d.getFullYear().toString()+"-"+(d.getMonth()+1).toString().padStart(2,'0')+"-"+d.getDate().toString().padStart(2,'0')
var DatabaseTablesEventCount = new Map()
async function getDatabaseLogs(server,db,date) {
    DatabaseTablesEventCount.clear()
    Object.keys(logs).forEach(table=>{
        DatabaseTablesEventCount.set(table,{"Total":0,"Type":[0,0,0]})
    })
    if(date === undefined) date = currentdate
    else{
        currentdate = date
        update('currentdate','<div class="-rotate-45 line">'+currentdate+'</div>')
    }

    const res = await fetch(APP_URL+'/api/servers/'+server+'/databases/'+db+'/logs/'+date);
    if(res.status == 200){
        var log = await res.json()
    }else{
        var log = []
    }
    
    getDatabaseLogFiles(server,db)
    for (const event of Object.values(log)) {
        if(event[0] != 3){
            var tableStatics = DatabaseTablesEventCount.get(event[1])
            tableStatics.Total++
            tableStatics.Type[event[0]]++
            DatabaseTablesEventCount.set(event[1],tableStatics)
            var time = event[event.length-1]
            document.getElementById("currentEventTime").innerHTML = time.replace(/(\d{2})(\d{2})(\d{2})/, "$1:$2:$3")
        }
    };
    DatabaseTablesEventCount.forEach((v,k)=>{
        if(v.Total > 0) updateEventCounter(k)
    })
    CAN_RECIEVE_SSE = true
    sortTablesByEvents(SORTBY)
}

async function getDatabaseLogFiles(server,db) {
    const res = await fetch(APP_URL+'/api/servers/'+server+'/databases/'+db+'/files');
    var files = await res.json()

    var dates = document.getElementById('dates')
    dates.innerHTML = ''
    if(files == null){
        dates.innerHTML = `
        <div class="flex duration-200 p-4 text-center">
            No Backups
        </div>
        `
        return
    }
    Object.values(files).forEach(file=>{
        if(file['modification'].includes(file['name'])){
            file['modification'] = file['modification'].replace(file['name']+" ","")
        }
        var color = 'text-transparent'
        if(currentdate == file['name']){
            color = ""
        }
        dates.innerHTML = `
        <div class="flex duration-200">
            <a href="/api/servers/`+server+`/databases/`+db+`/sql/`+file['name']+`" target="_balanc" class="flex items-center space-x-2 p-4 border-r hover:bg-opacity-40 duration-200 hover:bg-black border-white/5">
                <svg class="h-5 w-5 opacity-70" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 7.5l3 2.25-3 2.25m4.5 0h3m-9 8.25h13.5A2.25 2.25 0 0021 18V6a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 6v12a2.25 2.25 0 002.25 2.25z" />
                </svg>
                <span>SQL</span>
            </a>     
            <button class="p-4 line flex space-x-2 items-center text-left grow hover:bg-opacity-40 duration-200 hover:bg-black" onclick="getDatabaseLogs('`+server+`','`+db+`','`+file['name']+`')">
                <div class="grow">
                    <h1 class="font-semibold leading-none line">`+file['name']+`</h1>
                    <div class="flex items-center space-x-4">
                        <span class="opacity-80 text-xs font-bold">`+fileSizeSI(file['size'])+`</span>
                        <span class="text-xs">`+file['modification']+`</span>
                    </div>
                </div>
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-4 h-4 duration-200 `+color+`">
                    <path fill-rule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm13.36-1.814a.75.75 0 10-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 00-1.06 1.06l2.25 2.25a.75.75 0 001.14-.094l3.75-5.25z" clip-rule="evenodd" />
                </svg>
            </button>
        </div>
        `+dates.innerHTML
    })
}

function toggleAutoScroll(){
    AUTO_SCROLL_DOWN = !AUTO_SCROLL_DOWN
    document.getElementById("autoScrollDownArrow").classList.toggle("-rotate-90")
}

function handleEvent(event){
    if(TABLE != null && PATH[2] == event[1]){
        if(event[0] == 0){ // insert
            createTableRow(TABLE.table.columns,event[2])
        }else if(event[0] == 1){ // update
            updateTableRow(TABLE.table.columns,event)
        }else if(event[0] == 2){ // delete
            deleteTableRow(event)
        }
        if(AUTO_SCROLL_DOWN){
            window.scrollTo(0, document.body.scrollHeight);
        }
        var time = event[event.length-1]
        document.getElementById("currentEventTime").innerHTML = time.replace(/(\d{2})(\d{2})(\d{2})/, "$1:$2:$3")
    }
}

async function shutdown(){
    try {
        await fetch(APP_URL+'/api/app/shutdown');
    }catch(err){}
    window.close()
}

function generateRowUniqueId(row){
    var uniqueId = ""
    PRIMARY_LIST.forEach((v,k)=>{
        uniqueId += "["+v+":"+row[k]+"]"
    })
    return uniqueId
}

function getPrimarys(columnsList){
    var i = 0
    PRIMARY_LIST = new Map()
    columnsList.forEach(column => {
        if(column.key == "PRI"){
            PRIMARY_LIST.set(i.toString(),column.name)
        }
        i++
    });
}
var SORTBY = "name"
function sortTablesByEvents(by){
    SORTBY = by
    const parentDiv = document.getElementById('tables')
    const childDivs = Array.from(parentDiv.children)
    childDivs.sort((a, b) => {
        var a,b
        if(by == "total"){
            a = parseInt(DatabaseTablesEventCount.get(a.id.slice(1)).Total) || 0
            b = parseInt(DatabaseTablesEventCount.get(b.id.slice(1)).Total) || 0
        }else if(by == "insert"){
            a = parseInt(DatabaseTablesEventCount.get(a.id.slice(1)).Type[0]) || 0
            b = parseInt(DatabaseTablesEventCount.get(b.id.slice(1)).Type[0]) || 0
        }else if(by == "update"){
            a = parseInt(DatabaseTablesEventCount.get(a.id.slice(1)).Type[1]) || 0
            b = parseInt(DatabaseTablesEventCount.get(b.id.slice(1)).Type[1]) || 0
        }else if(by == "delete"){
            a = parseInt(DatabaseTablesEventCount.get(a.id.slice(1)).Type[2]) || 0
            b = parseInt(DatabaseTablesEventCount.get(b.id.slice(1)).Type[2]) || 0
        }else{
            a = a.id.slice(1) || ""
            b = b.id.slice(1) || ""
            return a.localeCompare(b)
        }
        return  b-a;
    });
    childDivs.forEach((childDiv, index) => {
        parentDiv.appendChild(childDiv)
        childDiv.dataset.order = index + 1
    })
    recountEventsTotals()
}
function recountEventsTotals(){
    var total = 0
    var updates = 0
    var inserts = 0
    var deletes = 0
    var sortByName = document.getElementById("sortByName")
    var sortByTotal = document.getElementById("sortByTotal")
    var sortByInserts = document.getElementById("sortByInserts")
    var sortByUpdates = document.getElementById("sortByUpdates")
    var sortByDeletes = document.getElementById("sortByDeletes")
    DatabaseTablesEventCount.forEach((v,k)=>{
        inserts += v.Type[0]
        updates += v.Type[1]
        deletes += v.Type[2]
        total += v.Total
    })
    if(total>0) update("totalCount","("+total+")")
    if(inserts>0) update("insertsCount","("+inserts+")")
    if(updates>0) update("updatesCount","("+updates+")")
    if(deletes>0) update("deletesCount","("+deletes+")")

    if(SORTBY == "name"){
        sortByName.classList.remove("bg-opacity-0")
        sortByName.classList.add("bg-opacity-30")

        sortByTotal.classList.remove("bg-opacity-30")
        sortByTotal.classList.add("bg-opacity-0")
        sortByInserts.classList.remove("bg-opacity-30")
        sortByInserts.classList.add("bg-opacity-0")
        sortByUpdates.classList.remove("bg-opacity-30")
        sortByUpdates.classList.add("bg-opacity-0")
        sortByDeletes.classList.remove("bg-opacity-30")
        sortByDeletes.classList.add("bg-opacity-0")
    }else if(SORTBY == "total"){
        sortByTotal.classList.remove("bg-opacity-0")
        sortByTotal.classList.add("bg-opacity-30")

        sortByName.classList.remove("bg-opacity-30")
        sortByName.classList.add("bg-opacity-0")
        sortByInserts.classList.remove("bg-opacity-30")
        sortByInserts.classList.add("bg-opacity-0")
        sortByUpdates.classList.remove("bg-opacity-30")
        sortByUpdates.classList.add("bg-opacity-0")
        sortByDeletes.classList.remove("bg-opacity-30")
        sortByDeletes.classList.add("bg-opacity-0")
    }else if(SORTBY == "insert"){
        sortByInserts.classList.remove("bg-opacity-0")
        sortByInserts.classList.add("bg-opacity-30")

        sortByName.classList.remove("bg-opacity-30")
        sortByName.classList.add("bg-opacity-0")
        sortByTotal.classList.remove("bg-opacity-30")
        sortByTotal.classList.add("bg-opacity-0")
        sortByUpdates.classList.remove("bg-opacity-30")
        sortByUpdates.classList.add("bg-opacity-0")
        sortByDeletes.classList.remove("bg-opacity-30")
        sortByDeletes.classList.add("bg-opacity-0")
    }else if(SORTBY == "update"){
        sortByUpdates.classList.remove("bg-opacity-0")
        sortByUpdates.classList.add("bg-opacity-30")

        sortByName.classList.remove("bg-opacity-30")
        sortByName.classList.add("bg-opacity-0")
        sortByInserts.classList.remove("bg-opacity-30")
        sortByInserts.classList.add("bg-opacity-0")
        sortByTotal.classList.remove("bg-opacity-30")
        sortByTotal.classList.add("bg-opacity-0")
        sortByDeletes.classList.remove("bg-opacity-30")
        sortByDeletes.classList.add("bg-opacity-0")
    }else if(SORTBY == "delete"){
        sortByDeletes.classList.remove("bg-opacity-0")
        sortByDeletes.classList.add("bg-opacity-30")

        sortByName.classList.remove("bg-opacity-30")
        sortByName.classList.add("bg-opacity-0")
        sortByInserts.classList.remove("bg-opacity-30")
        sortByInserts.classList.add("bg-opacity-0")
        sortByUpdates.classList.remove("bg-opacity-30")
        sortByUpdates.classList.add("bg-opacity-0")
        sortByTotal.classList.remove("bg-opacity-30")
        sortByTotal.classList.add("bg-opacity-0")
    }

}
