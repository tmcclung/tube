// common variables
let iBytesImported = 0
let iBytesTotal = 0
let iPreviousBytesLoaded = 0
let iMaxFilesize = 104857600 // 100MB
let timer = 0
let importInProgress = 'n/a'
let isProcessing = false
let url = null

/* CACHED ELEMENTS */

const importForm = document.getElementById('import-form')
const importInput = document.getElementById('import-input')
const importMessageLabel = document.getElementById('import-message')
const importButtonWrapper = document.getElementById('import-button-wrapper')
const importButton = document.getElementById('import-button')
const importProgressContainer = document.getElementById('import-progress-container')
const importProgressBar = document.getElementById('import-progress')
const importProgressLabel = document.getElementById('import-progress-label')
const importStopped = document.getElementById('import-stopped')
const importStarted = document.getElementById('import-started')

/* HELPERS */

const setProgress = (_progress) => {
    importProgressContainer.style.display = _progress > 0 ? 'flex' : 'none'
    importProgressBar.style.width = `${_progress}%`
    importProgressLabel.innerText =  _progress >= 15 ? `${_progress}%` : ''
}

const setMessage = (_message, isError) => {
    importMessageLabel.style.display = _message ? 'block' : 'none'
    importMessageLabel.innerHTML = _message
    if (isError) {
        importMessageLabel.classList.add('error')
    } else {
        importMessageLabel.classList.remove('error')
    }
}

const setImportState = (_importInProgress) => {
    importInProgress = _importInProgress

    importStarted.style.display = _importInProgress ? 'inline-block' : 'none'
    importStopped.style.display = _importInProgress ? 'none' : 'block'

    if (_importInProgress) {
        importButton.classList.add('transparent')
        timer = setInterval(doInnerUpdates, 300)
    } else {
        importButton.classList.remove('transparent')
        clearInterval(timer)
    }
}

const secondsToTime = (secs) => {
    let hr = Math.floor(secs / 3600)
    let min = Math.floor((secs - (hr * 3600)) / 60)
    let sec = Math.floor(secs - (hr * 3600) - (min * 60))
    if (hr < 10) hr = `0${hr}`
    if (min < 10) min = `0${min}`
    if (sec < 10) sec = `0${sec}`
    if (hr) hr = '00'
    return `${hr}:${min}:${sec}`
}

const bytesToSize = (bytes) => {
    const sizes = ['Bytes', 'KB', 'MB']
    if (bytes == 0) return 'n/a'
    const i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)))
    return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + sizes[i]
}

/* MAIN */

document.addEventListener('DOMContentLoaded', () => {
  importInput.addEventListener("keyup", (e) => {
    if (e.keyCode === 13) {
      e.preventDefault();
      importButton.click();
    }
  });
}, false)

const labelClicked = (e) => {
    if (importInProgress === true) {
        e.preventDefault()
        return false
    }
}

const urlSelected = (_url) => {
    url = _url || importInput.value

    setMessage('')
    setProgress(0)
    
    importButtonWrapper.style.display = 'block'
    setImportState(false)
}

const startImporting = () => {
    if (importInProgress === true) return
    if (!url) return

    isProcessing = false
    iPreviousBytesLoaded = 0
    setMessage('')
    setProgress(0)
    setImportState(true)

    const formData = new FormData()
    formData.append('url', url)
    const xhr = new XMLHttpRequest()

    xhr.upload.addEventListener('progress', importProgress, false)
    xhr.addEventListener('load', importFinish, false)
    xhr.addEventListener('error', importError, false)
    xhr.addEventListener('abort', importAbort, false)

    xhr.open('POST', '/import')
    xhr.send(formData)

    timer = setInterval(doInnerUpdates, 300)
}

const doInnerUpdates = () => { // we will use this function to display import speed
    if (isProcessing) {
        clearInterval(timer)
        return
    }

    let iDiff = iBytesImported - iPreviousBytesLoaded
    // if nothing new loaded - exit
    if (iDiff == 0)
        return
    iPreviousBytesLoaded = iBytesImported
    iDiff = iDiff * 2
    const iBytesRem = iBytesTotal - iPreviousBytesLoaded
    const secondsRemaining = iBytesRem / iDiff
    // update speed info
    let iSpeed = iDiff.toString() + 'B/s'
    if (iDiff > 1024 * 1024) {
        iSpeed = (Math.round(iDiff * 100/(1024*1024))/100).toString() + 'MB/s'
    } else if (iDiff > 1024) {
        iSpeed =  (Math.round(iDiff * 100/1024)/100).toString() + 'KB/s'
    }

    const speedMessage = `${iSpeed} | ${secondsToTime(secondsRemaining)}`
    setMessage(speedMessage)
}

function importProgress(e) { // import process in progress
    if (e.lengthComputable) {
        iBytesImported = e.loaded
        iBytesTotal = e.total

        const iPercentComplete = Math.round(iBytesImported / iBytesTotal * 100)
        setProgress(iPercentComplete)
        if (iPercentComplete === 100) {
            isProcessing = true
            setMessage('Processing video... please wait')
        }
    } else {
        setMessage('Unable to compute progress.')
    }
}

const importFinish = (e) => { // import successfully finished
    const message = e.target.responseText
    const isSuccess = e.target.status < 400

    setProgress(isSuccess ? 100 : 0)
    setMessage(message, !isSuccess)
    setImportState(false)
    if (isSuccess) removeFile(null, true)
}

const importError = () => { // import error
    setMessage('An error occurred while importing the url.', true)
    setProgress(0)
    setImportState(false)
}

const importAbort = () => { // import abort
    setMessage('The import has been canceled by the user or the browser dropped the connection.', true)
    setProgress(0)
    setImportState(false)
}
