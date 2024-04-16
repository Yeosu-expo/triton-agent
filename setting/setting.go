package setting

/* ----- Server Setting ----- */
const ServerPort string = "7000" // Edit this
const ModelsPath string = "./models"

/* ----- Triton Server Setting ----- */
const triton string = "100.0.0.2"
const TritonUrl string = triton + ":8000" // Edit this
const ModelRepository string = "test"
const TritonSSH string = triton + ":22"
const TritonUser string = "root"
const TritonPassword string = "ahri"

/* ----- Scheduler Setting ----- */
// If you are not using a scheduler, change the 'SchedulerActive' variable to false.
const SchedulerActive bool = false           // Edit this
const SchedulerUrl string = "localhost:8000" // Edit this

/* ----- Model Store Setting ----- */
var ModelStoreUrl string = "localhost:8700" // Edit this

const ManagerActive bool = true
const ManagerUrl string = "210.125.31.176:80"

const TorrentUrl string = "host.docker.internal:7001"

var LoadedModel string = ""
