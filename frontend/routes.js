async function getServers() {
    var dash = document.getElementById('dashboard')
    const result = await fetch(APP_URL+'/api/servers');
    const servers = await result.json()
    dash.innerHTML = `
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4" id="servers"></div>
    `

    Object.keys(servers).forEach(key => {
        var serversContainer = document.getElementById('servers')
        serversContainer.innerHTML += createServer(servers[key])
    });
    await sleep(2000)
    document.getElementById('on').classList.add("opacity-0")
    await sleep(300)
    document.getElementById('on').classList.add("hidden")

}

function collapseFiles(){
    var files = document.getElementById('files')
    if(files.classList.contains("w-0")){
        document.getElementById('body').classList.add('space-x-4')
        files.classList.remove("w-0","hidden","pr-4")
    } else {
        document.getElementById('body').classList.remove('space-x-4')
        files.classList.add("w-0","hidden","pr-4")
    }
}

async function getDatabase(server,db) {
    PATH[0] = server
    PATH[1] = db
    PATH.splice(2)
    TABLE = null
    var dash = document.getElementById('dashboard')
    //dash.innerHTML = loading()
    const result = await fetch(APP_URL+'/api/servers/'+server+'/databases/'+db);
    const res = await result.json()

    var loadingDivs = ""
    for (let index = 0; index < 15; index++) {
        loadingDivs += `
            <div class="bg-clip-padding animate-pulse backdrop-filter backdrop-blur-xl bg-opacity-20 duration-200 hover:bg-opacity-40 rounded-lg bg-black border-black overflow-hidden p-4">
                <h2 class="font-semibold text-transparent">--</h2>
                <small>
                    <span class="opacity-50 text-transparent">events:</span>
                    <span class="font-semibold opacity-50 text-transparent">--+--</span>
                </small>
            </div>
        `
    }
    dash.innerHTML = `
        <div class="flex grow items-center text-center animate-fade h-8 rounded-lg overflow-hidden no-scrollbar line">
            <div class="overflow-x-scroll overflow-y-hidden flex text-sm grow -space-x-4 items-center" id="navigation">
            </div>
            <button id="currentdate" onclick="collapseFiles()" class="text-xs font-semibold py-20 w-28 bg-clip-padding backdrop-filter backdrop-blur-xl bg-black bg-opacity-20 rotate-45">
                <div class="-rotate-45 line">`+currentdate+`</div>
            </button>
        </div>

        <div class="flex md:space-x-4 pb-12 overflow-x-hidden no-scrollbar" id="body">
            <div class="hidden w-0 duration-300 md:block md:w-auto" id="files">
                <h2 class="px-4 py-2 text-sm">Events per day:</h2>
                <div class="animate-fade bg-clip-padding backdrop-filter backdrop-blur-xl bg-opacity-20 rounded-lg bg-black overflow-hidden" id="dates"></div>
            </div>
            <div class="grow animate-fade space-y-4">
                <div class="md:flex items-center md:space-x-2">
                    <h2 class="px-4 text-sm grow hidden md:block">Tables</h2>
                    <h2 class="text-sm mb-4 md:mb-0">Sort by DESC</h2>
                    <div class="grid grid-cols-5 divide-x divide-white/5 text-xs overflow-hidden text-white duration-200 bg-clip-padding backdrop-filter backdrop-blur-sm bg-opacity-20 rounded-lg bg-black">
                        <button onclick="sortTablesByEvents('name')" id="sortByName" class="py-2 px-4 bg-black hover:bg-opacity-40 bg-clip-padding backdrop-filter backdrop-blur-sm bg-opacity-0 duration-200">
                            Name
                        </button>
                        <button onclick="sortTablesByEvents('total')" id="sortByTotal" class="py-2 px-4 bg-black hover:bg-opacity-40 bg-clip-padding backdrop-filter backdrop-blur-sm bg-opacity-0 duration-200">
                            <h1>Total</h1>
                            <small id="totalCount" class="text-opacity-50 font-semibold">(0)</small>
                        </button>
                        <button onclick="sortTablesByEvents('insert')" id="sortByInserts" class="py-2 px-4 bg-black hover:bg-opacity-40 bg-clip-padding backdrop-filter backdrop-blur-sm bg-opacity-0 duration-200">
                            <h1>Inserts</h1>
                            <small id="insertsCount" class="text-opacity-50 font-semibold">(0)</small>
                        </button>
                        <button onclick="sortTablesByEvents('update')" id="sortByUpdates" class="py-2 px-4 bg-black hover:bg-opacity-40 bg-clip-padding backdrop-filter backdrop-blur-sm bg-opacity-0 duration-200">
                            <h1>Updates</h1>
                            <small id="updatesCount" class="text-opacity-50 font-semibold">(0)</small>
                        </button>
                        <button onclick="sortTablesByEvents('delete')" id="sortByDeletes" class="py-2 px-4 bg-black hover:bg-opacity-40 bg-clip-padding backdrop-filter backdrop-blur-sm bg-opacity-0 duration-200">
                            <h1>Deletes</h1>
                            <small id="deletesCount" class="text-opacity-50 font-semibold">(0)</small>
                        </button>
                    </div>
                </div>
                <div class="grow grid min-w-full grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-5 2xl:grid-cols-6 gap-4 animate-fade line" id="tables">
                    `+loadingDivs+`
                </div>
            </div>
        </div>
    `
    
    getDatabaseLogs(server,db,currentdate)
    document.getElementById('navigation').innerHTML += createHeader(res)
    var tables = ""
    Object.keys(res.database.tables).forEach(table => {
        tables += createTable(server,db,res.database.tables[table])
        logs[table] = []
        DatabaseTablesEventCount.set(table,{"Total":0,"Type":[0,0,0]})
    });
    document.getElementById('tables').innerHTML = tables
    sortTablesByEvents(SORTBY)
}

async function getTable(server,db,table) {
    PATH[0] = server
    PATH[1] = db
    PATH[2] = table
    var dash = document.getElementById('body')
    //dash.innerHTML = loading()
    const result = await fetch(APP_URL+'/api/servers/'+server+'/databases/'+db+'/tables/'+table);
    TABLE = await result.json()
    dash.innerHTML = `
        <div class="w-full animate-fade overflow-x-scroll overflow-hidden no-scrollbar rounded-lg p-4 bg-clip-padding backdrop-filter backdrop-blur-xl bg-opacity-20 duration-200 bg-black font-semibold" id="table"></div>
    `
    var div = document.createElement('div')
    var btn = document.createElement('button')
    btn.classList.add("py-12","w-48","pr-2","duration-200","bg-clip-padding","backdrop-filter","backdrop-blur-xl","bg-black","bg-opacity-20","hover:bg-opacity-40","-rotate-45")
    var text = document.createElement('div')
    text.classList.add("rotate-45")
    text.textContent = table
    btn.appendChild(text)
    //btn.setAttribute("id","t"+table)
    btn.onclick = function(){
        getTable(server,db,table)
    }
    document.getElementById('backbtn').onclick = function(){
        getDatabase(server,db)
    }

    var check = document.getElementById('t'+table)
    if(check === null){
        div.appendChild(btn)
        document.getElementById('navigation').appendChild(div)
    }
    document.getElementById('table').innerHTML += createTableLog(TABLE.table)
    for(const index of Object.keys(TABLE.table.columns)){
        if(TABLE.table.columns[index].key == "PRI") primary = index
        break
    }
    
    const gg = await fetch(APP_URL+'/api/servers/'+server+'/databases/'+db+'/tables/'+table+'/logs/'+currentdate);
    var events = await gg.json()

    var percentage = 0
    var eventsNum = events.length
    var index = 1
    var auther = document.getElementById("auther")
    for (const event of events) {
        auther.innerText = percentage+"%"
        await sleep(DELAY_BETWEEN_EVENT)
        handleEvent(event)
        percentage = Math.ceil(index / eventsNum * 100)
        index++

        if(PATH.length != 3){
            auther.innerText = "Thabet.dev"
            return
        }
    }
    auther.innerText = "Thabet.dev"
}

function settings(){
    var settings = document.getElementById('settings')
    settings.classList.toggle("flex")
    settings.classList.toggle("hidden")
}